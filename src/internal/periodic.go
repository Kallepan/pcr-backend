package internal

import (
	"log/slog"
	"sync"
	"time"

	"gitlab.com/kallepan/pcr-backend/app/repository"
)

var syncLock sync.Mutex

func initSynchronization(synchronizeRepo repository.SynchronizeRepository, interval time.Duration) {
	slog.Info("Synchronization started")

	go func() {
		syncLock.Lock()
		defer syncLock.Unlock()
		synchronizeRepo.Synchronize()
		time.Sleep(interval)
	}()

}
