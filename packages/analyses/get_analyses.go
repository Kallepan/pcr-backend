package analyses

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/kaka/pcr-backend/common/database"
	"gitlab.com/kaka/pcr-backend/common/models"
)

func FetchAnalysisInformationFromDatabase(analysisID string) (*models.Analysis, error) {
	var analysis models.Analysis

	query :=
		`
		SELECT analysis_id, display_name, ready_mix, is_active
		FROM analyses
		WHERE analysis_id = $1;
		`
	err := database.Instance.QueryRow(query, analysisID).Scan(&analysis.AnalysisId, &analysis.DisplayName, &analysis.ReadyMix, &analysis.IsActive)

	if err != nil {
		return nil, err
	}

	return &analysis, nil
}

func GetAnalyses(ctx *gin.Context) {
	var analyses []models.Analysis

	analysis_id := ctx.Param("analysis_id")

	// Get all analyses
	query :=
		`
		SELECT analysis_id, display_name, ready_mix, is_active
		FROM analyses
		WHERE 1 = 1
		`
	var params []interface{}

	if analysis_id != "" {
		query += "AND analysis_id = ?"
		params = append(params, analysis_id)
	} else {
		query += "ORDER BY display_name LIMIT 100;"
	}

	// Add pagination
	rows, err := database.Instance.QueryContext(ctx, query, params...)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	for rows.Next() {
		var analysis models.Analysis

		if err := rows.Scan(&analysis.AnalysisId, &analysis.DisplayName, &analysis.ReadyMix, &analysis.IsActive); err != nil {
			break
		}
		analyses = append(analyses, analysis)
	}

	if err = rows.Err(); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, &analyses)
}
