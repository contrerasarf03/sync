package sync

import (
	"time"
)

// LinearTask will run a task continuously sleeping at a provided interval.
// This returns a channel that can be used to stop the task by closing it
// or sending any message to it. When the task is ended, the taskEnded channel
// will be closed.
func LinearTask(interval time.Duration, failOnError bool, task func() error) (endTask, taskEnded chan interface{}) {
	endTask = make(chan interface{})
	taskEnded = make(chan interface{})
	hasStopped := false
	go func() {
		for !hasStopped {
			select {
			case <-endTask:
				hasStopped = true
				close(taskEnded)
			default:
				if err := task(); err != nil {
					if failOnError {
						hasStopped = true
						close(taskEnded)
					}
				}
			}
			time.Sleep(interval)
		}
	}()
	return
}
