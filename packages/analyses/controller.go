package analyses

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.RouterGroup) {
	router.POST("", AddAnalysis)
	router.PUT("/:analysis_id", UpdateAnalysis)
}
