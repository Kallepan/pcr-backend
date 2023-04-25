package main

import (
	"strings"

	"github.com/gin-gonic/gin"
	"gitlab.com/kaka/pcr-backend/common/controllers"
	"gitlab.com/kaka/pcr-backend/common/database"
	"gitlab.com/kaka/pcr-backend/common/middlewares"
	"gitlab.com/kaka/pcr-backend/jwt"
	"gitlab.com/kaka/pcr-backend/packages/analyses"
	sampleanalyses "gitlab.com/kaka/pcr-backend/packages/sample_analyses"
	"gitlab.com/kaka/pcr-backend/packages/samples"
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

	v1 := router.Group("/api/v1")
	{
		v1.GET("/samples", samples.GetSamples)
		v1.GET("/sampleanalyses", sampleanalyses.GetSampleAnalyses)
		v1.GET("/analyses", analyses.GetAllAnalyses)
		v1.GET("/ping", controllers.Ping)
	}

	secured := router.Group("/api/v1")
	secured.Use(middlewares.AuthMiddleware())
	{
		samples.RegisterRoutes(secured.Group("/samples"))
		analyses.RegisterRoutes(secured.Group("/analyses"))
		sampleanalyses.RegisterRoutes(secured.Group("/sampleanalyses"))
	}

	return router
}
