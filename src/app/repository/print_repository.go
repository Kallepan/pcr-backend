package repository

import (
	"database/sql"

	"gitlab.com/kallepan/pcr-backend/app/domain/dco"
	"gitlab.com/kallepan/pcr-backend/driver"
)

type PrintRepository interface {
	QuerySample(sampleID string, panelID string) (*dco.PrintData, error)
}

type PrintRepositoryImpl struct {
	db *sql.DB
}

func PrintRepositoryInit() *PrintRepositoryImpl {
	return &PrintRepositoryImpl{
		db: driver.DB,
	}
}

func (p PrintRepositoryImpl) QuerySample(sampleID string, panelID string) (*dco.PrintData, error) {
	/* Query sample data */
	var printData dco.PrintData

	var runDate sql.NullTime
	var run sql.NullString
	var device sql.NullString
	var fullName sql.NullString
	var position sql.NullString

	query := `
		SELECT position, run_date, samples.full_name, device, run
		FROM samplespanels
		LEFT JOIN samples ON samplespanels.sample_id = samples.sample_id
		WHERE samplespanels.sample_id = $1 AND panel_id = $2
	`
	err := p.db.QueryRow(query, sampleID, panelID).Scan(&position, &runDate, &fullName, &device, &run)
	if err != nil {
		return nil, err
	}

	// Parse the attributes
	if runDate.Valid {
		printData.Date = runDate.Time.Format("2006-01-02")
	} else {
		printData.Date = ""
	}
	if run.Valid {
		printData.Run = run.String
	} else {
		printData.Run = ""
	}
	if device.Valid {
		printData.Device = device.String
	} else {
		printData.Device = ""
	}
	if fullName.Valid {
		printData.Name = fullName.String
	} else {
		printData.Name = ""
	}
	if position.Valid {
		printData.Position = position.String
	} else {
		printData.Position = ""
	}

	return &printData, nil
}
