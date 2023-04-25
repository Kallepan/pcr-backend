package sampleanalyses

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/kaka/pcr-backend/common/database"
	"gitlab.com/kaka/pcr-backend/common/models"
)

func GetSampleAnalyses(ctx *gin.Context) {
	var sampleAnalyses []models.SampleAnalysis

	query := `
		WITH sample_query AS (
			SELECT sampleanalyses.sample_id, samples.firstname, samples.lastname, samples.created_at, users.username AS created_by
			FROM sampleanalyses
			LEFT JOIN samples ON sampleanalyses.sample_id = samples.sample_id
			LEFT JOIN users ON samples.created_by = users.user_id
			GROUP BY sampleanalyses.sample_id, samples.firstname, samples.lastname, samples.created_at, users.username
		) 
		SELECT sampleanalyses.sample_id, sample_query.firstname, sample_query.lastname, sample_query.created_at, sample_query.created_by, sampleanalyses.analysis_id, analyses.analyt, analyses.material, analyses.assay, analyses.ready_mix, sampleanalyses.run, sampleanalyses.device, sampleanalyses.completed, sampleanalyses.created_at, users.username
		FROM sampleanalyses
		LEFT JOIN sample_query ON sampleanalyses.sample_id = sample_query.sample_id
		LEFT JOIN analyses ON sampleanalyses.analysis_id = analyses.analysis_id
		LEFT JOIN users ON sampleanalyses.created_by = users.user_id

		ORDER BY $1 DESC LIMIT $2
	`

	rows, err := database.Instance.Query(query, "sampleanalyses.created_at", 100)

	if err != nil {
		ctx.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	for rows.Next() {
		var sampleAnalysis models.SampleAnalysis
		var sample models.Sample
		var analysis models.Analysis

		if err := rows.Scan(
			&sample.SampleID, &sample.FirstName, &sample.LastName, &sample.CreatedAt, &sample.CreatedBy,
			&analysis.AnalysisID, &analysis.Analyt, &analysis.Material, &analysis.Assay, &analysis.ReadyMix,
			&sampleAnalysis.Run, &sampleAnalysis.Device, &sampleAnalysis.Completed, &sampleAnalysis.CreatedAt, &sampleAnalysis.CreatedBy); err != nil {

			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		sampleAnalysis.Sample = sample
		sampleAnalysis.Analysis = analysis
		sampleAnalyses = append(sampleAnalyses, sampleAnalysis)
	}

	if err = rows.Err(); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, &sampleAnalyses)
}
