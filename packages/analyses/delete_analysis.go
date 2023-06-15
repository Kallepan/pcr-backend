package analyses

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/kaka/pcr-backend/common/database"
)

func DeleteAnalysis(ctx *gin.Context) {
	analyt := ctx.Param("analyt")
	material := ctx.Param("material")
	assay := ctx.Param("assay")

	if !AnalysisExists(analyt, material, assay) {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "analysis not found"})
		return
	}

	query := "DELETE FROM analyses WHERE analyt = $1 AND material = $2 AND assay = $3"
	_, err := database.Instance.Exec(query, analyt, material, assay)

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
