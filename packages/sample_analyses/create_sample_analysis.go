package sampleanalyses

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/kaka/pcr-backend/common/database"
	"gitlab.com/kaka/pcr-backend/common/models"
	"gitlab.com/kaka/pcr-backend/packages/analyses"
	"gitlab.com/kaka/pcr-backend/packages/samples"
)

type AddAnalysisToSampleRequest struct {
	SampleID   string `json:"sample_id" binding:"required"`
	AnalysisID string `json:"analysis_id" binding:"required"`
	Run        string `json:"run" binding:"required"`
	Device     string `json:"device" binding:"required"`
}

func AddAnalysisToSample(ctx *gin.Context) {
	user_id := ctx.MustGet("user_id").(string)

	body := AddAnalysisToSampleRequest{}

	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if analysis exists
	if !analyses.AnalysisExistsByID(body.AnalysisID) {
		error_message := fmt.Sprintf("analysis %s does not exist", body.AnalysisID)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": error_message})
		return
	}

	// Check if sample exists
	if !samples.SampleExists(body.SampleID) {
		error_message := fmt.Sprintf("sample %s does not exist", body.SampleID)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": error_message})
		return
	}

	// Create sample analysis
	sample_analysis := models.SampleAnalysis{}
	sample_analysis.Sample = models.Sample{}
	sample_analysis.Analysis = models.Analysis{}
	sample_analysis.Sample.SampleID = body.SampleID
	sample_analysis.Analysis.AnalysisID = body.AnalysisID
	sample_analysis.Run = body.Run
	sample_analysis.Device = body.Device
	sample_analysis.Completed = false

	// Run query
	query := `
		WITH new_sample_analysis as 
			(INSERT INTO sampleanalyses (sample_id, analysis_id, run, device, created_by) 
			VALUES ($1, $2, $3, $4, $5)
			RETURNING sample_id, created_at, created_by, analysis_id)
			SELECT new_sample_analysis.created_at, users.username, analyses.analyt, analyses.material, analyses.assay, analyses.ready_mix, samples.firstname, samples.lastname
			FROM new_sample_analysis
			LEFT JOIN samples ON new_sample_analysis.sample_id = samples.sample_id
			LEFT JOIN users ON new_sample_analysis.created_by = users.user_id
			LEFT JOIN analyses ON new_sample_analysis.analysis_id = analyses.analysis_id
		`
	err := database.Instance.QueryRow(
		query,
		&sample_analysis.Sample.SampleID,
		&sample_analysis.Analysis.AnalysisID,
		&sample_analysis.Run,
		&sample_analysis.Device,
		user_id).Scan(
		&sample_analysis.CreatedAt,
		&sample_analysis.CreatedBy,
		&sample_analysis.Analysis.Analyt,
		&sample_analysis.Analysis.Material,
		&sample_analysis.Analysis.Assay,
		&sample_analysis.Analysis.ReadyMix,
		&sample_analysis.Sample.FirstName,
		&sample_analysis.Sample.LastName)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, &sample_analysis)
}
