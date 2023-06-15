package samplesanalyses

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/kaka/pcr-backend/common/database"
)

func DeleteSampleAnalysis(ctx *gin.Context) {
	sample_id := ctx.Param("sample_id")
	analysis_id := ctx.Param("analysis_id")

	query := `DELETE FROM samplesanalyses WHERE sample_id = $1 AND analysis_id = $2`

	_, err := database.Instance.Exec(query, sample_id, analysis_id)

	switch err {
	case nil:
		break
	case sql.ErrNoRows:
		error_message := fmt.Sprintf("Analyse %s ist nicht mit Probe %s verknüpft", analysis_id, sample_id)
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": error_message})
		return
	default:
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}
