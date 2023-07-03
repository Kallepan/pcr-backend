package panels

import (
	"database/sql"

	"gitlab.com/kaka/pcr-backend/common/database"
	"gitlab.com/kaka/pcr-backend/common/models"
)

func PanelExists(panelId string) bool {
	query := `
		SELECT EXISTS(
			SELECT panel_id
			FROM panels
			WHERE panel_id = $1
		)`

	var exists bool

	if err := database.Instance.QueryRow(query, panelId).Scan(&exists); err != nil && err != sql.ErrNoRows {
		return false
	}

	return exists
}

func FetchPanelInformationFromDatabase(panelId string) (*models.Panel, error) {
	query := `
		SELECT panel_id, display_name, ready_mix
		FROM panels
		WHERE panel_id = $1`

	var panel models.Panel

	if err := database.Instance.QueryRow(query, panelId).Scan(&panel.PanelId, &panel.DisplayName, &panel.ReadyMix); err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return &panel, nil
}
