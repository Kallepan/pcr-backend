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

	query := "INSERT INTO users (email, firstname, lastname, username, password) VALUES ($1, $2, $3, $4, $5) RETURNING user_id"
	err := database.Instance.QueryRow(query, user.Email, user.FirstName, user.LastName, user.Username, user.Password).Scan(&user.UserId)

	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"userId": user.UserId, "email": user.Email, "username": user.Username})
}
