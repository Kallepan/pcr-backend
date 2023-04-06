package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/kaka/pcr-backend/database"
	"gitlab.com/kaka/pcr-backend/models"
)

func RegisterUser(context *gin.Context) {
	// Validate the input from user, hash password and send 201 status code

	var user models.User

	if err := context.ShouldBindJSON(&user); err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := user.HashPassword(user.Password); err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	record := database.Instance.Create(&user)
	if record.Error != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": record.Error.Error()})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"userId": user.ID, "email": user.Email, "username": user.Username})
}
