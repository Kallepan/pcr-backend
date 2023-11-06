package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/xuri/excelize/v2"
	"gitlab.com/kallepan/pcr-backend/app/domain/dco"
	"gitlab.com/kallepan/pcr-backend/driver"
)

type RunRepository interface {
	CreateRun(request []dco.ExportData, device string, run string, date string) (string, error)

	// Utility functions
	UpdateSamplePanelInDatabase(sampleID string, panelID string, run string, device string, date string) (int, error)
	GetPositionForSampleAnalysis(sampleID string, panelID string) (int, error)
	IsAlreadyInRun(sampleID string, panelID string) bool
	GetLastRunId(sampleID string) (string, error)
}

type RunRepositoryImpl struct {
	db *sql.DB
}

func RunRepositoryInit() *RunRepositoryImpl {
	return &RunRepositoryImpl{
		db: driver.DB,
	}
}

// Utility functions
func (r RunRepositoryImpl) GetLastRunId(sampleID string) (string, error) {
	query := `
		SELECT CONCAT(device, '-POS', position, '-', run_date)
		FROM samplespanels
		WHERE
			sample_id = $1 AND
			run_date IS NOT NULL AND
			device IS NOT NULL AND
			position IS NOT NULL
		ORDER BY run_date ASC
		LIMIT 1;`

	var runID sql.NullString
	if err := r.db.QueryRow(query, sampleID).Scan(&runID); err == sql.ErrNoRows {
		return "", nil
	} else if err != nil {
		return "", err
	}

	if runID.Valid {
		return runID.String, nil
	}

	return "", errors.New("runID is null")
}

func (r RunRepositoryImpl) UpdateSamplePanelInDatabase(sampleID string, panelID string, run string, device string, date string) (int, error) {
	var position sql.NullInt32
	query := `
		UPDATE samplespanels
		SET
			run_date = $1,
			device = $2,
			run = $3
		WHERE
			sample_id = $4 AND
			panel_id = $5
		RETURNING position
	`

	if err := r.db.QueryRow(query, date, device, run, sampleID, panelID).Scan(&position); err != nil {
		return 0, err
	}

	if position.Valid {
		return int(position.Int32), nil
	}

	return 0, errors.New("position is null")
}

func (r RunRepositoryImpl) GetPositionForSampleAnalysis(sampleID string, panelID string) (int, error) {
	// Fetch position for sample analysis
	var nullPosition sql.NullInt32
	query := `
		SELECT position
		FROM samplespanels
		WHERE
			sample_id = $1 AND
			panel_id = $2 AND (
				run IS NOT NULL AND
				device IS NOT NULL
			)	
	`
	if err := r.db.QueryRow(query, sampleID, panelID).Scan(&nullPosition); err != nil {
		return 0, err
	}
	if nullPosition.Valid {
		return int(nullPosition.Int32), nil
	} else {
		return 0, errors.New("position is null")
	}
}

func (r RunRepositoryImpl) IsAlreadyInRun(sampleID string, panelID string) bool {
	var exists bool
	query := `
		SELECT EXISTS(
			SELECT sample_id
			FROM samplespanels
			WHERE
				sample_id = $1 AND
				panel_id = $2 AND (
					run IS NOT NULL AND
					device IS NOT NULL
				)
		)`

	if err := r.db.QueryRow(query, sampleID, panelID).Scan(&exists); err != nil {
		return false
	}

	return exists
}

// Core Functions
func (r RunRepositoryImpl) CreateRun(elements []dco.ExportData, device string, run string, date string) (string, error) {
	// Create transaction to rollback if something goes wrong
	tx, err := r.db.Begin()
	if err != nil {
		return "", err
	}

	// Recovery
	defer func() {
		if r := recover(); r != nil {
			slog.Error("Error creating run. Recovering...", r)
			tx.Rollback()
			panic(r)
		}
	}()

	// Load template path from env
	templatePath := os.Getenv("TEMPLATE_PATH")
	if templatePath == "" {
		templatePath = "/app/templates/v1.xlsm"
	}
	outputPath := fmt.Sprintf("/tmp/%s.xlsm", time.Now().Format("20060102150405"))

	// Check if template exists
	if _, err := os.Stat(templatePath); os.IsNotExist(err) {
		slog.Error("Template does not exist", err)
		tx.Rollback()
		return "", err
	}

	// Open copy of template
	file, err := excelize.OpenFile(templatePath)
	if err != nil {
		return "", err
	}
	// Close file
	defer func() {
		file.Close()
	}()

	// Iterate over elements
	for idx, element := range elements {
		if element.IsControl {
			// If a control set the control id and description
			file.SetCellValue(
				"Lauf",
				fmt.Sprintf("A%d", idx+12),
				"NA",
			)
			file.SetCellValue(
				"Lauf",
				fmt.Sprintf("D%d", idx+12),
				*element.Description,
			)
			continue
		}

		// SamplePanel
		// Update sample panel in database
		position, err := r.UpdateSamplePanelInDatabase(element.Sample.SampleID, element.Panel.PanelID, run, device, date)
		if err != nil {
			tx.Rollback()
			return "", err
		}

		// Insert sample panel into excel
		file.SetCellValue(
			"Lauf",
			fmt.Sprintf("A%d", idx+12),
			position,
		)

		name := element.GetFormattedName()
		material := element.GetFormattedMaterial()
		birthdate := element.GetFormattedBirthdate()
		sampleID := element.GetFormattedSampleID()
		comment := element.GetFormattedComment()

		// Insert name, material, birthdate into excel
		file.SetCellValue(
			"Lauf",
			fmt.Sprintf("B%d", idx+12),
			fmt.Sprintf("%s, %s - %s", sampleID, name, birthdate),
		)
		file.SetCellValue(
			"Lauf",
			fmt.Sprintf("C%d", idx+12),
			fmt.Sprintf("%s (%s)", element.Panel.DisplayName, material),
		)
		file.SetCellValue(
			"Lauf",
			fmt.Sprintf("E%d", idx+12),
			comment,
		)
		// LastRunId
		file.SetCellValue(
			"Lauf",
			fmt.Sprintf("D%d", idx+12),
			element.LastRunId,
		)
	}
	// Insert date, device, run into excel
	file.SetCellValue(
		"Lauf",
		"B9",
		date,
	)
	file.SetCellValue(
		"Lauf",
		"C9",
		device,
	)
	file.SetCellValue(
		"Lauf",
		"D9",
		run,
	)

	// Save file
	if err := file.SaveAs(outputPath); err != nil {
		slog.Error("Error saving file", err)
		tx.Rollback()
		return "", err
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		slog.Error("Error committing transaction", err)
		return "", err
	}

	return outputPath, nil
}
