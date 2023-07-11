/*
Synchronizes the ingenious table with the sampleanalyses and samples table
*/
package samplespanels

import (
	"log"
	"sync"
	"time"

	"gitlab.com/kaka/pcr-backend/common/database"
)

var syncLock sync.Mutex

func StartSynchronize(interval time.Duration) {
	// Execute once before starting the interval
	synchronize()
	log.Println("Synchronization started")

	// Start the interval and synchronize every interval by starting a goroutine
	go func() {
		for {
			synchronize()
			time.Sleep(interval)
		}
	}()
}

func synchronize() {
	syncLock.Lock()
	defer syncLock.Unlock()

	// Ensure that the analyses and panels table has the same entries
	_, err := database.Instance.Exec(`
	BEGIN TRANSACTION;

	WITH new_analyses AS (
		SELECT DISTINCT ingenious.usi AS panel_id
		FROM ingenious
		LEFT JOIN analyses
		ON analyses.analysis_id = ingenious.usi
		WHERE analyses.analysis_id IS NULL
	)
	INSERT INTO panels (panel_id, display_name) 
		SELECT panel_id, panel_id
		FROM new_analyses;

	WITH new_analyses AS (
		SELECT DISTINCT ingenious.usi AS panel_id
		FROM ingenious
		LEFT JOIN analyses
		ON analyses.analysis_id = ingenious.usi
		WHERE analyses.analysis_id IS NULL
	)
	INSERT INTO analyses (analysis_id, panel_id)
		SELECT panel_id, panel_id
		FROM new_analyses;

	COMMIT;`)
	if err != nil {
		log.Println(err)
		return
	}

	_, err = database.Instance.Exec(`
	BEGIN TRANSACTION;

	INSERT INTO samples (sample_id, birthdate, full_name, created_by)
		SELECT DISTINCT ON (ingenious.barcode, ingenious.birthdate, ingenious.patient) ingenious.barcode, ingenious.birthdate, ingenious.patient, users.user_id
		FROM ingenious
		LEFT JOIN samples
		ON samples.sample_id = ingenious.barcode
		LEFT JOIN (
			SELECT user_id
			FROM users
			LIMIT 1
		) users ON 1=1
		WHERE samples.sample_id IS NULL;
	
	WITH filtered_samples AS (
		SELECT DISTINCT ingenious.barcode, analyses.panel_id, users.user_id
		FROM ingenious
		LEFT JOIN analyses
		ON analyses.analysis_id = ingenious.usi
		LEFT JOIN (
			SELECT user_id
			FROM users
			LIMIT 1
		) users ON 1=1
	) 
	INSERT INTO samplespanels (sample_id, panel_id, created_by)
		SELECT DISTINCT filtered_samples.barcode, filtered_samples.panel_id, filtered_samples.user_id
		FROM filtered_samples
		LEFT JOIN samplespanels
		ON samplespanels.sample_id = filtered_samples.barcode AND
		samplespanels.panel_id = filtered_samples.panel_id
		WHERE samplespanels.sample_id IS NULL AND  samplespanels.panel_id IS NULL;
	COMMIT;
	`)
	if err != nil {
		log.Println(err)
		return
	}

	log.Printf("Synchronization completed: %s", time.Now().Format(time.RFC3339))
}
