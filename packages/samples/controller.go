package samples

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/:sample_id/analyses", GetAnalysesForSample)
	router.POST("", AddSample)
	router.PUT("/:sample_id", UpdateSample)
	router.DELETE("/:sample_id", DeleteSample)
}
