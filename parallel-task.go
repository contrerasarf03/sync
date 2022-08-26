package sync

import (
	"sync"
	"time"
)

// ParallelTask (for lack of a better work) will run the provided `task` in parallel by
// spinning goroutines to the number of workers specified. This function returns a channel
// that can be used to stop the workers.
//
// task should be wrapped in a plain function and should not spin up goroutines or else
// unexpected behaviour will occur. (eg. more than the specified workers will be running and
// may end up spinning too many gorouties).
//
// This is a simple implementation where it uses sync.WaitGroup to create create goroutines
// up to the workers specified. Each goroutine will call the task function then mark the job
// done before it creates another set of goroutines after the specified intervalDuration.
//
// Sending a stop signal via the channel will stop the parallel work. Any running parallel
// work will have to completed before this actually stops.
//
// the stopped channel will be closed when the parallel work has stopped.
//
// when an error occurs in one of the tasks, the parallel work also attempt to stop (no more new jobs)
func ParallelTask(workers int, intervalDuration time.Duration, failOnError bool, task func() error) (stop chan interface{}, stopped chan interface{}) {
	stop = make(chan interface{})
	stopped = make(chan interface{})
	errChan := make(chan error, workers)
	go func() {
		hasStopped := false
		wg := &sync.WaitGroup{}
		for !hasStopped {
			select {
			case <-stop:
				hasStopped = true
			case <-errChan:
				if failOnError {
					hasStopped = true
				}
			default:
				processTaskInParallel(wg, workers, errChan, task)
			}
			time.Sleep(intervalDuration)
		}
		wg.Wait()
		close(stopped)
	}()
	return
}

func processTaskInParallel(wg *sync.WaitGroup, workers int, errChan chan error, task func() error) {
	wg.Add(workers)

	for i := 0; i < workers; i++ {
		go func() {
			if err := task(); err != nil {
				errChan <- err
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
