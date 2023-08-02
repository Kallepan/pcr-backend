package samplespanels

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

/*

func removeWhitespaces(s string) string {
	// Remove tabs and newlines
	s = strings.ReplaceAll(s, "\t", "")
	s = strings.ReplaceAll(s, "\n", "")
	s = strings.ReplaceAll(s, " ", "")
	return s
}
*/

func TestBuildGetQuery(t *testing.T) {
	// Empty Sample ID
	query, params := buildGetQuery("")

	targetQuery := `
	WITH sample_query AS (
		SELECT samplespanels.sample_id, samples.full_name, samples.created_at, users.username AS created_by
		FROM samplespanels
		LEFT JOIN samples ON samplespanels.sample_id = samples.sample_id
		LEFT JOIN users ON samples.created_by = users.user_id
		GROUP BY samplespanels.sample_id, samples.full_name, samples.created_at, users.username
	) 
	SELECT samplespanels.sample_id, sample_query.full_name, sample_query.created_at, sample_query.created_by, samplespanels.panel_id, panels.display_name, panels.ready_mix, samplespanels.run, samplespanels.device, samplespanels.position, samplespanels.run_date, samplespanels.created_at, users.username
	FROM samplespanels
	LEFT JOIN sample_query ON samplespanels.sample_id = sample_query.sample_id
	LEFT JOIN panels ON samplespanels.panel_id = panels.panel_id
	LEFT JOIN users ON samplespanels.created_by = users.user_id
	WHERE
		samplespanels.deleted = false
		AND samplespanels.run IS NULL AND
		samplespanels.device IS NULL AND
		samplespanels.position IS NULL ORDER BY samplespanels.created_at DESC, samplespanels.sample_id DESC`

	assert.Equal(t, targetQuery, query)
	assert.Equal(t, 0, len(params))

	// Non-empty Sample ID
	query, params = buildGetQuery("BLAH")

	targetQuery = `
	WITH sample_query AS (
		SELECT samplespanels.sample_id, samples.full_name, samples.created_at, users.username AS created_by
		FROM samplespanels
		LEFT JOIN samples ON samplespanels.sample_id = samples.sample_id
		LEFT JOIN users ON samples.created_by = users.user_id
		GROUP BY samplespanels.sample_id, samples.full_name, samples.created_at, users.username
	) 
	SELECT samplespanels.sample_id, sample_query.full_name, sample_query.created_at, sample_query.created_by, samplespanels.panel_id, panels.display_name, panels.ready_mix, samplespanels.run, samplespanels.device, samplespanels.position, samplespanels.run_date, samplespanels.created_at, users.username
	FROM samplespanels
	LEFT JOIN sample_query ON samplespanels.sample_id = sample_query.sample_id
	LEFT JOIN panels ON samplespanels.panel_id = panels.panel_id
	LEFT JOIN users ON samplespanels.created_by = users.user_id
	WHERE
		samplespanels.deleted = false AND samplespanels.sample_id = $1 ORDER BY samplespanels.created_at DESC, samplespanels.sample_id DESC`

	assert.Equal(t, targetQuery, query)
	assert.Equal(t, 1, len(params))
}
