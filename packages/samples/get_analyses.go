package samples

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/kaka/pcr-backend/common/database"
	"gitlab.com/kaka/pcr-backend/common/models"
)

func GetAnalysesForSample(ctx *gin.Context) {
	var analyses []models.Analysis

	query := `
		WITH samples_analyses AS (
			SELECT samplesanalyses.analysis_id
			FROM samplesanalyses
			WHERE samplesanalyses.sample_id = $1
		)
		SELECT analyses.analysis_id, analyses.analyt, analyses.material, analyses.assay, analyses.ready_mix
		FROM analyses
		RIGHT JOIN samples_analyses ON analyses.analysis_id = samples_analyses.analysis_id
	`

	rows, err := database.Instance.Query(query, ctx.Param("sample_id"))

	if err != nil {
		ctx.AbortWithStatusJSON(500, gin.H{"message": err.Error()})
		return
	}

	for rows.Next() {
		var analysis models.Analysis

		if err := rows.Scan(&analysis.AnalysisID, &analysis.Analyt, &analysis.Material, &analysis.Assay, &analysis.ReadyMix); err != nil {
			ctx.AbortWithStatusJSON(500, gin.H{"message": err.Error()})
		}

		analyses = append(analyses, analysis)
	}

	if err = rows.Err(); err != nil {
		ctx.AbortWithStatusJSON(500, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(200, analyses)
}
