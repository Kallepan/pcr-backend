package samplespanels

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/kaka/pcr-backend/common/database"
)

type Statistic struct {
	PanelID string `json:"panel_id"`
	Count   int    `json:"count"`
}

func buildStatisticsQuery() string {
	return `
	SELECT
		LEFT(panel_id,3),
		count(*)
	FROM
		samplespanels
	WHERE
		run_date IS NULL AND deleted = FALSE
	GROUP BY
		LEFT(panel_id,3)
	`
}

func getStats(ctx *gin.Context) {
	/*
		1. Get all samplespanels that have a run_date of NULL
		2. Group by panel_id
		3. Count the number of rows in each group
		4. Return the panel_id and the count
	*/

	query := buildStatisticsQuery()

	rows, err := database.Instance.Query(query)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	var statistics []Statistic
	for rows.Next() {
		var statistic Statistic
		err := rows.Scan(&statistic.PanelID, &statistic.Count)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		statistics = append(statistics, statistic)
	}

	ctx.JSON(http.StatusOK, statistics)
}
