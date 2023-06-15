package analyses

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/kaka/pcr-backend/common/database"
)

type UpdateAnalysisRequest struct {
	ReadyMix bool `json:"ready_mix"`
}

func UpdateAnalysis(ctx *gin.Context) {
	analysis_id := ctx.Param("analysis_id")

	var request UpdateAnalysisRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	query := "UPDATE analyses SET ready_mix = $1 WHERE analysis_id = $2"
	_, err := database.Instance.Exec(query, request.ReadyMix, analysis_id)

	switch err {
	case nil:
		break
	case sql.ErrNoRows:
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "analysis not found"})
		return
	default:
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}
