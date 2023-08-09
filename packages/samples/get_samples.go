package samples

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/kaka/pcr-backend/common/database"
	"gitlab.com/kaka/pcr-backend/common/models"
)

func GetSamples(ctx *gin.Context) {
	var samples []models.Sample

	sample_id := ctx.Param("sample_id")

	var params []interface{}

	query := `
		SELECT s.sample_id,s.full_name,s.birthdate,s.sputalysed,s.comment,s.created_at,u.username, string_agg(samplespanels.panel_id, ', ') AS panels
		FROM samples s
		LEFT JOIN users u ON s.created_by = u.user_id
		LEFT JOIN samplespanels ON samplespanels.sample_id = s.sample_id AND samplespanels.deleted = false
		WHERE 1 = 1
	`

	if sample_id != "" {
		query += `
			AND s.sample_id LIKE $1
		`
		param := fmt.Sprintf("%%%s%%", sample_id)
		params = append(params, param)
	}
	query += `
		AND s.created_at >= current_date - interval '14 day'
		GROUP BY s.sample_id, u.username ORDER BY s.created_at DESC, s.sample_id DESC;
	`

	rows, err := database.Instance.Query(query, params...)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	for rows.Next() {
		var sample models.Sample
		var panels sql.NullString
		if err := rows.Scan(&sample.SampleId, &sample.FullName, &sample.Birthdate, &sample.Sputalysed, &sample.Comment, &sample.CreatedAt, &sample.CreatedBy, &panels); err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		if panels.Valid {
			sample.Panels = panels.String
		} else {
			sample.Panels = "N/A"
		}

		samples = append(samples, sample)
	}

	// Empty array
	if len(samples) == 0 {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	ctx.JSON(http.StatusOK, &samples)
}
