package samples

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.RouterGroup) {
	router.GET("", GetSamples)
	router.GET("/:tagesnummer", GetSample)
	router.POST("", AddSample)
	router.PUT("/:tagesnummer", UpdateSample)
	router.DELETE("/:tagesnummer", DeleteSample)
}
