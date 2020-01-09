package sys

import (
	"bytes"
	"log"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

func CmdOut(name string, arg ...string) (string, error) {
	cmd := exec.Command(name, arg...)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	return out.String(), err
}

func CmdOutBytes(name string, arg ...string) ([]byte, error) {
	cmd := exec.Command(name, arg...)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	return out.Bytes(), err
}

func CmdOutNoLn(name string, arg ...string) (out string, err error) {
	out, err = CmdOut(name, arg...)
	if err != nil {
		return
	}
	return strings.TrimSpace(string(out)), nil
}

func CmdRunWithTimeout(cmd *exec.Cmd, timeout time.Duration) (error, bool) {
	var err error

	done := make(chan error)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case <-time.After(timeout):
		log.Printf("timeout, process:%s will be killed", cmd.Path)

		go func() {
			<-done // allow goroutine to exit
		}()

		/*
			将PGID设置成与PID相同的值，避免fork子进程杀不掉
			IMPORTANT: cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true} is nessary before cmd.Start()
			https://www.jianshu.com/p/1f3ec2f00b03
		*/
		err = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
		if err != nil {
			log.Println("kill failed, error:", err)
		}
		return err, true
	case err = <-done:
		return err, false
	}

}
