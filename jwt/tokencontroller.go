package jwt

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/kaka/pcr-backend/common/database"
	"gitlab.com/kaka/pcr-backend/common/models"
)

type TokenRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func GenerateJWTTokenController(context *gin.Context) {
	var request TokenRequest
	var user models.User

	if err := context.ShouldBindJSON(&request); err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	// Check if user exists and password is correct
	query := "SELECT username, password, email, user_id FROM users WHERE username = $1"
	err := database.Instance.QueryRow(query, request.Username).Scan(&user.Username, &user.Password, &user.Email, &user.UserId)

	if err != nil {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials1"})
		return
	}
	credentialsError := user.CheckPassword(request.Password)
	if credentialsError != nil {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	tokenString, err := GenerateJWTToken(user.Username, user.Email, user.UserId)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "error generating token"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"token": tokenString})
}
