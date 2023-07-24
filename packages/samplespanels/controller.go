package samplespanels

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.RouterGroup) {
	router.POST("", AddAnalysisToSample)
	router.POST("/reset", ResetSamplePanel)
	router.PATCH("/:sample_id/:panel_id", UpdateSamplePanel)
	router.DELETE("/:sample_id/:panel_id", DeleteSamplePanel)
	router.POST("/create-run", CreateRun)
}
