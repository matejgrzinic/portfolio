package refresher

import (
	"time"
)

func StartRefresher(errorChan chan error, stopChan chan struct{}, interval time.Duration, f func() error) {
	go func() {
		sendError(errorChan, f())
		for {
			select {
			case <-time.After(interval):
				sendError(errorChan, f())
			case <-stopChan:
				return
			}
		}
	}()
}

func sendError(errorChan chan error, err error) {
	select {
	case errorChan <- err:
		return
	default:
		return
	}
}
