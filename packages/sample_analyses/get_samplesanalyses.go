package samplesanalyses

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/kaka/pcr-backend/common/database"
	"gitlab.com/kaka/pcr-backend/common/models"
)

func GetSamplesAnalyses(ctx *gin.Context) {
	var samplesAnalyses []models.SampleAnalysis

	sample_id := ctx.Query("sample_id")

	var params []interface{}

	query := `
		WITH sample_query AS (
			SELECT samplesanalyses.sample_id, samples.full_name, samples.created_at, users.username AS created_by
			FROM samplesanalyses
			LEFT JOIN samples ON samplesanalyses.sample_id = samples.sample_id
			LEFT JOIN users ON samples.created_by = users.user_id
			GROUP BY samplesanalyses.sample_id, samples.full_name, samples.created_at, users.username
		) 
		SELECT samplesanalyses.sample_id, sample_query.full_name, sample_query.created_at, sample_query.created_by, samplesanalyses.analysis_id, analyses.display_name, analyses.ready_mix, analyses.is_active, samplesanalyses.run, samplesanalyses.device, samplesanalyses.position, samplesanalyses.created_at, users.username
		FROM samplesanalyses
		LEFT JOIN sample_query ON samplesanalyses.sample_id = sample_query.sample_id
		LEFT JOIN analyses ON samplesanalyses.analysis_id = analyses.analysis_id
		LEFT JOIN users ON samplesanalyses.created_by = users.user_id
		WHERE
			1 = 1 AND
		`
	if sample_id != "" {
		query += "samplesanalyses.sample_id = ?"
		params = append(params, sample_id)
	} else {
		query += `
			samplesanalyses.deleted = false AND 
			samplesanalyses.run IS NULL AND
			samplesanalyses.device IS NULL AND
			samplesanalyses.position IS NULL
			ORDER BY samplesanalyses.created_at DESC, samplesanalyses.sample_id DESC`
	}
	rows, err := database.Instance.Query(query, params...)

	if err != nil {
		ctx.AbortWithStatusJSON(500, gin.H{"message": err.Error()})
		return
	}

	for rows.Next() {
		var sampleAnalysis models.SampleAnalysis
		var sample models.Sample
		var analysis models.Analysis

		if err := rows.Scan(
			&sample.SampleID, &sample.FullName, &sample.CreatedAt, &sample.CreatedBy,
			&analysis.AnalysisId, &analysis.DisplayName, &analysis.ReadyMix, &analysis.IsActive,
			&sampleAnalysis.Run, &sampleAnalysis.Device, &sampleAnalysis.Position, &sampleAnalysis.CreatedAt, &sampleAnalysis.CreatedBy); err != nil {

			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		}

		sampleAnalysis.Sample = sample
		sampleAnalysis.Analysis = analysis
		samplesAnalyses = append(samplesAnalyses, sampleAnalysis)
	}

	if err = rows.Err(); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(200, &samplesAnalyses)
}
