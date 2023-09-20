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

func (s SampleRepositoryImpl) GetSamples(sampleID string) ([]dao.Sample, error) {
	/* Returns all samples with an optional filter for sample_id */
	samples := []dao.Sample{}

	// Get all samples and filter out inactive ones
	query := `
		SELECT s.sample_id, s.full_name, s.birthdate, s.comment, s.created_at, u.username, s.material, string_agg(samplespanels.panel_id, ', ') AS panels
		FROM samples s
		LEFT JOIN users u ON s.created_by = u.user_id
		LEFT JOIN samplespanels ON samplespanels.sample_id = s.sample_id AND samplespanels.deleted = false
		WHERE is_active = TRUE
	`

	// Add pagination and filters
	var params []interface{}
	if sampleID != "" {
		query += "AND sample_id = $1"

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
		if err := rows.Scan(&sample.SampleId,
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

func (s SampleRepositoryImpl) CreateSample(sample dco.AddSampleRequest, userID string) (dao.Sample, error) {
	/* Creates a new sample using a custom add sample request struct */
	newSample := dao.Sample{}

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
			RETURNING created_at, created_by 
		)
		SELECT created_at, users.username
		FROM new_sample
		LEFT JOIN users ON new_sample.created_by = users.user_id
	`
	err := s.db.QueryRow(
		query,
		sample.SampleId,
		sample.FullName,
		sample.Sputalysed,
		sample.Comment,
		sample.Birthdate,
		userID,
		sample.Material,
	).Scan(&newSample.CreatedAt, &newSample.CreatedBy)
	if err != nil {
		return newSample, err
	}

	// Return the new sample
	return newSample, nil
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
		query += fmt.Sprintf(" full_name = $%d,", len(params)+1)
		params = append(params, *sample.FullName)
	}
	if sample.Sputalysed != nil {
		query += fmt.Sprintf(" sputalysed = $%d,", len(params)+1)
		params = append(params, *sample.Sputalysed)
	}
	if sample.Comment != nil {
		query += fmt.Sprintf(" comment = $%d,", len(params)+1)
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
		&updatedSample.SampleId,
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
