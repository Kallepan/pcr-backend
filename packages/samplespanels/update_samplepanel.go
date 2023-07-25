package samplespanels

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/kaka/pcr-backend/common/database"
)

type UpdateSampleAnalysisRequest struct {
	Deleted *bool `json:"deleted" binding:"required"`
}

func UpdateSamplePanel(ctx *gin.Context) {
	sample_id := ctx.Param("sample_id")
	panel_id := ctx.Param("panel_id")
	body := UpdateSampleAnalysisRequest{}

	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// Check if sample_analysis exists
	if !SamplePanelExists(sample_id, panel_id) {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("SamplePanel with sample_id %s and panel_id %s not found", sample_id, panel_id)})
		return
	}

	// Run query
	query := `
		UPDATE samplespanels
		SET deleted = $3
		WHERE sample_id = $1 AND panel_id = $2;
	`
	_, err := database.Instance.Exec(query, sample_id, panel_id, *body.Deleted)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}

func UpdateSamplePanelDeletedStatus(sample_id string, analysis_id string, deleted bool) error {
	query := `
		UPDATE samplespanels
		SET deleted = $3
		WHERE sample_id = $1 AND panel_id = $2;
	`
	_, err := database.Instance.Exec(query, sample_id, analysis_id, deleted)

	return err
}
