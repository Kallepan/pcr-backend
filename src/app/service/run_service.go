package service

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"gitlab.com/kallepan/pcr-backend/app/constant"
	"gitlab.com/kallepan/pcr-backend/app/domain/dco"
	"gitlab.com/kallepan/pcr-backend/app/pkg"
	"gitlab.com/kallepan/pcr-backend/app/repository"
)

type RunService interface {
	CreateRun(ctx *gin.Context)
}

type RunServiceImpl struct {
	runRepo    repository.RunRepository
	sampleRepo repository.SampleRepository
	panelRepo  repository.PanelRepository
}

func RunServiceInit(runRepo repository.RunRepository, sampleRepo repository.SampleRepository, panelRepo repository.PanelRepository) *RunServiceImpl {
	return &RunServiceImpl{
		runRepo:    runRepo,
		sampleRepo: sampleRepo,
		panelRepo:  panelRepo,
	}
}

func (s RunServiceImpl) CreateRun(ctx *gin.Context) {
	/* Creates a new run */
	defer pkg.PanicHandler(ctx)
	slog.Info("Received request to create run")

	// Bind request
	var request dco.CreatRunRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		slog.Error("Error binding JSON", err)
		pkg.PanicException(constant.InvalidRequest)
	}

	// Validate request by looping through elements
	exportDataArr := []dco.ExportData{}
	for _, element := range request.Elements {
		// Fill export data
		var exportData dco.ExportData
		if element.SampleID != nil && element.PanelID != nil {

			// Check if sample exists
			if !s.sampleRepo.SampleExists(*element.SampleID) {
				slog.Error("Error checking if sample exists")
				pkg.PanicException(constant.UnknownError)
			}

			// Check if panel exists
			if !s.panelRepo.PanelExists(*element.PanelID) {
				slog.Error("Error checking if panel exists")
				pkg.PanicException(constant.UnknownError)
			}

			// Check if already in run
			if s.runRepo.IsAlreadyInRun(*element.SampleID, *element.PanelID) {
				slog.Error("Error checking if sample panel is already in run")
				pkg.PanicException(constant.UnknownError)
			}

			// Get sample
			sample, err := s.sampleRepo.GetSample(*element.SampleID)
			if err != nil {
				slog.Error("Error getting sample", err)
				pkg.PanicException(constant.InvalidRequest)
			}

			// Get panel
			panel, err := s.panelRepo.GetPanel(*element.PanelID)
			if err != nil {
				slog.Error("Error getting panel", err)
				pkg.PanicException(constant.InvalidRequest)
			}

			// Fetch last run id
			// This is the last "run id" in which the sample was used
			lastRunId, err := s.runRepo.GetLastRunId(*element.SampleID)
			if err != nil {
				slog.Error("Error getting last run id", err)
				pkg.PanicException(constant.UnknownError)
			}
			exportData.LastRunId = lastRunId
			// Fill export data
			exportData.Sample = sample
			exportData.Panel = panel
			exportData.IsControl = false

			// Append to export data array
			exportDataArr = append(exportDataArr, exportData)
		} else if element.ControlID != nil && element.Description != nil {
			// Handle Controls
			exportData.Description = element.Description
			exportData.IsControl = true
			exportDataArr = append(exportDataArr, exportData)
		} else {
			// Invalid Data
			slog.Error("Invalid data")
			pkg.PanicException(constant.InvalidRequest)
		}
	}

	// Create run
	pathToFile, err := s.runRepo.CreateRun(exportDataArr, request.Device, request.Run, request.Date)
	if err != nil {
		slog.Error("Error creating run")
		pkg.PanicException(constant.UnknownError)
	}

	ctx.File(pathToFile)

}
