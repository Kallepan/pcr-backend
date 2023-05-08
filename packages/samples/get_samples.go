package samples

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/kaka/pcr-backend/common/database"
	"gitlab.com/kaka/pcr-backend/common/models"
)

func GetSamples(ctx *gin.Context) {
	var samples []models.Sample

	query := `
		SELECT sample_id,samples.firstname,samples.lastname,sputalysed,comment,created_at,users.username
		FROM samples
		LEFT JOIN users ON samples.created_by = users.user_id
		WHERE created_at >= current_date - interval '10 day'
		ORDER BY $1 DESC LIMIT $2;
		`

	// TODO: Add pagination and filtering by params
	rows, err := database.Instance.Query(query, "created_at", 100)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	for rows.Next() {
		var sample models.Sample
		if err := rows.Scan(&sample.SampleID, &sample.FirstName, &sample.LastName, &sample.Sputalysed, &sample.Comment, &sample.CreatedAt, &sample.CreatedBy); err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		samples = append(samples, sample)
	}

	ctx.JSON(http.StatusOK, &samples)
}
