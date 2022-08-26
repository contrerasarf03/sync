package sync

import (
	"os"
	"os/signal"
)

// WaitForInterrupt will wait until an OS interrupt occurs (eg. Ctrl+C) then calls
// the cleanup functions for cleanup before closing the application. Should the
// functions returns an error, an exit error will be returned instead of zero.
func WaitForInterrupt(cleanups ...func() error) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	for _, cleanup := range cleanups {
		if err := cleanup(); err != nil {
			os.Exit(1)
		}
	}
	os.Exit(0)
}
