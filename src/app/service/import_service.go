package service

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"gitlab.com/kallepan/pcr-backend/app/constant"
	"gitlab.com/kallepan/pcr-backend/app/domain/dco"
	"gitlab.com/kallepan/pcr-backend/app/pkg"
	"gitlab.com/kallepan/pcr-backend/app/repository"
)

type ImportService interface {
	ImportSamplePanel(ctx *gin.Context)
}

type ImportServiceImpl struct {
	importRepository repository.ImportRepository
}

func ImportServiceInit(importRepository repository.ImportRepository) *ImportServiceImpl {
	return &ImportServiceImpl{
		importRepository: importRepository,
	}
}

func (i ImportServiceImpl) ImportSamplePanel(ctx *gin.Context) {
	/* Import sample panel from  */
	defer pkg.PanicHandler(ctx)
	slog.Info("Received request to import sample panel")

	var samplePanelRequest dco.SamplePanelRequest
	if err := ctx.ShouldBindJSON(&samplePanelRequest); err != nil {
		slog.Error("Error binding JSON", err)
		pkg.PanicException(constant.InvalidRequest)
	}

	// Validate request
	for _, samplePanel := range samplePanelRequest.SamplePanel {
		if err := samplePanel.Validate(); err != nil {
			slog.Error("Error validating sample panel", err)
			pkg.PanicException(constant.InvalidRequest)
		}
		if !i.importRepository.PanelExists(samplePanel.PanelID) {
			slog.Error("Panel does not exist")
			pkg.PanicException(constant.InvalidRequest)
		}
	}

	// Save sample panel
	if err := i.importRepository.Save(samplePanelRequest.SamplePanel); err != nil {
		slog.Error("Error saving sample panel", err)
		pkg.PanicException(constant.UnknownError)
	}

	ctx.JSON(200, pkg.BuildResponse(constant.Success, pkg.Null()))
}
