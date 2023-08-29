package main

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/kaka/pcr-backend/common/controllers"
	"gitlab.com/kaka/pcr-backend/common/database"
	"gitlab.com/kaka/pcr-backend/common/middlewares"
	"gitlab.com/kaka/pcr-backend/jwt"
	"gitlab.com/kaka/pcr-backend/packages/importer"
	"gitlab.com/kaka/pcr-backend/packages/panels"
	"gitlab.com/kaka/pcr-backend/packages/printer"
	"gitlab.com/kaka/pcr-backend/packages/samples"
	"gitlab.com/kaka/pcr-backend/packages/samplespanels"
	"gitlab.com/kaka/pcr-backend/utils"
)

func main() {
	connectionString := utils.GetDBConnectionString()
	database.Connect(connectionString)
	database.Migrate()
	defer database.Instance.Close()

	jwt.CreateAdminUser()

	interval := time.Minute * 5
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
		v1.GET("/samplespanels/:sample_id", samplespanels.GetSamplesPanels)
		v1.GET("/samplespanels", samplespanels.GetSamplesPanels)

		// analyses
		v1.GET("/panels", panels.GetPanels)
		v1.GET("/panels/:panel_id", panels.GetPanels)

		// ping
		v1.GET("/ping", controllers.Ping)

		// importer
		v1.POST("/import", importer.PostSampleMaterial)
	}

	secured := router.Group("/api/v1")
	secured.Use(middlewares.AuthMiddleware())
	{
		samples.RegisterRoutes(secured.Group("/samples"))
		samplespanels.RegisterRoutes(secured.Group("/samplespanels"))
		printer.RegisterRoutes(secured.Group("/printer"))
	}

	return router
}
