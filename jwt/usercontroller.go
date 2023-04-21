package jwt

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/kaka/pcr-backend/common/database"
	"gitlab.com/kaka/pcr-backend/common/models"
)

func CheckIfUserExists(username string) bool {
	var exists bool

	query := `
		SELECT EXISTS (
			SELECT 1 FROM users WHERE username = $1
		);
	`
	err := database.Instance.QueryRow(query, username).Scan(&exists)

	if err != nil {
		return false
	}

	return exists
}

func RegisterUser(context *gin.Context) {
	// Validate the input from user, hash password and send 201 status code

	var user models.User

	if err := context.ShouldBindJSON(&user); err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if CheckIfUserExists(user.Username) {
		error_message := fmt.Sprintf("User with username %s already exists", user.Username)
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": error_message})
		return
	}

	if err := user.HashPassword(user.Password); err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	query := "INSERT INTO users (email, firstname, lastname, username, password) VALUES ($1, $2, $3, $4, $5) RETURNING user_id"
	err := database.Instance.QueryRow(query, &user.Email, &user.FirstName, &user.LastName, &user.Username, &user.Password).Scan(&user.UserId)

	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"userId": user.UserId, "email": user.Email, "username": user.Username, "firstname": user.FirstName, "lastname": user.LastName})
}
