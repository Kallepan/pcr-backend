package internal

import (
	"time"

	"gitlab.com/kallepan/pcr-backend/config"
)

func Init(init *config.Initialization, interval time.Duration) {
	initSynchronization(init.SynchroRepo, interval)
}
