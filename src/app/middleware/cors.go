package middleware

import (
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func compareOrigin(c *gin.Context) string {
	/* Compares the origin of the request with the allowed origins */
	origin := c.Request.Header.Get("Origin")

	allowedOrigins := os.Getenv("ALLOWED_ORIGINS")

	// If allowedOrigins is empty, allow all origins
	if allowedOrigins == "" {
		return "*"
	}

	// If allowedOrigins is not empty, compare the origin with the allowed origins
	allowedOriginsArr := strings.Split(allowedOrigins, ",")
	for _, allowedOrigin := range allowedOriginsArr {
		if origin == allowedOrigin {
			return origin
		}
	}

	// If the origin is not allowed, return empty string
	return ""
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Set the allowed origins
		origin := compareOrigin(c)

		// Set the allowed headers
		c.Header("Access-Control-Allow-Origin", origin)
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		c.Header("Access-Control-Allow-Credentials", "true")

		// If the request method is OPTIONS, return with status 200
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}

		// Continue to the next middleware
		c.Next()
	}
}
