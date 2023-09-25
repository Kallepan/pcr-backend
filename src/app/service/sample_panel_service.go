package service

import (
	"fmt"
	"log/slog"

	"github.com/gin-gonic/gin"
	"gitlab.com/kallepan/pcr-backend/app/constant"
	"gitlab.com/kallepan/pcr-backend/app/domain/dco"
	"gitlab.com/kallepan/pcr-backend/app/pkg"
	"gitlab.com/kallepan/pcr-backend/app/repository"
)

type SamplePanelService interface {
	GetSamplePanels(ctx *gin.Context)
	ResetSamplePanel(ctx *gin.Context)
	UpdateSamplePanel(ctx *gin.Context)
	CreateSamplePanel(ctx *gin.Context)

	GetStatistics(ctx *gin.Context)
}

type SamplePanelServiceImpl struct {
	sp repository.SamplePanelRepository
	p  repository.PanelRepository
	s  repository.SampleRepository
}

func SamplePanelServiceInit(samplePanelRepo repository.SamplePanelRepository, sampleRepo repository.SampleRepository, panelRepo repository.PanelRepository) *SamplePanelServiceImpl {
	return &SamplePanelServiceImpl{
		sp: samplePanelRepo,
		p:  panelRepo,
		s:  sampleRepo,
	}
}

func (s SamplePanelServiceImpl) GetStatistics(ctx *gin.Context) {
	defer pkg.PanicHandler(ctx)
	slog.Info("Received request to get statistics")

	// Get statistics
	statistics, err := s.sp.GetStatistics()
	if err != nil {
		slog.Error("Error getting statistics", err)
		pkg.PanicException(constant.UnknownError)
	}

	ctx.JSON(200, pkg.BuildResponse(constant.Success, statistics))
}

func (s SamplePanelServiceImpl) CreateSamplePanel(ctx *gin.Context) {
	defer pkg.PanicHandler(ctx)
	slog.Info("Received request to create sample panel")

	var addSamplePanelRequest dco.AddSamplePanelRequest
	if err := ctx.ShouldBindJSON(&addSamplePanelRequest); err != nil {
		slog.Error("Error binding JSON", err)
		pkg.PanicException(constant.InvalidRequest)
	}

	// Check if sample exists
	if !s.s.SampleExists(addSamplePanelRequest.SampleID) {
		errorMessage := fmt.Sprintf("Sample %s does not exist", addSamplePanelRequest.SampleID)
		slog.Error(errorMessage)
		pkg.PanicExceptionWithMessage(constant.InvalidRequest, errorMessage)
	}

	// Check if panel exists
	if !s.p.PanelExists(addSamplePanelRequest.PanelID) {
		errorMessage := fmt.Sprintf("Panel %s does not exist", addSamplePanelRequest.PanelID)
		slog.Error(errorMessage)
		pkg.PanicExceptionWithMessage(constant.InvalidRequest, errorMessage)
	}

	// Validate request
	if s.sp.SamplePanelExists(addSamplePanelRequest.SampleID, addSamplePanelRequest.PanelID) {
		// Set deleted to false
		if err := s.sp.UndeleteSamplePanel(addSamplePanelRequest.SampleID, addSamplePanelRequest.PanelID); err != nil {
			slog.Error("Error undeleting sample panel", err)
			pkg.PanicException(constant.UnknownError)
		}
		ctx.JSON(200, pkg.BuildResponse(constant.Success, pkg.Null()))
		return
	}

	// Create sample panel
	samplePanel, err := s.sp.CreateSamplePanel(addSamplePanelRequest, ctx.GetString("user_id"))
	if err != nil {
		slog.Error("Error creating sample panel", err)
		pkg.PanicException(constant.UnknownError)
	}

	ctx.JSON(200, pkg.BuildResponse(constant.Success, samplePanel))
}

func (s SamplePanelServiceImpl) UpdateSamplePanel(ctx *gin.Context) {
	defer pkg.PanicHandler(ctx)
	slog.Info("Received request to update sample panel")

	// Get values from path params
	sampleID := ctx.Param("sample_id")
	panelID := ctx.Param("panel_id")
	if sampleID == "" || panelID == "" {
		errorMessage := fmt.Sprintf("Sample: %s with Panel %s is invalid", sampleID, panelID)
		slog.Error(errorMessage)
		pkg.PanicExceptionWithMessage(constant.InvalidRequest, errorMessage)
	}

	var updateSamplePanelRequest dco.UpdateSamplePanelRequest
	if err := ctx.ShouldBindJSON(&updateSamplePanelRequest); err != nil {
		slog.Error("Error binding JSON", err)
		pkg.PanicException(constant.InvalidRequest)
	}

	// Validate request
	if !s.sp.SamplePanelExists(sampleID, panelID) {
		errorMessage := fmt.Sprintf("Sample %s with panel %s does not exist", sampleID, panelID)
		slog.Error(errorMessage)
		pkg.PanicExceptionWithMessage(constant.InvalidRequest, errorMessage)
	}

	// Update sample panel
	if err := s.sp.UpdateSamplePanel(sampleID, panelID, updateSamplePanelRequest); err != nil {
		slog.Error("Error updating sample panel", err)
		pkg.PanicException(constant.UnknownError)
	}

	ctx.JSON(200, pkg.BuildResponse(constant.Success, pkg.Null()))
}

func (s SamplePanelServiceImpl) ResetSamplePanel(ctx *gin.Context) {
	defer pkg.PanicHandler(ctx)
	slog.Info("Received request to reset sample panel")

	var resetSamplePanelRequest dco.ResetSamplePanelRequest
	if err := ctx.ShouldBindJSON(&resetSamplePanelRequest); err != nil {
		slog.Error("Error binding JSON", err)
		pkg.PanicException(constant.InvalidRequest)
	}

	// Validate request
	if !s.sp.SamplePanelExists(resetSamplePanelRequest.SampleID, resetSamplePanelRequest.PanelID) {
		errorMessage := fmt.Sprintf("Sample %s with panel %s does not exist", resetSamplePanelRequest.SampleID, resetSamplePanelRequest.PanelID)
		slog.Error(errorMessage)
		pkg.PanicExceptionWithMessage(constant.InvalidRequest, errorMessage)
	}

	// Reset sample panel
	if err := s.sp.ResetSamplePanel(resetSamplePanelRequest); err != nil {
		slog.Error("Error resetting sample panel", err)
		pkg.PanicException(constant.UnknownError)
	}

	ctx.JSON(200, pkg.BuildResponse(constant.Success, pkg.Null()))
}

func (s SamplePanelServiceImpl) GetSamplePanels(ctx *gin.Context) {
	defer pkg.PanicHandler(ctx)
	slog.Info("Received request to get sample panels")
	// Get values from query params
	sampleID := ctx.Query("sample_id")
	runDate := ctx.Query("run_date")
	device := ctx.Query("device")
	run := ctx.Query("run")

	// Get sample panels
	samplePanels, err := s.sp.GetSamplePanels(sampleID, runDate, device, run)
	if err != nil {
		slog.Error("Error getting sample panels", err)
		pkg.PanicException(constant.UnknownError)
	}

	ctx.JSON(200, pkg.BuildResponse(constant.Success, samplePanels))
}
