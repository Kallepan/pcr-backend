package controller

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/kallepan/pcr-backend/app/service"
)

type PanelController interface {
	GetPanels(ctx *gin.Context)
}

type PanelControllerImpl struct {
	svc service.PanelService
}

func PanelControllerInit(panelService service.PanelService) *PanelControllerImpl {
	return &PanelControllerImpl{
		svc: panelService,
	}
}

func (p PanelControllerImpl) GetPanels(ctx *gin.Context) {
	p.svc.GetPanels(ctx)
}
