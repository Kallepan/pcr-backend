package samples

import (
	"database/sql"

	"gitlab.com/kaka/pcr-backend/common/database"
)

func SampleExists(sample_id string) bool {
	query := `
		SELECT EXISTS(
			SELECT sample_id 
			FROM samples 
			WHERE sample_id = $1
		)`

	var exists bool
	err := database.Instance.QueryRow(query, sample_id).Scan(&exists)

	if err != nil && err != sql.ErrNoRows {
		return false
	}

	return exists
}
