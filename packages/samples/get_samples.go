package samples

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/kaka/pcr-backend/common/database"
	"gitlab.com/kaka/pcr-backend/common/models"
)

func FetchSampleInformationFromDatabase(sampleID string) (*models.Sample, error) {
	var sample models.Sample

	query :=
		`SELECT sample_id,samples.full_name,created_at,users.username,birthdate,sputalysed,comment
		FROM samples 
		LEFT JOIN users ON samples.created_by = users.user_id
		WHERE sample_id = $1;`

	row := database.Instance.QueryRow(query, sampleID)

	if err := row.Scan(&sample.SampleId, &sample.FullName, &sample.CreatedAt, &sample.CreatedBy, &sample.Birthdate, &sample.Sputalysed, &sample.Comment); err != nil {
		return nil, err
	}

	sample.Birthdate = sample.Birthdate[:10]

	return &sample, nil
}

func GetSamples(ctx *gin.Context) {
	var samples []models.Sample

	sample_id := ctx.Query("sample_id")

	var params []interface{}

	query := `
		SELECT sample_id,samples.full_name,birthdate,sputalysed,comment,created_at,users.username
		FROM samples
		LEFT JOIN users ON samples.created_by = users.user_id
		WHERE 1 = 1 AND
		`

	if sample_id != "" {
		query += `
			sample_id = ?
		`
		params = append(params, sample_id)
	} else {
		query += `
		created_at >= current_date - interval '14 day'
		ORDER BY created_at DESC, sample_id DESC;
		`
	}

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
		sample.Birthdate = sample.Birthdate[:10]
		samples = append(samples, sample)
	}

	ctx.JSON(http.StatusOK, &samples)
}
