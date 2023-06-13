package samplesanalyses

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/kaka/pcr-backend/common/database"
)

type UpdateSampleAnalysisRequest struct {
	Completed *bool `json:"completed" binding:"required"`
}

func UpdateSampleAnalysis(ctx *gin.Context) {
	sample_id := ctx.Param("sample_id")
	analysis_id := ctx.Param("analysis_id")

	body := UpdateSampleAnalysisRequest{}

	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// Check if sample_analysis exists
	if !SampleAnalysisExists(sample_id, analysis_id) {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "sample analysis not found"})
		return
	}

	if *body.Completed {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "cannot set completed to true"})
		return
	}

	// Run query
	query := `
		UPDATE samplesanalyses
		SET run_id = NULL, device_id = NULL, position = NULL
		WHERE sample_id = $1 AND analysis_id = $2;
	`
	_, err := database.Instance.Exec(query, sample_id, analysis_id)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	ctx.Status(http.StatusOK)
}
