package main

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/kaka/pcr-backend/common/controllers"
	"gitlab.com/kaka/pcr-backend/common/database"
	"gitlab.com/kaka/pcr-backend/common/middlewares"
	"gitlab.com/kaka/pcr-backend/jwt"
	"gitlab.com/kaka/pcr-backend/packages/panels"
	"gitlab.com/kaka/pcr-backend/packages/samples"
	"gitlab.com/kaka/pcr-backend/packages/samplespanels"
	"gitlab.com/kaka/pcr-backend/utils"
)

func main() {
	connectionString := utils.GetConnectionString()
	database.Connect(connectionString)
	database.Migrate()
	defer database.Instance.Close()

	jwt.CreateAdminUser()

	interval := time.Minute * 7
	samplespanels.StartSynchronize(interval)

	router := initRouter()
	router.Run(":8080")
}

func initRouter() *gin.Engine {
	router := gin.Default()

	router.NoRoute(middlewares.NoRouteHandler)
	router.Use(middlewares.ErrorHandler)
	router.Use(middlewares.CORSMiddleware())
	router.SetTrustedProxies(strings.Split(utils.GetValueFromEnv("TRUSTED_PROXIES", ","), ","))

	auth := router.Group("/api")
	{
		auth.POST("/token", jwt.GenerateJWTTokenController)
		auth.POST("/register", jwt.RegisterUser)
	}

	v1 := router.Group("/api/v1")
	{
		// samples
		v1.GET("/samples/:sample_id", samples.GetSamples)
		v1.GET("/samples", samples.GetSamples)

		// samples-panels
		v1.GET("/samples-panels/:sample_id", samplespanels.GetSamplesPanels)
		v1.GET("/samples-panels", samplespanels.GetSamplesPanels)

		// analyses
		v1.GET("/panels", panels.GetPanels)
		v1.GET("/panels/:panel_id", panels.GetPanels)

		// ping
		v1.GET("/ping", controllers.Ping)
	}

	secured := router.Group("/api/v1")
	secured.Use(middlewares.AuthMiddleware())
	{
		samples.RegisterRoutes(secured.Group("/samples"))
		samplespanels.RegisterRoutes(secured.Group("/samples-panels"))
	}

	return router
}
