package analyses

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.RouterGroup) {
	router.GET("", GetAllAnalyses)
	router.POST("", AddAnalysis)
	router.PUT("/:analysis_id", UpdateAnalysis)
	router.DELETE("/:analysis_id", DeleteAnalysis)
}
