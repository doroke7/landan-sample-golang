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
		// Step2: 檢查pid，如果沒有 pid 就 執行
		pid, ok := readPid()
		if !ok || !runningLauncher(pid, exe) {
			if pid, err = spawnLauncher(exe, args); err != nil {
				return err
			}
			log.Printf("[supervisor] 已啟動 (pid %d)", pid)
		} else {
			log.Printf("[supervisor] 沿用執行中的副程序 (pid %d)", pid)
		}

		// Step3:
		// 每隔 1 秒用 pid 檢查副程序還在不在; 在就持續檢查，外部卡住
		// 每隔 1 秒用 pid 檢查副程序還在不在; 不在就持續檢查，回傳 false
		// 同時等停止訊號的 channel。
		if stopped := waitLauncher(pid, exe, sig); stopped {
			stopLauncher(pid, exe)
			return nil
		}

		// Step4:pid 不在了(副程序死了),就等待 5 秒後再一次迴圈。
		log.Printf("[supervisor] pid %d 已結束,%s 後重啟", pid, 5*time.Second)
		os.Remove(pidPath())
		select {
		case <-sig:
			return nil
		case <-time.After(5 * time.Second):
		}
	}
}

/*
有三種情況
 1. 如果 pid 正常運行，責會一直跑迴圈定時檢查，不回傳，外面會看起來卡住。
 2. 如果收到 ctrl + C 訊號，就回傳stop 為 true
 3. 如果pid 運行不正常，就回傳 stop false
*/
func waitLauncher(pid int, exe string, sig <-chan os.Signal) (stopped bool) {

	for runningLauncher(pid, exe) {
		select {
		case <-sig:
			return true
		case <-time.After(time.Second):
		}
	}
	return false
}

// spawnLauncher runs exe with args in the background with nohup and returns its pid, which is also written to the pid file.
// The output goes to the log file; nohup makes it ignore SIGHUP, so closing the terminal does not kill it.
// $0 is the log file, "$@" is the command to run.
func spawnLauncher(exe string, args []string) (int, error) {

	aShells := append([]string{"-c", `nohup "$@" >>"$0" 2>&1 </dev/null & echo $!`, logPath(), exe}, args...)
	out, err := exec.Command("sh", aShells...).Output()
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

// runningLauncher is true if pid is alive AND is our executable (a stale pid file may point at an unrelated process).
func runningLauncher(pid int, exe string) bool {

	out, err := exec.Command("ps", "-p", strconv.Itoa(pid), "-o", "command=").Output()
	return err == nil && strings.Contains(string(out), exe)
}

// stopLauncher asks the background process to quit by its pid and waits for it; it is killed if it does not quit in time.
func stopLauncher(pid int, exe string) {

	_ = syscall.Kill(pid, syscall.SIGTERM)
	for i := 0; i < int(5*time.Second/(100*time.Millisecond)) && runningLauncher(pid, exe); i++ {
		time.Sleep(100 * time.Millisecond)
	}
	if runningLauncher(pid, exe) {
		_ = syscall.Kill(pid, syscall.SIGKILL)
	}
	os.Remove(pidPath())
}
