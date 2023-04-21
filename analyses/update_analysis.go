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
	analyt := ctx.Param("analyt")
	material := ctx.Param("material")
	assay := ctx.Param("assay")

	var request UpdateAnalysisRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := "UPDATE analyses SET ready_mix = $1 WHERE analyt = $2 AND material = $3 AND assay = $4"
	_, err := database.Instance.Exec(query, request.ReadyMix, analyt, material, assay)

	switch err {
	case nil:
		break
	case sql.ErrNoRows:
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "analysis not found"})
		return
	default:
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}
