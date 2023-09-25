package main

import (
	"context"
	"os"
	"time"

	"gitlab.com/kallepan/pcr-backend/app/router"
	"gitlab.com/kallepan/pcr-backend/auth"
	"gitlab.com/kallepan/pcr-backend/config"
	"gitlab.com/kallepan/pcr-backend/driver"
	"gitlab.com/kallepan/pcr-backend/internal"
)

func init() {
	config.InitLog()
}

func main() {
	ctx := context.Background()
	port := os.Getenv("PORT")

	driver.Init(ctx)
	auth.Init()
	init := config.Init()
	internal.Init(init, 5*time.Minute)
	app := router.Init(init)

	app.Run(":" + port)
}
