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

type SampleService interface {
	AddSample(ctx *gin.Context)
	UpdateSample(ctx *gin.Context)
	DeleteSample(ctx *gin.Context)
	GetSamples(ctx *gin.Context)
}

type SampleServiceImpl struct {
	repo repository.SampleRepository
}

func SampleServiceInit(repo repository.SampleRepository) *SampleServiceImpl {
	return &SampleServiceImpl{
		repo: repo,
	}
}

func (s SampleServiceImpl) GetSamples(ctx *gin.Context) {
	/* Returns all samples with an optional filter for sample_id */
	defer pkg.PanicHandler(ctx)
	slog.Info("Received request to get samples")

	// Get sample_id from query param
	sampleID := ctx.Query("sample_id")

	// Get samples
	samples, err := s.repo.GetSamples(sampleID)
	if err != nil {
		slog.Error("Error getting samples", err)
		pkg.PanicException(constant.UnknownError)
	}

	ctx.JSON(200, pkg.BuildResponse(constant.Success, samples))
}

func (s SampleServiceImpl) AddSample(ctx *gin.Context) {
	/* Add a new sample */
	defer pkg.PanicHandler(ctx)
	slog.Info("Received request to add sample")

	var sampleRequest dco.AddSampleRequest
	if err := ctx.ShouldBindJSON(&sampleRequest); err != nil {
		slog.Error("Error binding JSON", err)
		pkg.PanicException(constant.InvalidRequest)
	}
	// Validate request
	if err := sampleRequest.Validate(); err != nil {
		slog.Error("Error validating sample", err)
		pkg.PanicException(constant.InvalidRequest)
	}

	userID := ctx.MustGet("user_id").(string)

	// Check if sample already exists
	if s.repo.SampleExists(sampleRequest.SampleId) {
		errorMessage := fmt.Sprintf("Sample %s already exists", sampleRequest.SampleId)
		slog.Error(errorMessage)
		pkg.PanicExceptionWithMessage(constant.InvalidRequest, errorMessage)
	}

	// Create sample
	sample, err := s.repo.CreateSample(sampleRequest, userID)
	if err != nil {
		slog.Error("Error creating sample", err)
		pkg.PanicException(constant.UnknownError)
	}

	ctx.JSON(200, pkg.BuildResponse(constant.Success, sample))
}

func (s SampleServiceImpl) UpdateSample(ctx *gin.Context) {
	/* Update sample */
	defer pkg.PanicHandler(ctx)
	slog.Info("Received request to update sample")

	sampleID := ctx.Param("sample_id")
	var sampleRequest dco.UpdateSampleRequest

	if err := ctx.ShouldBindJSON(&sampleRequest); err != nil {
		slog.Error("Error binding JSON", err)
		pkg.PanicException(constant.InvalidRequest)
	}

	// Validate request
	if err := sampleRequest.Validate(); err != nil {
		errorMessage := fmt.Sprintf("Error validating sample: %s", err)
		slog.Error(errorMessage)
		pkg.PanicExceptionWithMessage(constant.InvalidRequest, errorMessage)
	}

	// Check if sample exists
	if !s.repo.SampleExists(sampleID) {
		errorMessage := fmt.Sprintf("Sample %s does not exist", sampleID)
		slog.Error(errorMessage)
		pkg.PanicExceptionWithMessage(constant.InvalidRequest, errorMessage)
	}

	// Update sample
	sample, err := s.repo.UpdateSample(sampleRequest, sampleID)
	if err != nil {
		slog.Error("Error updating sample", err)
		pkg.PanicException(constant.UnknownError)
	}

	ctx.JSON(200, pkg.BuildResponse(constant.Success, sample))
}

func (s SampleServiceImpl) DeleteSample(ctx *gin.Context) {
	/* Delete sample */
	defer pkg.PanicHandler(ctx)
	slog.Info("Received request to delete sample")

	sampleID := ctx.Param("sample_id")

	// Check if sample exists
	if !s.repo.SampleExists(sampleID) {
		errorMessage := fmt.Sprintf("Sample %s does not exist", sampleID)
		slog.Error(errorMessage)
		pkg.PanicExceptionWithMessage(constant.InvalidRequest, errorMessage)
	}

	// Delete sample
	if err := s.repo.DeleteSample(sampleID); err != nil {
		slog.Error("Error deleting sample", err)
		pkg.PanicException(constant.UnknownError)
	}

	ctx.JSON(200, pkg.BuildResponse(constant.Success, pkg.Null()))
}
