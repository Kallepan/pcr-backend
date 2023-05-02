package samplesanalyses

import (
	"database/sql"
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

	// Run query
	query := `UPDATE samplesanalyses SET completed = $1 WHERE sample_id = $2 AND analysis_id = $3`

	_, err := database.Instance.Exec(query, body.Completed, sample_id, analysis_id)

	switch err {
	case nil:
		break
	case sql.ErrNoRows:
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "sample analysis not found"})
		return
	default:
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}
