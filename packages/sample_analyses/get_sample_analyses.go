package samplesanalyses

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/kaka/pcr-backend/common/database"
	"gitlab.com/kaka/pcr-backend/common/models"
)

func GetSamplesAnalyses(ctx *gin.Context) {
	var samplesAnalyses []models.SampleAnalysis

	query := `
		WITH sample_query AS (
			SELECT samplesanalyses.sample_id, samples.firstname, samples.lastname, samples.created_at, users.username AS created_by
			FROM samplesanalyses
			LEFT JOIN samples ON samplesanalyses.sample_id = samples.sample_id
			LEFT JOIN users ON samples.created_by = users.user_id
			GROUP BY samplesanalyses.sample_id, samples.firstname, samples.lastname, samples.created_at, users.username
		) 
		SELECT samplesanalyses.sample_id, sample_query.firstname, sample_query.lastname, sample_query.created_at, sample_query.created_by, samplesanalyses.analysis_id, analyses.analyt, analyses.material, analyses.assay, analyses.ready_mix, samplesanalyses.run, samplesanalyses.device, samplesanalyses.completed, samplesanalyses.created_at, users.username
		FROM samplesanalyses
		LEFT JOIN sample_query ON samplesanalyses.sample_id = sample_query.sample_id
		LEFT JOIN analyses ON samplesanalyses.analysis_id = analyses.analysis_id
		LEFT JOIN users ON samplesanalyses.created_by = users.user_id
		WHERE samplesanalyses.completed = false
		ORDER BY $1 DESC LIMIT $2
	`

	rows, err := database.Instance.Query(query, "samplesanalyses.created_at", 100)

	if err != nil {
		ctx.AbortWithStatusJSON(500, gin.H{"message": err.Error()})
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
