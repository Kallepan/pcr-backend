package samples

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/kaka/pcr-backend/common/database"
	"gitlab.com/kaka/pcr-backend/common/models"
)

func GetAssociatedAnalysis(sample_id string) []models.SampleAnalysis {
	// Get all analyses associated with a sample
	query := `
		SELECT analyses.analyt,analyses.assay,analyses.material,analyses.ready_mix,sampleanalyses.run,sampleanalyses.device,sampleanalyses.created_at,users.username
		FROM sampleanalyses
		LEFT JOIN analyses ON sampleanalyses.analysis_id = analyses.analysis_id
		LEFT JOIN users ON sampleanalyses.created_by = users.user_id
		WHERE sampleanalyses.sample_id = $1
		ORDER BY sampleanalyses.created_at DESC;
		`

	rows, err := database.Instance.Query(query, sample_id)

	switch err {
	case nil:
		break
	default:
		return nil
	}

	var sampleanalyses []models.SampleAnalysis
	for rows.Next() {
		var sampleanalysis models.SampleAnalysis

		if err := rows.Scan(&sampleanalysis.Analyt, &sampleanalysis.Assay, &sampleanalysis.Material, &sampleanalysis.ReadyMix, &sampleanalysis.Run, &sampleanalysis.Device, &sampleanalysis.CreatedAt, &sampleanalysis.CreatedBy); err != nil {
			break
		}
		sampleanalyses = append(sampleanalyses, sampleanalysis)
	}

	if err = rows.Err(); err != nil {
		return nil
	}

	return sampleanalyses
}

func GetSample(ctx *gin.Context) {
	sample_id := ctx.Param("sample_id")

	var sample models.Sample

	query :=
		`SELECT sample_id,samples.firstname,samples.lastname,created_at,users.username
		FROM samples 
		LEFT JOIN users ON samples.created_by = users.user_id
		WHERE sample_id = $1;`

	row := database.Instance.QueryRow(query, sample_id)

	switch err := row.Scan(&sample.SampleID, &sample.FirstName, &sample.LastName, &sample.CreatedAt, &sample.CreatedBy); err {
	case sql.ErrNoRows:
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "sample not found"})
		return
	case nil:
		break
	default:
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	sampleAnalyses := GetAssociatedAnalysis(sample_id)

	sample.AssociatedAnalyses = sampleAnalyses
	ctx.JSON(http.StatusOK, sample)
}
