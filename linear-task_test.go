package sync

import (
	"errors"
	"testing"
	"time"
)

func TestLinearTask(t *testing.T) {
	counter := 0
	task := func() error {
		time.Sleep(1 * time.Second)
		counter++ // should be safe since this is syncrhonous
		t.Log("task done")
		return nil
	}
	endTask, taskEnded := LinearTask(0, true, task)
	<-time.After(3 * time.Second)
	close(endTask)
	<-taskEnded
	if counter > 3 {
		t.Fail()
	}
}

func TestLinearTaskFail(t *testing.T) {
	counter := 0
	task := func() error {
		time.Sleep(1 * time.Second)
		counter++ // should be safe since this is syncrhonous
		t.Log("task done")
		return errors.New("something went wrong")
	}
	endTask, taskEnded := LinearTask(0, true, task)
	<-time.After(2 * time.Second)
	close(endTask)
	<-taskEnded
	if counter > 1 {
		t.Fail()
	}
}
