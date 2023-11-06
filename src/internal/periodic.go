package internal

import (
	"log/slog"
	"time"

	"gitlab.com/kallepan/pcr-backend/app/repository"
)

func initSynchronization(synchronizeRepo repository.SynchronizeRepository, interval time.Duration) {
	slog.Info("Synchronization started")

	go func() {
		for {
			synchronizeRepo.Synchronize()
			time.Sleep(interval)
		}
	}()

}
