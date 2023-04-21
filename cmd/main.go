package main

import (
	"strings"

	"github.com/gin-gonic/gin"
	"gitlab.com/kaka/pcr-backend/analyses"
	"gitlab.com/kaka/pcr-backend/common/controllers"
	"gitlab.com/kaka/pcr-backend/common/database"
	"gitlab.com/kaka/pcr-backend/common/middlewares"
	"gitlab.com/kaka/pcr-backend/jwt"
	"gitlab.com/kaka/pcr-backend/samples"
	"gitlab.com/kaka/pcr-backend/utils"
)

func main() {
	connectionString := utils.GetConnectionString()
	database.Connect(connectionString)

	router := initRouter()
	router.Run(":8080")
}

func initRouter() *gin.Engine {
	router := gin.Default()

	router.Use(middlewares.ErrorHandler)

	router.SetTrustedProxies(strings.Split(utils.GetValueFromEnv("TRUSTED_PROXIES", ","), ","))

	auth := router.Group("/api")
	{
		auth.POST("/token", jwt.GenerateJWTTokenController)
		auth.POST("/register", jwt.RegisterUser)
	}

	api := router.Group("/api/v1")
	api.Use(middlewares.AuthMiddleware())
	{
		samples.RegisterRoutes(api.Group("/samples"))
		analyses.RegisterRoutes(api.Group("/analyses"))
		api.GET("/ping", controllers.Ping)
	}

	return router
}
