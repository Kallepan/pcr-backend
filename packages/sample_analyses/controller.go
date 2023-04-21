package sampleanalyses

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.RouterGroup) {
	router.GET("", GetSampleAnalyses)
	router.POST("", AddAnalysisToSample)
	router.PATCH("/:sample_id/:analysis_id", UpdateSampleAnalysis)
}
