package samplesanalyses

import (
	"database/sql"

	"gitlab.com/kaka/pcr-backend/common/database"
)

func SampleAnalysisExists(sample_id string, analysis_id string) bool {
	query := `
		SELECT EXISTS(
			SELECT sample_id
			FROM samplesanalyses
			WHERE sample_id = $1 AND analysis_id = $2
		)`
	var exists bool

	err := database.Instance.QueryRow(query, sample_id, analysis_id).Scan(&exists)

	if err != nil && err != sql.ErrNoRows {
		return false
	}

	return exists
}
