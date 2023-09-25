package controller

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/kallepan/pcr-backend/app/service"
)

type ImportController interface {
	ImportSample(ctx *gin.Context)
}

type ImportControllerImpl struct {
	svc service.ImportService
}

func ImportControllerInit(ImportService service.ImportService) *ImportControllerImpl {
	return &ImportControllerImpl{
		svc: ImportService,
	}
}

func (ctrl ImportControllerImpl) ImportSample(ctx *gin.Context) {
	ctrl.svc.ImportSamplePanel(ctx)
}
