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
		SELECT analysis_id,ready_mix,is_active
		FROM analyses
		WHERE analysis_id = $1;
		`
	err := database.Instance.QueryRow(query, analysisID).Scan(&analysis.AnalysisId, &analysis.AnalysisId, &analysis.ReadyMix, &analysis.IsActive)

	if err != nil {
		return nil, err
	}

	return &analysis, nil
}

func GetAnalysis(ctx *gin.Context) {

	analysis_id := ctx.Param("analysis_id")

	analysis, err := FetchAnalysisInformationFromDatabase(analysis_id)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, &analysis)
}

type AnalysisQueryParams struct {
	AnalysisId string `form:"analysis"`

	// Pagination
	OrderBy string `form:"order_by"`
	Limit   string `form:"limit"`
}

func getFilterString(query_params *AnalysisQueryParams) (string, []interface{}) {
	var filters []string
	var params []interface{}
	if query_params.AnalysisId != "" {
		filters = append(filters, "analysis_id = $3")
		params = append(params, query_params.AnalysisId)
	}

	query := ""
	if len(filters) > 0 {
		query += " WHERE "
		for i, filter := range filters {
			query += filter
			if i != len(filters)-1 {
				query += " AND "
			}
		}
	}

	return query, params
}

func GetAllAnalyses(ctx *gin.Context) {
	var analyses []models.Analysis
	var query_params = AnalysisQueryParams{
		OrderBy: "analysis_id",
		Limit:   "100",
	}
	ctx.ShouldBindQuery(&query_params)

	// Get all analyses
	query :=
		`
		SELECT analysis_id,ready_mix,is_active
		FROM analyses
		`
	// Add filters
	filter_string, filter_params := getFilterString(&query_params)
	query += filter_string

	// Add pagination
	query += ` ORDER BY $1 LIMIT $2;`
	filter_params = append([]interface{}{query_params.OrderBy, query_params.Limit}, filter_params...)
	rows, err := database.Instance.QueryContext(ctx, query, filter_params...)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	for rows.Next() {
		var analysis models.Analysis

		if err := rows.Scan(&analysis.AnalysisId, &analysis.ReadyMix, &analysis.IsActive); err != nil {
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
