package repository

import (
	"database/sql"
	"log/slog"

	"gitlab.com/kallepan/pcr-backend/app/domain/dco"
	"gitlab.com/kallepan/pcr-backend/driver"
)

type ImportRepository interface {
	Save(sample []dco.SamplePanel) error
}

type ImportRepositoryImpl struct {
	db *sql.DB
}

func ImportRepositoryInit() *ImportRepositoryImpl {
	return &ImportRepositoryImpl{
		db: driver.DB,
	}
}

func (i ImportRepositoryImpl) Save(sample []dco.SamplePanel) error {
	tx, err := i.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			slog.Error("Recovered in error:", r)
		}
	}()

	for _, samplePanel := range sample {
		query := `
		WITH inserted_sample AS (
			INSERT INTO samples (sample_id, full_name, birthdate, material, created_by)
			VALUES ($1, $2, $3, $4, $5)
			ON CONFLICT (sample_id) DO UPDATE SET
				full_name = $2,
				birthdate = $3,
				material = $4
			RETURNING sample_id
		)
		INSERT INTO samplespanels (sample_id, panel_id, created_by)
		VALUES ((SELECT sample_id FROM inserted_sample), $6, $5)
		ON CONFLICT (sample_id, panel_id) DO NOTHING
		`

		// Check if material is nil
		material := "NA"
		if samplePanel.Material == nil {
			samplePanel.Material = &material
		}

		if _, err := tx.Exec(query, samplePanel.SampleID, samplePanel.Name, samplePanel.Birthdate, material, 1, samplePanel.PanelID); err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
