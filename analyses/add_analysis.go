package analyses

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/kaka/pcr-backend/common/database"
	"gitlab.com/kaka/pcr-backend/common/models"
)

type AddAnalysisRequest struct {
	Analyt   string `json:"analyt"`
	Material string `json:"material"`
	Assay    string `json:"assay"`
	ReadyMix bool   `json:"ready_mix"`
}

func AnalysisExists(analysis models.Analysis) bool {
	query := `
		SELECT EXISTS(
		SELECT * 
		FROM analyses 
		WHERE analyt = $1 AND material = $2 AND assay = $3)`

	var exists bool
	err := database.Instance.QueryRow(query, analysis.Analyt, analysis.Material, analysis.Assay).Scan(&exists)

	return exists && err != sql.ErrNoRows
}

func AddAnalysis(ctx *gin.Context) {
	var request AddAnalysisRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	analysis := models.Analysis{
		Analyt:   request.Analyt,
		Material: request.Material,
		Assay:    request.Assay,
		ReadyMix: request.ReadyMix,
	}

	// Check if analysis already exists
	if AnalysisExists(analysis) {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "analysis already exists"})
		return
	}

	// Insert analysis
	query := "INSERT INTO analyses (analyt,material,assay,ready_mix) VALUES ($1,$2,$3,$4) RETURNING analyt,material,assay,ready_mix"
	err := database.Instance.QueryRow(query, analysis.Analyt, analysis.Material, analysis.Assay, analysis.ReadyMix).Scan(&analysis.Analyt, &analysis.Material, &analysis.Assay, &analysis.ReadyMix)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, analysis)
}
