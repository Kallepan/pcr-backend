package analyses

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.RouterGroup) {
	router.GET("", GetAllAnalyses)
	router.POST("", AddAnalysis)
	router.PUT("/:analyt/:material/:assay", UpdateAnalysis)
	router.DELETE("/:analyt/:material/:assay", DeleteAnalysis)
}
