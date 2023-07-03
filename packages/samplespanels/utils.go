package samplespanels

import (
	"database/sql"
	"math/rand"

	"gitlab.com/kaka/pcr-backend/common/database"
)

func SampleAnalysisExists(sample_id string, analysis_id string) bool {
	query := `
		SELECT EXISTS(
			SELECT sample_id
			FROM samplespanels
			WHERE sample_id = $1 AND panel_id = $2
		)`
	var exists bool

	if err := database.Instance.QueryRow(query, sample_id, analysis_id).Scan(&exists); err != nil && err != sql.ErrNoRows {
		return false
	}

	return exists
}

// Generate Hash
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
