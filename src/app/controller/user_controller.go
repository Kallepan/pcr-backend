package controller

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/kallepan/pcr-backend/app/service"
)

type UserController interface {
	RegisterUser(ctx *gin.Context)
	LoginUser(ctx *gin.Context)
}

type UserControllerImpl struct {
	svc service.UserService
}

func UserControllerInit(userService service.UserService) *UserControllerImpl {
	return &UserControllerImpl{
		svc: userService,
	}
}

func (u UserControllerImpl) LoginUser(ctx *gin.Context) {
	u.svc.LoginUser(ctx)
}

func (u UserControllerImpl) RegisterUser(ctx *gin.Context) {
	u.svc.RegisterUser(ctx)
}
