package repository

import (
	"database/sql"
	"fmt"

	"gitlab.com/kallepan/pcr-backend/app/domain/dao"
	"gitlab.com/kallepan/pcr-backend/app/domain/dco"
	"gitlab.com/kallepan/pcr-backend/driver"
)

type SamplePanelRepository interface {
	GetSamplePanels(sampleID string, run_date string, run string, device string) ([]dao.SamplePanel, error)
	SamplePanelExists(sampleID string, panelID string) bool
	ResetSamplePanel(samplePanel dco.ResetSamplePanelRequest) error
	UpdateSamplePanel(sampleID string, panelID string, request dco.UpdateSamplePanelRequest) error
	CreateSamplePanel(samplePanel dco.AddSamplePanelRequest, userID string) (string, error)
	UndeleteSamplePanel(sampleID string, panelID string) error

	// Statistics
	GetStatistics() ([]dao.Statistic, error)
}

type SamplePanelRepositoryImpl struct {
	db *sql.DB
}

func SamplePanelRepositoryInit() *SamplePanelRepositoryImpl {
	return &SamplePanelRepositoryImpl{
		db: driver.DB,
	}
}

func (s SamplePanelRepositoryImpl) GetStatistics() ([]dao.Statistic, error) {
	/* Returns the number of unfinished samples per analysis */
	query := `
		SELECT
			LEFT(panel_id,3),
			count(*)
		FROM
			samplespanels
		WHERE
			run_date IS NULL AND deleted = FALSE
		GROUP BY
			LEFT(panel_id,3)
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}

	statistics := []dao.Statistic{}
	for rows.Next() {
		var statistic dao.Statistic
		err := rows.Scan(&statistic.PanelID, &statistic.Count)

		if err != nil {
			return nil, err
		}

		statistics = append(statistics, statistic)
	}

	return statistics, nil
}

func (s SamplePanelRepositoryImpl) UndeleteSamplePanel(sampleID string, panelID string) error {
	/* Sets the deleted field to false depending on the current value */

	query := `
		UPDATE samplespanels
		SET deleted = FALSE
		WHERE sample_id = $1 AND panel_id = $2;
	`
	_, err := s.db.Exec(query, sampleID, panelID)

	return err
}

func (s SamplePanelRepositoryImpl) CreateSamplePanel(request dco.AddSamplePanelRequest, userID string) (string, error) {
	/*
	* Creates a samplepanel with the given sample_id and panel_id
	 */
	query := `
	WITH sample AS (
		SELECT samples.sample_id, samples.full_name, samples.created_at, samples.sputalysed, users.username AS created_by
		FROM samples
		LEFT JOIN users ON samples.created_by = users.user_id
		WHERE samples.sample_id = $1
		GROUP BY samples.sample_id, samples.full_name, samples.created_at, samples.sputalysed, users.username
	)
	INSERT INTO samplespanels (sample_id, panel_id, created_by) 
	VALUES ($1, $2, $3)
	RETURNING sample_id, panel_id`

	var (
		panelID  string
		sampleID string
	)

	if err := s.db.QueryRow(query, request.SampleID, request.PanelID, userID).Scan(
		&sampleID,
		&panelID,
	); err != nil {
		return "", err
	}

	message := fmt.Sprintf("Sample %s added to panel %s", sampleID, panelID)
	return message, nil
}

func (s SamplePanelRepositoryImpl) UpdateSamplePanel(sampleID string, panelID string, request dco.UpdateSamplePanelRequest) error {
	/*
		* Updates a samplepanel with the given sample_id and panel_id,
		currently only the deleted field can be updated
	*/

	query := `
		UPDATE samplespanels
		SET deleted = $3
		WHERE sample_id = $1 AND panel_id = $2;
	`
	_, err := s.db.Exec(query, sampleID, panelID, *request.Deleted)

	return err
}

func (s SamplePanelRepositoryImpl) SamplePanelExists(sampleID string, panelID string) bool {
	/* Checks if a sample is associated with a panel */
	query := `
		SELECT EXISTS(
			SELECT sample_id
			FROM samplespanels
			WHERE sample_id = $1 AND panel_id = $2
		)`
	var exists bool
	if err := s.db.QueryRow(query, sampleID, panelID).Scan(&exists); err != nil && err != sql.ErrNoRows {
		return false
	}

	return exists
}

func (s SamplePanelRepositoryImpl) ResetSamplePanel(request dco.ResetSamplePanelRequest) error {
	/*
	* Resets the samplepanel with the given sample_id and panel_id
	 */

	query := `
		UPDATE samplespanels
		SET run = NULL, device = NULL, position = NULL, run_date = NULL
		WHERE
			sample_id = $1 AND
			panel_id = $2 AND
			deleted = false
	`

	if _, err := s.db.Exec(query, request.SampleID, request.PanelID); err != nil {
		return err
	}

	return nil
}

func (s SamplePanelRepositoryImpl) GetSamplePanels(sampleID string, runDate string, run string, device string) ([]dao.SamplePanel, error) {
	/*
	* Returns all samplespanels with an optional filter for sample_id or the
	* triple (run_date, run, device)
	 */

	samplespanels := []dao.SamplePanel{}

	query, params := buildGetQuery(sampleID, runDate, run, device)

	rows, err := s.db.Query(query, params...)
	if err != nil {
		return samplespanels, err
	}

	for rows.Next() {
		var samplePanel dao.SamplePanel
		var sample dao.Sample
		var panel dao.Panel

		if err := rows.Scan(
			&sample.SampleID, &sample.FullName, &sample.CreatedAt, &sample.CreatedBy, &sample.Material,
			&panel.PanelID, &panel.DisplayName, &panel.ReadyMix,
			&samplePanel.Run, &samplePanel.Device, &samplePanel.Position, &samplePanel.RunDate, &samplePanel.CreatedAt, &samplePanel.CreatedBy); err != nil {

			return samplespanels, err
		}

		samplePanel.Sample = sample
		samplePanel.Panel = panel
		samplespanels = append(samplespanels, samplePanel)
	}

	if err = rows.Err(); err != nil {
		return samplespanels, err
	}

	return samplespanels, nil
}

func buildGetQuery(sampleID string, run_date string, device string, run string) (string, []interface{}) {
	var params []interface{}
	paramCounter := 1
	query := `
	WITH sample_query AS (
		SELECT samplespanels.sample_id, samples.full_name, samples.created_at, users.username AS created_by, samples.material 
		FROM samplespanels
		LEFT JOIN samples ON samplespanels.sample_id = samples.sample_id
		LEFT JOIN users ON samples.created_by = users.user_id
		GROUP BY samplespanels.sample_id, samples.full_name, samples.created_at, users.username, samples.material
	) 
	SELECT samplespanels.sample_id, sample_query.full_name, sample_query.created_at, sample_query.created_by, sample_query.material, 
	samplespanels.panel_id, panels.display_name, panels.ready_mix, samplespanels.run, samplespanels.device, samplespanels.position, samplespanels.run_date, samplespanels.created_at, users.username
	FROM samplespanels
	LEFT JOIN sample_query ON samplespanels.sample_id = sample_query.sample_id
	LEFT JOIN panels ON samplespanels.panel_id = panels.panel_id
	LEFT JOIN users ON samplespanels.created_by = users.user_id
	WHERE
		samplespanels.deleted = false`

	if sampleID != "" {
		query += fmt.Sprintf(" AND samplespanels.sample_id = $%d", paramCounter)
		paramCounter++
		params = append(params, sampleID)
	}

	if run_date != "" {
		query += fmt.Sprintf(" AND samplespanels.run_date = $%d", paramCounter)
		paramCounter++
		params = append(params, run_date)
	}

	if run != "" {
		query += fmt.Sprintf(" AND samplespanels.run = $%d", paramCounter)
		paramCounter++
		params = append(params, run)
	}

	if device != "" {
		query += fmt.Sprintf(" AND samplespanels.device = $%d", paramCounter)
		paramCounter++
		params = append(params, device)
	}

	if sampleID == "" && run_date == "" && run == "" && device == "" {
		query += `
		AND samplespanels.run IS NULL AND
		samplespanels.device IS NULL AND
		samplespanels.position IS NULL`
	}

	// Order by
	query += " ORDER BY samplespanels.created_at ASC, samplespanels.sample_id DESC LIMIT 100"

	return query, params
}
