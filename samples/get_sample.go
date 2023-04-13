package samples

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/kaka/pcr-backend/common/database"
	"gitlab.com/kaka/pcr-backend/common/models"
)

func GetSample(ctx *gin.Context) {
	tagesnummer := ctx.Param("tagesnummer")

	var sample models.Sample

	query := `SELECT tagesnummer,name,created_at,created_by FROM samples WHERE tagesnummer = $1`
	row := database.Instance.QueryRow(query, tagesnummer)

	switch err := row.Scan(&sample.Tagesnummer, &sample.Name, &sample.CreatedAt, &sample.CreatedBy); err {
	case sql.ErrNoRows:
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "sample not found"})
	case nil:
		ctx.JSON(http.StatusOK, &sample)
	default:
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}
