package samples

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/kaka/pcr-backend/common/database"
	"gitlab.com/kaka/pcr-backend/common/models"
)

func GetSample(ctx *gin.Context) {
	sample_id := ctx.Param("sample_id")

	var sample models.Sample

	query :=
		`SELECT sample_id,samples.firstname,samples.lastname,created_at,users.username,sputalysed,comment
		FROM samples 
		LEFT JOIN users ON samples.created_by = users.user_id
		WHERE sample_id = $1;`

	row := database.Instance.QueryRow(query, sample_id)

	switch err := row.Scan(&sample.SampleID, &sample.FirstName, &sample.LastName, &sample.CreatedAt, &sample.CreatedBy, &sample.Sputalysed, &sample.Comment); err {
	case sql.ErrNoRows:
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "sample not found"})
		return
	case nil:
		break
	default:
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, sample)
}
