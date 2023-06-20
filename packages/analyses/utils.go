package analyses

import (
	"database/sql"

	"gitlab.com/kaka/pcr-backend/common/database"
)

func AnalysisExists(AnalysisId string) bool {
	query := `
		SELECT EXISTS(
		SELECT * 
		FROM analyses 
		WHERE analysis_id = $1)`

	var exists bool
	err := database.Instance.QueryRow(query, AnalysisId).Scan(&exists)

	return exists && err != sql.ErrNoRows
}
