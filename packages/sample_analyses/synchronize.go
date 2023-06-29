/*
Synchronizes the ingenious table with the sampleanalyses and samples table
*/
package samplesanalyses

import (
	"log"
	"time"

	"gitlab.com/kaka/pcr-backend/common/database"
)

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
	_, err := database.Instance.Exec(`
	INSERT INTO analyses (analysis_id)
	SELECT DISTINCT ingenious.usi
	FROM ingenious
	LEFT JOIN analyses
	ON analyses.analysis_id = ingenious.usi
	WHERE analyses.analysis_id IS NULL
	`)
	if err != nil {
		log.Println(err)
		return
	}

	_, err = database.Instance.Exec(`
	INSERT INTO samples (sample_id, birthdate, full_name, created_by)
	SELECT DISTINCT ingenious.barcode,ingenious.birthdate, ingenious.patient, users.user_id
	FROM ingenious
	LEFT JOIN samples
	ON samples.sample_id = ingenious.barcode
	LEFT JOIN (
		SELECT user_id
		FROM users
		WHERE users.is_admin = true
		LIMIT 1
	) users ON 1=1
	WHERE samples.sample_id IS NULL
	`)
	if err != nil {
		log.Println(err)
		return
	}

	_, err = database.Instance.Exec(`
	WITH created_sample AS (
		INSERT INTO samples (sample_id, birthdate, full_name, created_by)
		SELECT DISTINCT ingenious.barcode,ingenious.birthdate, ingenious.patient, users.user_id
		FROM ingenious
		LEFT JOIN samples
		ON samples.sample_id = ingenious.barcode
		LEFT JOIN (
			SELECT user_id
			FROM users
			WHERE users.is_admin = true
			LIMIT 1
		) users ON 1=1
		WHERE samples.sample_id IS NULL
	)
	INSERT INTO samplesanalyses (sample_id, analysis_id, created_by)
	SELECT DISTINCT ingenious.barcode, ingenious.usi,users.user_id
	FROM ingenious
	LEFT JOIN samplesanalyses
	ON samplesanalyses.sample_id = ingenious.barcode AND 
	samplesanalyses.analysis_id = ingenious.usi
	LEFT JOIN (
		SELECT user_id
		FROM users
		LIMIT 1
	) users ON 1=1
	WHERE samplesanalyses.sample_id IS NULL AND  
	samplesanalyses.analysis_id IS NULL
	`)
	if err != nil {
		log.Println(err)
		return
	}

}
