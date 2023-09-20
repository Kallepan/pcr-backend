package main

import (
	"context"
	"os"

	"gitlab.com/kallepan/pcr-backend/app/router"
	"gitlab.com/kallepan/pcr-backend/auth"
	"gitlab.com/kallepan/pcr-backend/config"
	"gitlab.com/kallepan/pcr-backend/driver"
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
	app := router.Init(init)

	app.Run(":" + port)
}
