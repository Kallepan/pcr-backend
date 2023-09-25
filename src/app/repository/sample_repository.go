package repository

import (
	"database/sql"
	"fmt"

	"gitlab.com/kallepan/pcr-backend/app/domain/dao"
	"gitlab.com/kallepan/pcr-backend/app/domain/dco"
	"gitlab.com/kallepan/pcr-backend/driver"
)

type SampleRepository interface {
	GetSamples(sampleID string) ([]dao.Sample, error)
	GetSample(sampleID string) (dao.Sample, error)
	SampleExists(sampleID string) bool
	CreateSample(sample dco.AddSampleRequest, userID string) (dao.Sample, error)
	UpdateSample(sample dco.UpdateSampleRequest, sampleID string) (dao.Sample, error)
	DeleteSample(sampleID string) error
}

type SampleRepositoryImpl struct {
	db *sql.DB
}

func SampleRepositoryInit() *SampleRepositoryImpl {
	return &SampleRepositoryImpl{
		db: driver.DB,
	}
}

func (s SampleRepositoryImpl) GetSample(sampleID string) (dao.Sample, error) {
	/* Returns a single sample */
	sample := dao.Sample{}

	// Get all samples
	query := `
		SELECT s.sample_id, s.full_name, s.birthdate, s.sputalysed, s.comment, s.created_at, u.username, s.material
		FROM samples s
		LEFT JOIN users u ON s.created_by = u.user_id
		WHERE s.sample_id = $1
	`

	// query
	err := s.db.QueryRow(query, sampleID).Scan(&sample.SampleID,
		&sample.FullName,
		&sample.Birthdate,
		&sample.Sputalysed,
		&sample.Comment,
		&sample.CreatedAt,
		&sample.CreatedBy,
		&sample.Material,
	)
	if err != nil {
		return sample, err
	}

	return sample, nil
}

func (s SampleRepositoryImpl) GetSamples(sampleID string) ([]dao.Sample, error) {
	/* Returns all samples with an optional filter for sample_id */
	samples := []dao.Sample{}

	// Get all samples and filter out inactive ones
	query := `
		SELECT s.sample_id, s.full_name, s.birthdate, s.sputalysed, s.comment, s.created_at, u.username, s.material, string_agg(samplespanels.panel_id, ', ') AS panels
		FROM samples s
		LEFT JOIN users u ON s.created_by = u.user_id
		LEFT JOIN samplespanels ON samplespanels.sample_id = s.sample_id AND samplespanels.deleted = false
		WHERE 1 = 1
	`

	// Add pagination and filters
	var params []interface{}
	if sampleID != "" {
		query += "AND s.sample_id LIKE $1"

		// Format param for LIKE query
		param := fmt.Sprintf("%%%s%%", sampleID)
		params = append(params, param)
	}

	// Additional Filters:
	// samples younger than 14 days
	// Group by sample_id and order by created_at
	query += `	
		AND s.created_at >= current_date - interval '14 day'
		GROUP BY s.sample_id, u.username ORDER BY s.created_at DESC, s.sample_id DESC LIMIT 1000;
	`

	// query
	rows, err := s.db.Query(query, params...)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var sample dao.Sample
		var panels sql.NullString
		if err := rows.Scan(&sample.SampleID,
			&sample.FullName,
			&sample.Birthdate,
			&sample.Sputalysed,
			&sample.Comment,
			&sample.CreatedAt,
			&sample.CreatedBy,
			&sample.Material,
			&panels); err != nil {
			return nil, err
		}

		if panels.Valid {
			sample.Panels = panels.String
		} else {
			sample.Panels = "N/A"
		}

		samples = append(samples, sample)
	}

	return samples, nil
}

func (s SampleRepositoryImpl) SampleExists(sampleID string) bool {
	/* Returns true if sample exists */

	query := `
		SELECT EXISTS(
			SELECT sample_id 
			FROM samples 
			WHERE sample_id = $1
		)`
	var exists bool
	err := s.db.QueryRow(query, sampleID).Scan(&exists)
	if err != nil {
		return false
	}

	return exists
}

func (s SampleRepositoryImpl) CreateSample(request dco.AddSampleRequest, userID string) (dao.Sample, error) {
	/* Creates a new sample using a custom add sample request struct */
	sample := dao.Sample{}

	// Insert sample
	query := `
		WITH new_sample AS
		(
			INSERT INTO samples (
				sample_id,
				full_name,
				sputalysed,
				comment,
				birthdate,
				created_by,
				manual,
				material
			)
			VALUES ($1, $2, $3, $4, $5, $6, TRUE, $7)
			RETURNING sample_id, full_name, created_at, sputalysed, comment, birthdate, created_by, material
		)
		SELECT sample_id, full_name, created_at, sputalysed, comment, birthdate, users.username, material
		FROM new_sample
		LEFT JOIN users ON new_sample.created_by = users.user_id
	`
	err := s.db.QueryRow(
		query,
		request.SampleId,
		request.FullName,
		request.Sputalysed,
		request.Comment,
		request.Birthdate,
		userID,
		request.Material,
	).Scan(
		&sample.SampleID,
		&sample.FullName,
		&sample.CreatedAt,
		&sample.Sputalysed,
		&sample.Comment,
		&sample.Birthdate,
		&sample.CreatedBy,
		&sample.Material,
	)
	if err != nil {
		return sample, err
	}

	// Return the new sample
	return sample, nil
}

func (s SampleRepositoryImpl) UpdateSample(sample dco.UpdateSampleRequest, sampleID string) (dao.Sample, error) {
	/* Updates a sample using a custom update sample request struct */
	updatedSample := dao.Sample{}

	// Update sample, I know this is ugly but I don't know how to do it better
	// without using a library like squirrel
	var params []interface{}
	query := `
		WITH updated_sample AS (
			UPDATE samples SET`
	if sample.FullName != nil {
		query += fmt.Sprintf(" full_name = $%d", len(params)+1)
		params = append(params, *sample.FullName)
	}
	if sample.Sputalysed != nil {
		query += fmt.Sprintf(", sputalysed = $%d", len(params)+1)
		params = append(params, *sample.Sputalysed)
	}
	if sample.Comment != nil {
		query += fmt.Sprintf(", comment = $%d", len(params)+1)
		params = append(params, *sample.Comment)
	}

	query += fmt.Sprintf(" WHERE sample_id = $%d", len(params)+1)
	params = append(params, sampleID)
	query += ` RETURNING sample_id, full_name, created_at, sputalysed, comment, birthdate, created_by, material
		)
		SELECT sample_id, full_name, created_at, sputalysed, comment, birthdate, users.username, material
		FROM updated_sample
		LEFT JOIN users ON updated_sample.created_by = users.user_id;
	`
	err := s.db.QueryRow(
		query,
		params...,
	).Scan(
		&updatedSample.SampleID,
		&updatedSample.FullName,
		&updatedSample.CreatedAt,
		&updatedSample.Sputalysed,
		&updatedSample.Comment,
		&updatedSample.Birthdate,
		&updatedSample.CreatedBy,
		&updatedSample.Material,
	)

	if err != nil {
		return updatedSample, err
	}

	// Return the updated sample
	return updatedSample, nil
}

func (s SampleRepositoryImpl) DeleteSample(sampleID string) error {
	/* Deletes a sample */
	query := `DELETE FROM samples WHERE sample_id = $1`

	_, err := s.db.Exec(query, sampleID)

	return err
}
