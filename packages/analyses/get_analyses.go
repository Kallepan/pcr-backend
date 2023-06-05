package analyses

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/kaka/pcr-backend/common/database"
	"gitlab.com/kaka/pcr-backend/common/models"
)

func GetAnalysis(ctx *gin.Context) {
	var anlysis models.Analysis

	analysis_id := ctx.Param("analysis_id")

	query :=
		`
		SELECT analysis_id,analyt,assay,material,ready_mix
		FROM analyses
		WHERE analysis_id = $1;
		`
	err := database.Instance.QueryRow(query, analysis_id).Scan(&anlysis.AnalysisID, &anlysis.Analyt, &anlysis.Assay, &anlysis.Material, &anlysis.ReadyMix)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, &anlysis)
}

type AnalysisQueryParams struct {
	Analyt   string `form:"analyt"`
	Assay    string `form:"assay"`
	Material string `form:"material"`

	// Pagination
	OrderBy string `form:"order_by"`
	Limit   string `form:"limit"`
}

func getFilterString(query_params *AnalysisQueryParams) (string, []interface{}) {
	var filters []string
	var params []interface{}
	if query_params.Analyt != "" {
		filters = append(filters, "analyt = $3")
		params = append(params, query_params.Analyt)
	}
	if query_params.Assay != "" {
		filters = append(filters, "assay = $4")
		params = append(params, query_params.Assay)
	}
	if query_params.Material != "" {
		filters = append(filters, "material = $5")
		params = append(params, query_params.Material)
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
		SELECT analysis_id,analyt,assay,material,ready_mix
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

		if err := rows.Scan(&analysis.AnalysisID, &analysis.Analyt, &analysis.Assay, &analysis.Material, &analysis.ReadyMix); err != nil {
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
