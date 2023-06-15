package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HttpErrorResponse struct {
	Message     string `json:"message"`
	Status      int    `json:"status"`
	Description string `json:"description"`
}

func ErrorHandler(ctx *gin.Context) {
	ctx.Next()

	for _, err := range ctx.Errors {

		if err.Err == nil {
			continue
		}

		httpErrorResponse := HttpErrorResponse{
			Message:     "Internal Server Error",
			Status:      http.StatusInternalServerError,
			Description: err.Err.Error(),
		}

		ctx.AbortWithStatusJSON(httpErrorResponse.Status, httpErrorResponse)
		return
	}
}
