package repository

import (
	"database/sql"
	"log/slog"

	"gitlab.com/kallepan/pcr-backend/app/domain/dao"
	"gitlab.com/kallepan/pcr-backend/driver"
)

type PanelRepository interface {
	GetPanels(panelID string) ([]dao.Panel, error)
	PanelExists(panelID string) bool
}

type PanelRepositoryImpl struct {
	db *sql.DB
}

func PanelRepositoryInit() *PanelRepositoryImpl {
	return &PanelRepositoryImpl{
		db: driver.DB,
	}
}

func (p PanelRepositoryImpl) GetPanels(panelID string) ([]dao.Panel, error) {
	/* Returns all panels with an optional filter for panel_id */
	panels := []dao.Panel{}

	// Get all panels and filter out inactive ones
	query := `
		SELECT panel_id, display_name, ready_mix
		FROM panels
		WHERE is_active = TRUE
	`

	// Add pagination and filters
	var params []interface{}
	if panelID != "" {
		query += "AND panel_id = $1"
		params = append(params, panelID)
	}
	query += " ORDER BY display_name LIMIT 100;"
	slog.Info(query)

	// query
	rows, err := p.db.Query(query, params...)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var panel dao.Panel

		if err := rows.Scan(&panel.PanelID, &panel.DisplayName, &panel.ReadyMix); err != nil {
			break
		}
		panels = append(panels, panel)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return panels, nil
}

func (p PanelRepositoryImpl) PanelExists(panelID string) bool {
	/* Returns true if panel exists */
	query := `
			SELECT EXISTS(
				SELECT panel_id
				FROM panels
				WHERE panel_id = $1
			)
		`
	var exists bool
	err := p.db.QueryRow(query, panelID).Scan(&exists)
	if err != nil {
		slog.Error("Error checking if panel exists", err)
		return false
	}

	return exists
}
