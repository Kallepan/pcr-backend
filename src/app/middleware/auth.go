package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/kallepan/pcr-backend/auth"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString := ctx.GetHeader("Authorization")
		if tokenString == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Bitte zuerst Einloggen"})
			return
		}

		err := auth.ValidateJWTToken(tokenString)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Session ist abgelaufen. Bitte neu einloggen"})
			return
		}

		// Set user_id if possible
		user_id, err := auth.GetUserIdFromToken(ctx.GetHeader("Authorization"))
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Session ist ungültig. Bitte neu einloggen"})
			return
		}
		ctx.Set("user_id", user_id)

		ctx.Next()
	}
}
