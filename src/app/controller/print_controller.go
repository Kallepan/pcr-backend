package controller

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/kallepan/pcr-backend/app/service"
)

type PrintController interface {
	PrintSample(ctx *gin.Context)
}

type PrintControllerImpl struct {
	svc service.PrintService
}

func PrintControllerInit(svc service.PrintService) *PrintControllerImpl {
	return &PrintControllerImpl{
		svc: svc,
	}
}

func (p PrintControllerImpl) PrintSample(ctx *gin.Context) {
	p.svc.PrintSample(ctx)
}
