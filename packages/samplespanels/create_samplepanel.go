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
	PanelId  string `json:"panel_id" binding:"required"`
}

func buildCreateQuery() string {
	query := `
		WITH sample_query AS (
			SELECT samples.sample_id, samples.full_name, samples.created_at, samples.sputalysed, users.username AS created_by
			FROM samples
			LEFT JOIN users ON samples.created_by = users.user_id
			WHERE samples.sample_id = $1
			GROUP BY samples.sample_id, samples.full_name, samples.created_at, samples.sputalysed, users.username
		), sample_panel AS ( 
			INSERT INTO samplespanels (sample_id, panel_id, created_by) 
			VALUES ($1, $2, $3)
			RETURNING sample_id, panel_id, created_at, created_by
		)
			SELECT sample_panel.created_at, users.username, panels.panel_id, panels.ready_mix, sample_query.full_name, sample_query.sputalysed, sample_query.created_at, sample_query.created_by
			FROM sample_panel
			LEFT JOIN sample_query ON sample_panel.sample_id = sample_query.sample_id
			LEFT JOIN users ON sample_panel.created_by = users.user_id
			LEFT JOIN panels ON sample_panel.panel_id = panels.panel_id
		`
	return query
}

func executeCreateQuery(query string, userID string, sampleAnalysis *models.SampleAnalysis) error {
	err := database.Instance.QueryRow(
		query,
		&sampleAnalysis.Sample.SampleId,
		&sampleAnalysis.Panel.PanelId,
		userID).Scan(
		&sampleAnalysis.CreatedAt,
		&sampleAnalysis.CreatedBy,
		&sampleAnalysis.Panel.PanelId,
		&sampleAnalysis.Panel.ReadyMix,
		&sampleAnalysis.Sample.FullName,
		&sampleAnalysis.Sample.Sputalysed,
		&sampleAnalysis.Sample.CreatedAt,
		&sampleAnalysis.Sample.CreatedBy)

	return err
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
	if SamplePanelExists(body.SampleId, body.PanelId) {
		// Sample already exists, update deleted status and return 200
		err := UpdateSamplePanelDeletedStatus(body.SampleId, body.PanelId, false)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		}
		ctx.Status(http.StatusOK)
		return
	}

	// Create sample analysis
	sampleAnalysis := models.SampleAnalysis{}
	sampleAnalysis.Sample = models.Sample{}
	sampleAnalysis.Panel = models.Panel{}
	sampleAnalysis.Sample.SampleId = body.SampleId
	sampleAnalysis.Panel.PanelId = body.PanelId

	// Run query
	query := buildCreateQuery()

	// Execute query
	err := executeCreateQuery(query, user_id, &sampleAnalysis)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, &sampleAnalysis)
}
