package samplesanalyses

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.RouterGroup) {
	router.POST("", AddAnalysisToSample)
	router.PATCH("/:sample_id/:analysis_id", UpdateSampleAnalysis)
	router.DELETE("/:sample_id/:analysis_id", DeleteSampleAnalysis)
}
