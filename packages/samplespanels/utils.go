package samplespanels

import (
	"database/sql"
	"math/rand"
	"strings"

	"gitlab.com/kaka/pcr-backend/common/database"
)

func ExtractLastRunId(sample_id string) (string, error) {
	query := `
		SELECT CONCAT(device, '-POS', position, '-', run_date)
		FROM samplespanels
		WHERE 
			sample_id = $1 AND
			run_date IS NOT NULL AND
			device IS NOT NULL AND
			position IS NOT NULL
		ORDER BY run_date ASC
		LIMIT 1;
	`
	var run_id *string

	if err := database.Instance.QueryRow(query, sample_id).Scan(&run_id); err != nil && err != sql.ErrNoRows {
		return "-", err
	}

	if run_id == nil {
		return "-", nil
	}

	return *run_id, nil
}

// SamplePanelExists checks if a sample is associated with a panel
func SamplePanelExists(sample_id string, analysis_id string) bool {
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

// splitBySpace splits the string by space
func splitBySpace(s string) []string {
	return strings.Split(s, " ")
}
