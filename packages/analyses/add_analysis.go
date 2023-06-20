package analyses

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/kaka/pcr-backend/common/database"
	"gitlab.com/kaka/pcr-backend/common/models"
)

type AddAnalysisRequest struct {
	AnalysisId string `json:"analysis_id" binding:"required"`
	ReadyMix   *bool  `json:"ready_mix"`
}

func AddAnalysis(ctx *gin.Context) {
	var request AddAnalysisRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	analysis := models.Analysis{
		AnalysisId: request.AnalysisId,
		ReadyMix:   *request.ReadyMix,
	}

	// Check if analysis already exists
	if AnalysisExists(analysis.AnalysisId) {
		ctx.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": "analysis already exists"})
		return
	}

	// Insert analysis
	query := `
			INSERT INTO analyses (analysis_id, ready_mix) 
			VALUES ($1,$2) 
			RETURNING analysis_id;`
	err := database.Instance.QueryRow(query, analysis.AnalysisId, analysis.ReadyMix).Scan(&analysis.AnalysisId)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, analysis)
}
