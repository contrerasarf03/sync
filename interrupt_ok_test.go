package sync

import (
	"os"
	"syscall"
	"testing"
	"time"

	"os/exec"
)

func TestInterrupt(t *testing.T) {
	if os.Getenv("FLAG") == "1" {
		task := func() error {
			return nil
		}
		go WaitForInterrupt(task)
		<-time.After(1 * time.Second)
		pid := syscall.Getpid()
		p, err := os.FindProcess(pid)
		if err != nil {
			t.Fail()
		}
		p.Signal(os.Interrupt)
		t.Fail()
	}

	cmd := exec.Command("go", "test", "-v", "-test.run=TestInterrupt")
	cmd.Env = append(os.Environ(), "FLAG=1")
	err := cmd.Run()
	t.Log(err)
	if err != nil {
		t.Fail()
	}
}
