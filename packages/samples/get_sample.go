package samples

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/kaka/pcr-backend/common/database"
	"gitlab.com/kaka/pcr-backend/common/models"
)

func GetAssociatedAnalysis(tagesnummer string) []models.SampleAnalysis {
	// Get all analyses associated with a sample
	query := `
		SELECT analyses.analyt,analyses.assay,analyses.material,analyses.ready_mix,
		sampleanalyses.id,sampleanalyses.run,sampleanalyses.device,sampleanalyses.created_at,users.username
		FROM sampleanalyses
		LEFT JOIN analyses ON sampleanalyses.analysis_id = analyses.analysis_id
		LEFT JOIN users ON sampleanalyses.created_by = users.user_id
		WHERE sampleanalyses.tagesnummer = $1
		ORDER BY sampleanalyses.created_at DESC;
		`

	rows, err := database.Instance.Query(query, tagesnummer)

	switch err {
	case nil:
		break
	default:
		return nil
	}

	var sampleanalyses []models.SampleAnalysis
	for rows.Next() {
		var sampleanalysis models.SampleAnalysis

		if err := rows.Scan(&sampleanalysis.Analyt, &sampleanalysis.Assay, &sampleanalysis.Material, &sampleanalysis.ReadyMix, &sampleanalysis.ID, &sampleanalysis.Run, &sampleanalysis.Device, &sampleanalysis.CreatedAt, &sampleanalysis.CreatedBy); err != nil {
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
	tagesnummer := ctx.Param("tagesnummer")

	var sample models.Sample

	query :=
		`SELECT tagesnummer,name,created_at,users.username
		FROM samples 
		LEFT JOIN users ON samples.created_by = users.user_id
		WHERE tagesnummer = $1;`

	row := database.Instance.QueryRow(query, tagesnummer)

	switch err := row.Scan(&sample.Tagesnummer, &sample.Name, &sample.CreatedAt, &sample.CreatedBy); err {
	case sql.ErrNoRows:
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "sample not found"})
		return
	case nil:
		break
	default:
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	sampleAnalyses := GetAssociatedAnalysis(tagesnummer)

	sample.AssociatedAnalyses = sampleAnalyses
	ctx.JSON(http.StatusOK, sample)
}
