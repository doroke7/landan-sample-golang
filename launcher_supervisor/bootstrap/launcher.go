package bootstrap

import (
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"
)

const (
	restartWait = 5 * time.Second // 副程序結束後,等多久才重啟
	checkEvery  = time.Second     // 多久檢查一次副程序的 pid 還在不在
	stopWait    = 5 * time.Second // 收到停止訊號後等副程序收尾多久,超過就強制結束
)

var dir = filepath.Join("runtime", "launcher_supervisor") // 副程序的 pid 檔和日誌都放在這裡

func pidPath() string { return filepath.Join(dir, "child.pid") }
func logPath() string { return filepath.Join(dir, "child.log") }

// SuperviseLauncher runs this executable again with args as a background process (nohup), keeps its pid in a file,
// and runs it again whenever that pid is gone, until this process gets SIGTERM or Ctrl-C
// (which is passed on to the background process by its pid).
func SuperviseLauncher(args ...string) error {

	exe, err := os.Executable()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, os.Interrupt)
	defer signal.Stop(sig)

	// Step1:開 for loop,不斷的重試。
	for {
		// Step2:副程序是背景程式,沒有 Wait() 可以等,所以一切靠 pid。
		// pid 檔裡的副程序還活著(例如上一個主程序留下的)就沿用它,不開第二個;否則以 nohup 啟動,把 pid 寫進 pid 檔。
		pid, ok := readPid()
		if !ok || !running(pid, exe) {
			if pid, err = spawn(exe, args); err != nil {
				return err
			}
			log.Printf("[supervisor] 已啟動 (pid %d)", pid)
		} else {
			log.Printf("[supervisor] 沿用執行中的副程序 (pid %d)", pid)
		}

		// Step3:每隔 checkEvery 用 pid 檢查副程序還在不在;同時等停止訊號。
		if stopped := waitGone(pid, exe, sig); stopped {
			stopLauncher(pid, exe)
			return nil
		}

		// Step4:pid 不在了(副程序死了),就等待 restartWait 後再一次迴圈。
		log.Printf("[supervisor] pid %d 已結束,%s 後重啟", pid, restartWait)
		os.Remove(pidPath())
		select {
		case <-sig:
			return nil
		case <-time.After(restartWait):
		}
	}
}

// waitGone blocks until pid is gone (returns false) or a stop signal arrives (returns true).
func waitGone(pid int, exe string, sig <-chan os.Signal) (stopped bool) {

	for running(pid, exe) {
		select {
		case <-sig:
			return true
		case <-time.After(checkEvery):
		}
	}
	return false
}

// spawn runs exe with args in the background with nohup and returns its pid, which is also written to the pid file.
// The output goes to the log file; nohup makes it ignore SIGHUP, so closing the terminal does not kill it.
// $0 is the log file, "$@" is the command to run.
func spawn(exe string, args []string) (int, error) {

	shell := append([]string{"-c", `nohup "$@" >>"$0" 2>&1 </dev/null & echo $!`, logPath(), exe}, args...)
	out, err := exec.Command("sh", shell...).Output()
	if err != nil {
		return 0, err
	}
	pid, err := strconv.Atoi(strings.TrimSpace(string(out)))
	if err != nil {
		return 0, err
	}
	return pid, os.WriteFile(pidPath(), []byte(strconv.Itoa(pid)), 0o644)
}

func readPid() (int, bool) {
	data, err := os.ReadFile(pidPath())
	if err != nil {
		return 0, false
	}
	pid, err := strconv.Atoi(strings.TrimSpace(string(data)))
	return pid, err == nil && pid > 0
}

// running is true if pid is alive AND is our executable (a stale pid file may point at an unrelated process).
func running(pid int, exe string) bool {

	out, err := exec.Command("ps", "-p", strconv.Itoa(pid), "-o", "command=").Output()
	return err == nil && strings.Contains(string(out), exe)
}

// stopLauncher asks the background process to quit by its pid and waits for it; it is killed if it does not quit in time.
func stopLauncher(pid int, exe string) {

	_ = syscall.Kill(pid, syscall.SIGTERM)
	for i := 0; i < int(stopWait/(100*time.Millisecond)) && running(pid, exe); i++ {
		time.Sleep(100 * time.Millisecond)
	}
	if running(pid, exe) {
		_ = syscall.Kill(pid, syscall.SIGKILL)
	}
	os.Remove(pidPath())
}
