package samplespanels

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/kaka/pcr-backend/common/database"
	"gitlab.com/kaka/pcr-backend/common/models"
)

func GetSamplesPanels(ctx *gin.Context) {
	var samplespanels []models.SampleAnalysis

	sample_id := ctx.Param("sample_id")

	var params []interface{}

	query := `
		WITH sample_query AS (
			SELECT samplespanels.sample_id, samples.full_name, samples.created_at, users.username AS created_by
			FROM samplespanels
			LEFT JOIN samples ON samplespanels.sample_id = samples.sample_id
			LEFT JOIN users ON samples.created_by = users.user_id
			GROUP BY samplespanels.sample_id, samples.full_name, samples.created_at, users.username
		) 
		SELECT samplespanels.sample_id, sample_query.full_name, sample_query.created_at, sample_query.created_by, samplespanels.panel_id, panels.display_name, panels.ready_mix, samplespanels.run, samplespanels.device, samplespanels.position, samplespanels.run_date, samplespanels.created_at, users.username
		FROM samplespanels
		LEFT JOIN sample_query ON samplespanels.sample_id = sample_query.sample_id
		LEFT JOIN panels ON samplespanels.panel_id = panels.panel_id
		LEFT JOIN users ON samplespanels.created_by = users.user_id
		WHERE
			1 = 1 AND
		`
	if sample_id != "" {
		query += "samplespanels.sample_id = $1 AND samplespanels.deleted = false"
		params = append(params, sample_id)
	} else {
		query += `
			samplespanels.deleted = false AND 
			samplespanels.run IS NULL AND
			samplespanels.device IS NULL AND
			samplespanels.position IS NULL
			ORDER BY samplespanels.created_at DESC, samplespanels.sample_id DESC`
	}
	rows, err := database.Instance.Query(query, params...)

	if err != nil {
		ctx.AbortWithStatusJSON(500, gin.H{"message": err.Error()})
		return
	}

	for rows.Next() {
		var sampleAnalysis models.SampleAnalysis
		var sample models.Sample
		var panel models.Panel

		if err := rows.Scan(
			&sample.SampleId, &sample.FullName, &sample.CreatedAt, &sample.CreatedBy,
			&panel.PanelId, &panel.DisplayName, &panel.ReadyMix,
			&sampleAnalysis.Run, &sampleAnalysis.Device, &sampleAnalysis.Position, &sampleAnalysis.RunDate, &sampleAnalysis.CreatedAt, &sampleAnalysis.CreatedBy); err != nil {

			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		}

		sampleAnalysis.Sample = sample
		sampleAnalysis.Panel = panel
		samplespanels = append(samplespanels, sampleAnalysis)
	}

	if err = rows.Err(); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	// Empty array
	if len(samplespanels) == 0 {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	ctx.JSON(200, &samplespanels)
}
