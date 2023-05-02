package analyses

import (
	"database/sql"

	"gitlab.com/kaka/pcr-backend/common/database"
	"gitlab.com/kaka/pcr-backend/common/models"
)

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

func AnalysisExistsByID(analysis_id string) (models.Analysis, error) {
	var analysis models.Analysis

	query := `
		SELECT analysis_id,analyt,assay,material,ready_mix 
		FROM analyses 
		WHERE analysis_id = $1`

	err := database.Instance.QueryRow(query, analysis_id).Scan(&analysis.AnalysisID, &analysis.Analyt, &analysis.Assay, &analysis.Material, &analysis.ReadyMix)

	return analysis, err
}
