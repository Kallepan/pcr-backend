package router

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/kallepan/pcr-backend/app/middleware"
	"gitlab.com/kallepan/pcr-backend/config"
)

func Init(init *config.Initialization) *gin.Engine {

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.Cors())

	auth := router.Group("/api")
	{
		auth.POST("/token", init.UserCtrl.LoginUser)
		auth.POST("/register", init.UserCtrl.RegisterUser)
	}

	api := router.Group("/api/v1")
	{
		api.GET("/ping", init.SysCtrl.GetPing)

		// Import
		api.POST("/import", init.ImportCtrl.ImportSample)
	}

	secured := router.Group("/api/v1")
	secured.Use(middleware.AuthMiddleware())
	{
		// Print
		secured.POST("/printer", init.PrintCtrl.PrintSample)

		// Panels
		secured.GET("/panels", init.PanelCtrl.GetPanels)

		// Samples
		samples := secured.Group("/samples")
		samples.GET("", init.SampleCtrl.GetSamples)
		samples.GET("/:sample_id", init.SampleCtrl.GetSamples)
		samples.POST("", init.SampleCtrl.AddSample)
		samples.PUT("/:sample_id", init.SampleCtrl.UpdateSample)
		samples.DELETE("/:sample_id", init.SampleCtrl.DeleteSample)

		// Sample Panels
		samplePanels := secured.Group("/samplespanels")
		samplePanels.POST("", init.SamplePanelCtrl.CreateSamplePanel)
		samplePanels.POST("/reset", init.SamplePanelCtrl.ResetSamplePanel)
		samplePanels.PATCH("/:sample_id/:panel_id", init.SamplePanelCtrl.UpdateSamplePanel)
		samplePanels.GET("", init.SamplePanelCtrl.GetSamplePanels)
		samplePanels.GET("/stats", init.SamplePanelCtrl.GetStatistics)
		samplePanels.POST("/create-run", init.SamplePanelCtrl.CreateRun)
	}

	return router
}
