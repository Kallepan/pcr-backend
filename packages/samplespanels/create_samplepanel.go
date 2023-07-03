package samplespanels

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/kaka/pcr-backend/common/database"
	"gitlab.com/kaka/pcr-backend/common/models"
	"gitlab.com/kaka/pcr-backend/packages/panels"
	"gitlab.com/kaka/pcr-backend/packages/samples"
)

type AddAnalysisToSampleRequest struct {
	SampleId string `json:"sample_id" binding:"required"`
	PanelId  string `json:"analysis_id" binding:"required"`
}

func AddAnalysisToSample(ctx *gin.Context) {
	user_id := ctx.MustGet("user_id").(string)

	body := AddAnalysisToSampleRequest{}

	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// Check if analysis exists
	if !panels.PanelExists(body.PanelId) {
		error_message := fmt.Sprintf("panel %s does not exist", body.PanelId)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": error_message})
		return
	}

	// Check if sample exists
	if !samples.SampleExists(body.SampleId) {
		error_message := fmt.Sprintf("sample %s does not exist", body.SampleId)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": error_message})
		return
	}

	// Check if sample analysis already exists
	if SampleAnalysisExists(body.SampleId, body.PanelId) {
		// Sample already exists, update deleted status and return 200
		err := UpdateSampleAnalysisDeletedStatus(body.SampleId, body.PanelId, false)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		}
		ctx.Status(http.StatusOK)
		return
	}

	// Create sample analysis
	sample_analysis := models.SampleAnalysis{}
	sample_analysis.Sample = models.Sample{}
	sample_analysis.Panel = models.Panel{}
	sample_analysis.Sample.SampleId = body.SampleId
	sample_analysis.Panel.PanelId = body.PanelId

	// Run query
	query := `
		WITH sample_query AS (
			SELECT samples.sample_id, samples.full_name, samples.created_at, samples.created_by, samples.sputalysed, users.username AS created_by
			FROM samples
			LEFT JOIN users ON samples.created_by = users.user_id
			WHERE samples.sample_id = $1
			GROUP BY samples.sample_id, samples.full_name, samples.created_at, samples.sputalysed, users.username
		), sample_panel AS ( 
			INSERT INTO samplespanels (sample_id, panel_id, created_by) 
			VALUES ($1, $2, $3)
			RETURNING sample_id, panel_id, created_at, created_by
		)
			SELECT new_sample_analysis.created_at, users.username, panels.panel_id, panels.ready_mix, sample_query.full_name, sample_query.sputalysed, sample_query.created_at, sample_query.created_by
			FROM new_sample_analysis
			LEFT JOIN sample_query ON new_sample_analysis.sample_id = sample_query.sample_id
			LEFT JOIN users ON new_sample_analysis.created_by = users.user_id
			LEFT JOIN analyses ON new_sample_analysis.panel_id = panels.panel_id
		`
	err := database.Instance.QueryRow(
		query,
		&sample_analysis.Sample.SampleId,
		&sample_analysis.Panel.PanelId,
		user_id).Scan(
		&sample_analysis.CreatedAt,
		&sample_analysis.CreatedBy,
		&sample_analysis.Panel.PanelId,
		&sample_analysis.Panel.ReadyMix,
		&sample_analysis.Sample.FullName,
		&sample_analysis.Sample.Sputalysed,
		&sample_analysis.Sample.CreatedAt,
		&sample_analysis.Sample.CreatedBy)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, &sample_analysis)
}
