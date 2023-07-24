package samplespanels

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/kaka/pcr-backend/common/database"
)

type ResetSamplePanelRequest struct {
	PanelId  string `json:"panel_id" binding:"required"`
	SampleId string `json:"sample_id" binding:"required"`
}

func ResetSamplePanel(ctx *gin.Context) {
	var request ResetSamplePanelRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.AbortWithStatusJSON(400, gin.H{"message": err.Error()})
		return
	}

	query := `
		UPDATE samplespanels
		SET run = NULL, device = NULL, position = NULL, run_date = NULL
		WHERE
			sample_id = $1 AND
			panel_id = $2 AND
			deleted = false
	`

	if _, err := database.Instance.Exec(query, request.SampleId, request.PanelId); err != nil {
		ctx.AbortWithStatusJSON(500, gin.H{"message": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}
