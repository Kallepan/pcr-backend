package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/kaka/pcr-backend/auth"
	"gitlab.com/kaka/pcr-backend/database"
	"gitlab.com/kaka/pcr-backend/models"
)

type TokenRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func GenerateJWTToken(context *gin.Context) {
	var request TokenRequest
	var user models.User

	if err := context.ShouldBindJSON(&request); err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	// Check if user exists and password is correct
	record := database.Instance.Where("username = ?", request.Username).First(&user)
	if record.Error != nil {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	credentialsError := user.CheckPassword(request.Password)
	if credentialsError != nil {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	tokenString, err := auth.GenerateJWTToken(user.Username, user.Email)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "error generating token"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"token": tokenString})
}
