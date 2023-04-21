package samples

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/kaka/pcr-backend/common/database"
)

func DeleteSample(ctx *gin.Context) {
	tagesnummer := ctx.Param("tagesnummer")

	query := `DELETE FROM samples WHERE tagesnummer = $1`

	_, err := database.Instance.Exec(query, tagesnummer)

	switch err {
	case nil:
		break
	case sql.ErrNoRows:
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "sample not found"})
		return
	default:
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}
