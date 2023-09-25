package service

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"gitlab.com/kallepan/pcr-backend/app/constant"
	"gitlab.com/kallepan/pcr-backend/app/pkg"
	"gitlab.com/kallepan/pcr-backend/app/repository"
)

type PanelService interface {
	GetPanels(ctx *gin.Context)
}

type PanelServiceImpl struct {
	panelRepository repository.PanelRepository
}

func PanelServiceInit(panelRepository repository.PanelRepository) *PanelServiceImpl {
	return &PanelServiceImpl{
		panelRepository: panelRepository,
	}
}

func (p PanelServiceImpl) GetPanels(ctx *gin.Context) {
	defer pkg.PanicHandler(ctx)
	slog.Info("Received request to get panels")

	// Get panel_id from query
	panelID := ctx.Query("panel_id")

	// Get panels
	panels, err := p.panelRepository.GetPanels(panelID)
	if err != nil {
		slog.Error("Error getting panels", err)
		pkg.PanicException(constant.UnknownError)
	}

	ctx.JSON(200, pkg.BuildResponse(constant.Success, panels))
}
