package bootstrap

import (
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
)

const (
	restartWait = 5 * time.Second // 副程序結束後,等多久才重啟
	stopWait    = 5 * time.Second // 收到停止訊號後等子程序收尾多久,超過就強制結束
)

// SuperviseLauncher runs this executable again with args as a child, and runs it again whenever it exits,
// until this process gets SIGTERM or Ctrl-C (which is passed on to the child).
func SuperviseLauncher(args ...string) error {

	exe, err := os.Executable()
	if err != nil {
		return err
	}

	sig := make(chan os.Signal, 1)
	/*
		os.Interrupt    │ SIGINT（2）    │ 在終端機按 Ctrl-C
		syscall.SIGTERM │ SIGTERM（15）  │ kill <pid>（
	*/
	signal.Notify(sig, syscall.SIGTERM, os.Interrupt)
	defer signal.Stop(sig)

	// Step1:開 for loop,不斷的重試。
	for {
		// Step2:主程序以 os command 執行副程序(副程序才是主要商務應用)。
		oCommand := exec.Command(exe, args...)
		oCommand.Stdout = os.Stdout // 子程序的標準輸出 → 副程序自己的標準輸出
		oCommand.Stderr = os.Stderr // 子程序的標準錯誤 → 副程序自己的標準錯誤

		if err := oCommand.Start(); err != nil {
			return err
		}
		log.Printf("[supervisor] 已啟動 (pid %d)", oCommand.Process.Pid)

		// Step3:開一個協程,用 chan done 接收 oCommand.Wait();底下用 select 等待 done。
		done := make(chan error, 1)
		go func() {
			done <- oCommand.Wait()
		}()

		select {
		case <-sig:
			stopLauncher(oCommand, done)
			return nil

		case err := <-done:
			// Step4:done 等到了(副程序死了的訊號),就等待 restartWait 後再一次迴圈。
			log.Printf("[supervisor] 已結束 (%v),%s 後重啟", err, restartWait)
		}

		select {
		case <-sig:
			return nil
		case <-time.After(restartWait):
			// 等待
		}
	}
}

// stopLauncher asks the child to quit and waits for it; it is killed if it does not quit in time.
func stopLauncher(oCommand *exec.Cmd, done <-chan error) {

	if err := oCommand.Process.Signal(syscall.SIGTERM); err != nil {
		_ = oCommand.Process.Kill() // Windows has no SIGTERM
	}

	select {
	case <-done:
	case <-time.After(stopWait):
		_ = oCommand.Process.Kill()
		<-done
	}
}
