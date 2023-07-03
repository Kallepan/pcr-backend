package panels

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/kaka/pcr-backend/common/database"
	"gitlab.com/kaka/pcr-backend/common/models"
)

func GetPanels(ctx *gin.Context) {
	var panels []models.Panel

	panel_id := ctx.Param("panel_id")

	// Get all panels and filter out inactive ones
	query := `
		SELECT panel_id, display_name, ready_mix
		FROM panels
		WHERE is_active = TRUE
	`
	// Add pagination and filters
	var params []interface{}
	if panel_id != "" {
		query += "AND panel_id = $1"
		params = append(params, panel_id)
	}
	query += " ORDER BY display_name LIMIT 100;"

	// query
	rows, err := database.Instance.QueryContext(ctx, query, params...)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	for rows.Next() {
		var panel models.Panel

		if err := rows.Scan(&panel.PanelId, &panel.DisplayName, &panel.ReadyMix); err != nil {
			break
		}
		panels = append(panels, panel)
	}

	if err = rows.Err(); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, &panels)
}
