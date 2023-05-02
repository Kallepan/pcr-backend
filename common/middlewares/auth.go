package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/kaka/pcr-backend/jwt"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString := ctx.GetHeader("Authorization")
		if tokenString == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Bitte zuerst Einloggen"})
			return
		}

		err := jwt.ValidateJWTToken(tokenString)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Session ist abgelaufen. Bitte neu einloggen"})
			return
		}

		// Set user_id if possible
		user_id, err := jwt.GetUserIdFromToken(ctx.GetHeader("Authorization"))
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Session ist ungültig. Bitte neu einloggen"})
			return
		}
		ctx.Set("user_id", user_id)

		ctx.Next()
	}
}
