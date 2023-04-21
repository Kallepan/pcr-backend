package samples

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/kaka/pcr-backend/common/database"
)

func DeleteSample(ctx *gin.Context) {
	sample_id := ctx.Param("sample_id")

	query := `DELETE FROM samples WHERE sample_id = $1`

	_, err := database.Instance.Exec(query, sample_id)

	switch err {
	case nil:
		break
	case sql.ErrNoRows:
		error_message := fmt.Sprintf("sample with id %s not found", sample_id)
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": error_message})
		return
	default:
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}
