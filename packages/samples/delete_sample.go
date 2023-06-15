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

	if !SampleExists(sample_id) {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "sample not found"})
		return
	}

	query := `DELETE FROM samples WHERE sample_id = $1`

	_, err := database.Instance.Exec(query, sample_id)

	switch err {
	case nil:
		break
	case sql.ErrNoRows:
		error_message := fmt.Sprintf("sample with id %s not found", sample_id)
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": error_message})
		return
	default:
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}
