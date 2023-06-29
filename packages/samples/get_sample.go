package samples

import (
	"database/sql"
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

	if err := row.Scan(&sample.SampleID, &sample.FullName, &sample.CreatedAt, &sample.CreatedBy, &sample.Birthdate, &sample.Sputalysed, &sample.Comment); err != nil {
		return nil, err
	}

	return &sample, nil
}

func GetSample(ctx *gin.Context) {
	sample_id := ctx.Param("sample_id")

	sample, err := FetchSampleInformationFromDatabase(sample_id)

	switch err {
	case sql.ErrNoRows:
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Sample not found"})
		return
	case nil:
		break
	default:
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, sample)
}
