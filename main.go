package main

import (
	"strings"

	"github.com/gin-gonic/gin"
	"gitlab.com/kaka/pcr-backend/controllers"
	"gitlab.com/kaka/pcr-backend/database"
	"gitlab.com/kaka/pcr-backend/middlewares"
	"gitlab.com/kaka/pcr-backend/utils"
)

func main() {
	connectionString := utils.GetConnectionString()
	database.Connect(connectionString)
	database.Migrate()

	router := initRouter()
	router.Run(":8080")
}

func initRouter() *gin.Engine {
	router := gin.Default()

	router.SetTrustedProxies(strings.Split(utils.GetValueFromEnv("TRUSTED_PROXIES", ","), ","))

	auth := router.Group("/api")
	{
		auth.POST("/token", controllers.GenerateJWTToken)
		auth.POST("/register", controllers.RegisterUser)
	}

	api := router.Group("/api/v1")
	{
		secured := api.Group("/secured").Use(middlewares.AuthMiddleware())
		{
			secured.GET("/ping", controllers.Ping)
		}
	}

	return router
}
