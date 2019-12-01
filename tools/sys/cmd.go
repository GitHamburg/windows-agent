package sys

import (
	"os/exec"
	"time"
	"log"
	"os"
)

func CmdRunWithTimeout(cmd *exec.Cmd, timeout time.Duration) (error, bool) {
	done := make(chan error)
	go func() {
		done <- cmd.Wait()
	}()

	var err error
	select {
	case <-time.After(timeout):
		log.Printf("timeout, process:%s will be killed", cmd.Path)

		go func() {
			<-done // allow goroutine to exit
		}()

		// timeout
		if err = cmd.Process.Signal(os.Kill); err != nil {
			log.Printf("failed to kill: %s, error: %s", cmd.Path, err)
		}

		return err, true
	case err = <-done:
		return err, false
	}
}
