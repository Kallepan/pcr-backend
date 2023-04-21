package samples

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/kaka/pcr-backend/common/database"
	"gitlab.com/kaka/pcr-backend/common/models"
)

type UpdateSampleRequest struct {
	Name string `json:"name" binding:"required"`
}

func UpdateSample(ctx *gin.Context) {
	tagesnummer := ctx.Param("tagesnummer")
	body := UpdateSampleRequest{}

	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := `
		WITH updated_sample as (UPDATE samples SET name = $1 WHERE tagesnummer = $2 returning *) 
		SELECT tagesnummer, name, users.username 
		FROM updated_sample 
		LEFT JOIN users ON updated_sample.created_by = users.user_id;`

	result := database.Instance.QueryRow(query, body.Name, tagesnummer)

	var sample models.Sample

	switch err := result.Scan(&sample.Tagesnummer, &sample.Name, &sample.CreatedBy); err {
	case sql.ErrNoRows:
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "sample not found"})
		return
	case nil:
		ctx.JSON(http.StatusOK, &sample)
		return
	default:
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}
