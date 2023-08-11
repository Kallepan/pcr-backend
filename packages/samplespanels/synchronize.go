/*
Synchronizes the ingenious table with the sampleanalyses and samples table
*/
package samplespanels

import (
	"database/sql"
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

func synchronizeAnalysesTable(tx *sql.Tx) error {
	/*
		Ensure that the analyses and panels table has the same entries
		This function creates a new panel if it does not exist in the panels table.
		Furthermore, it creates a new analysis if it does not exist in the analyses table.
	*/

	_, err := tx.Exec(`
	INSERT INTO panels (panel_id, display_name)
	SELECT DISTINCT ingenious.usi as panel_id, ingenious.usi AS display_name
	FROM ingenious
	LEFT JOIN analyses ON analyses.analysis_id = ingenious.usi
	WHERE analyses.analysis_id IS NULL;

	INSERT INTO analyses (analysis_id, panel_id)
	SELECT DISTINCT ingenious.usi AS analysis_id, ingenious.usi AS panel_id
	FROM ingenious
	LEFT JOIN analyses ON analyses.analysis_id = ingenious.usi
	WHERE analyses.analysis_id IS NULL;
	`)

	return err
}

func synchronizeSamples(tx *sql.Tx) error {
	/*
		Ensure that the samples table has the same entries as the ingenious table
	*/

	_, err := tx.Exec(`
	INSERT INTO samples (sample_id, birthdate, full_name, created_by)
		SELECT DISTINCT ON (ingenious.barcode, ingenious.birthdate, ingenious.patient) ingenious.barcode, ingenious.birthdate, ingenious.patient, users.user_id
		FROM ingenious
		LEFT JOIN samples ON samples.sample_id = ingenious.barcode
		LEFT JOIN (
			SELECT user_id
			FROM users
			LIMIT 1
		) users ON 1=1
		WHERE samples.sample_id IS NULL 
		AND ingenious.barcode IS NOT NULL AND ingenious.patient IS NOT NULL;
	
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
		WHERE ingenious.barcode IS NOT NULL AND ingenious.patient IS NOT NULL
	) 
	INSERT INTO samplespanels (sample_id, panel_id, created_by)
		SELECT DISTINCT filtered_samples.barcode, filtered_samples.panel_id, filtered_samples.user_id
		FROM filtered_samples
		LEFT JOIN samplespanels
		ON samplespanels.sample_id = filtered_samples.barcode AND
		samplespanels.panel_id = filtered_samples.panel_id
		WHERE samplespanels.sample_id IS NULL AND  samplespanels.panel_id IS NULL;
	`)

	return err
}

func deleteOutdatedSamplesPanels(tx *sql.Tx) error {
	/*
		Delete all samplespanels entries where the first three letters of the panel_id and the first ten digits of the sample_id are the same except the youngest entry.

		Afterwards, the same is done for the first three digits of the panel_id and the first ten digits of the sample_id.

		Example:
		Entries:
		- TBXA, 123456789002, 2020-01-01 00:00:00
		- TBXB, 123456789002, 2020-01-01 00:00:01
		- TBXA, 123456789003, 2020-01-01 00:00:02

		After the query, only the last entry will remain.
	*/

	_, err := tx.Exec(`
	-- Delete all samplespanels entries by the sample_id where the first ten digits and the first three letters of the panel_id are the same except the youngest entry.
	DELETE FROM samplespanels sm
	WHERE (LEFT(sm.panel_id, 3), LEFT(sm.sample_id, 10), sm.created_at) NOT IN (
		 SELECT LEFT(sm.panel_id, 3), LEFT(sm.sample_id, 10), MAX(sm.created_at)
		 FROM samplespanels sm
		 GROUP BY LEFT(sm.panel_id,3) , LEFT(sm.sample_id, 10)
	) AND (LEFT(sm.panel_id, 3), LEFT(sm.sample_id, 10)) IN (
		SELECT LEFT(sm.panel_id, 3), LEFT(sm.sample_id, 10)
		FROM samplespanels sm
		GROUP BY LEFT(sm.panel_id, 3), LEFT(sm.sample_id, 10)
		HAVING COUNT(*) > 1
	) AND sm.sample_id IN (
		SELECT sample_id
		FROM samples
		WHERE manual = true
	);
	`)

	return err
}

func deleteEmptySamples(tx *sql.Tx) error {
	/*
		Delete all samples entries where the manual field is false and the sample_id is not in the samplespanels table.
	*/

	_, err := tx.Exec(`
	DELETE FROM samples
	WHERE sample_id NOT IN (
		SELECT sample_id
		FROM samplespanels
	) AND manual = false;
	`)

	return err
}

func synchronize() {
	syncLock.Lock()
	defer syncLock.Unlock()
	tx, err := database.Instance.Begin()
	if err != nil {
		log.Println(err)
		return
	}
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
			log.Println(err)
		}
	}()

	if err := synchronizeAnalysesTable(tx); err != nil {
		log.Println("Error while synchronizing analyses table")
		log.Println(err)
		tx.Rollback()
		return
	}

	if err := synchronizeSamples(tx); err != nil {
		log.Println("Error while synchronizing samples table")
		log.Println(err)
		tx.Rollback()
		return
	}

	if err := deleteOutdatedSamplesPanels(tx); err != nil {
		log.Println("Error while deleting outdated samplespanels entries")
		log.Println(err)
		tx.Rollback()
		return
	}

	if err := deleteEmptySamples(tx); err != nil {
		log.Println("Error while deleting empty samples entries")
		log.Println(err)
		tx.Rollback()
		return
	}

	if err := tx.Commit(); err != nil {
		log.Println(err)
		return
	}

	// Analyze the database to update the statistics
	if _, err := database.Instance.Exec("ANALYZE"); err != nil {
		log.Println(err)
		return
	}

	log.Printf("Synchronization completed: %s", time.Now().Format(time.RFC3339))
}
