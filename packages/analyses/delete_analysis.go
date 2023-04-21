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

	query := "DELETE FROM analyses WHERE analyt = $1 AND material = $2 AND assay = $3"
	_, err := database.Instance.Exec(query, analyt, material, assay)

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
