package analyses

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/kaka/pcr-backend/common/database"
	"gitlab.com/kaka/pcr-backend/common/models"
)

func GetAllAnalyses(ctx *gin.Context) {
	var analyses []models.Analysis

	// Get all analyses
	query :=
		`
		SELECT analysis_id,analyt,assay,material,ready_mix
		FROM analyses
		ORDER BY analysis_id ASC;
		`

	rows, err := database.Instance.Query(query)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for rows.Next() {
		var analysis models.Analysis

		if err := rows.Scan(&analysis.AnalysisID, &analysis.Analyt, &analysis.Assay, &analysis.Material, &analysis.ReadyMix); err != nil {
			break
		}
		analyses = append(analyses, analysis)
	}

	if err = rows.Err(); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, &analyses)
}
