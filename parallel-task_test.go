package sync

import (
	"errors"
	"testing"
	"time"
)

func TestParallel(t *testing.T) {

	// long running task is just waiting 2 seconds
	task := func() error { longRunningFunction(t, 2*time.Second); return nil }

	// parallel run 10 tasks and get duration
	start := time.Now()
	stop, stopped := ParallelTask(10, 0, true, task)

	// wait for 5 seconds, this should have queued the initial 10 tasks, and another 10 more tasks
	<-time.After(5 * time.Second)
	close(stop)
	<-stopped
	end := time.Now()
	duration := end.Sub(start)
	t.Log("duration", duration)
	if duration > 7*time.Second {
		t.Error("since this is parallel work, should really finish less than 7 secods (4 secs + waiting time")
	}

}

func TestParallelError(t *testing.T) {

	// long running task is just waiting 2 seconds
	task := func() error { longRunningFunction(t, 2*time.Second); return errors.New("something went wrong") }

	// parallel run 10 tasks and get duration
	start := time.Now()
	stop, stopped := ParallelTask(10, 0, true, task)

	// wait for 1 seconds, this should have queued the initial 10 tasks, and nothing more
	<-time.After(3 * time.Second)
	close(stop)
	<-stopped
	end := time.Now()
	duration := end.Sub(start)
	t.Log("duration", duration)
	if duration > 4*time.Second {
		t.Error("this should not be more than 4 or else the error trap did not work and a new batch of jobs were ran")
	}

}

func longRunningFunction(t *testing.T, duration time.Duration) {
	start := time.Now()
	time.Sleep(duration)
	end := time.Now()
	t.Logf("completed task @ %v", end.Sub(start))
}
