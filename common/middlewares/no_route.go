package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func NoRouteHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusNotFound, gin.H{"message": "Not Found"})
}
