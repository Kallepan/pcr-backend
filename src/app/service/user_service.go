package service

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"gitlab.com/kallepan/pcr-backend/app/constant"
	"gitlab.com/kallepan/pcr-backend/app/domain/dao"
	"gitlab.com/kallepan/pcr-backend/app/pkg"
	"gitlab.com/kallepan/pcr-backend/app/repository"
	"gitlab.com/kallepan/pcr-backend/auth"
)

type UserService interface {
	RegisterUser(ctx *gin.Context)
	LoginUser(ctx *gin.Context)
}

type UserServiceImpl struct {
	userRepository repository.UserRepository
}

func UserServiceInit(userRepository repository.UserRepository) *UserServiceImpl {
	return &UserServiceImpl{
		userRepository: userRepository,
	}
}

func (u UserServiceImpl) LoginUser(ctx *gin.Context) {
	defer pkg.PanicHandler(ctx)
	slog.Info("Received request to login user")

	var request dao.TokenRequest

	// Validate request
	if err := ctx.ShouldBindJSON(&request); err != nil {
		slog.Error("Error binding JSON", err)
		pkg.PanicException(constant.InvalidRequest)
	}

	// Check if user exists and password is correct
	user, err := u.userRepository.GetUserByUsername(request.Username)
	if err != nil {
		slog.Error("Error getting user", err)
		pkg.PanicException(constant.InvalidCredentials)
	}

	// Check if password is correct
	credentialsError := user.CheckPassword(request.Password)
	if credentialsError != nil {
		slog.Error("Error checking password", credentialsError)
		pkg.PanicException(constant.InvalidCredentials)
	}

	// Generate token
	tokenString, err := auth.GenerateJWTToken(user.Username, user.Email, user.UserId)
	if err != nil {
		slog.Error("Error generating token", err)
		pkg.PanicException(constant.UnknownError)
	}

	ctx.JSON(200, pkg.BuildResponse(constant.Success, gin.H{"token": tokenString}))
}

func (u UserServiceImpl) RegisterUser(ctx *gin.Context) {
	defer pkg.PanicHandler(ctx)
	slog.Info("Received request to register user")

	var user dao.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		slog.Error("Error binding JSON", err)
		pkg.PanicException(constant.InvalidRequest)
	}

	// Check if user exists
	if u.userRepository.CheckIfUserExists(user.Username) {
		slog.Info("User already exists...", "user", user.Username)
		pkg.PanicException(constant.InvalidRequest)
	}

	// Hash password
	if err := user.HashPassword(); err != nil {
		slog.Error("Error hashing password", err)
		pkg.PanicException(constant.InvalidRequest)
	}

	// Register user and capture user_id
	userID, err := u.userRepository.RegisterUser(&user)
	if err != nil {
		slog.Error("Error registering user", err)
		pkg.PanicException(constant.UnknownError)
	}
	user.UserId = userID

	ctx.JSON(200, pkg.BuildResponse(constant.Success, user))
}
