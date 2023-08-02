package samples

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/kaka/pcr-backend/common/database"
	"gitlab.com/kaka/pcr-backend/common/models"
)

func GetSamples(ctx *gin.Context) {
	var samples []models.Sample

	sample_id := ctx.Param("sample_id")

	var params []interface{}

	query := `
		SELECT sample_id,samples.full_name,birthdate,sputalysed,comment,created_at,users.username
		FROM samples
		LEFT JOIN users ON samples.created_by = users.user_id
		WHERE 1 = 1
	`

	if sample_id != "" {
		query += `
			AND sample_id = $1
		`
		params = append(params, sample_id)
	}
	query += `
		AND created_at >= current_date - interval '14 day'
		ORDER BY created_at DESC, sample_id DESC;
	`

	rows, err := database.Instance.Query(query, params...)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	for rows.Next() {
		var sample models.Sample
		if err := rows.Scan(&sample.SampleId, &sample.FullName, &sample.Birthdate, &sample.Sputalysed, &sample.Comment, &sample.CreatedAt, &sample.CreatedBy); err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		samples = append(samples, sample)
	}

	// Empty array
	if len(samples) == 0 {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	ctx.JSON(http.StatusOK, &samples)
}
