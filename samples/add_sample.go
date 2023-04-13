package samples

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gitlab.com/kaka/pcr-backend/common/database"
	"gitlab.com/kaka/pcr-backend/common/models"
)

type AddSampleRequest struct {
	Tagesnummer string `json:"tagesnummer"`
	Name        string `json:"name,omitempty" binding:"required"`
}

func AddSample(ctx *gin.Context) {
	user_id := ctx.MustGet("user_id").(string)

	var request AddSampleRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sample := models.Sample{
		Tagesnummer: request.Tagesnummer,
		Name:        request.Name,
		CreatedBy:   user_id,
	}

	query := "WITH new_sample AS (INSERT INTO samples (tagesnummer,name,created_by) VALUES ($1, $2, $3) RETURNING *) SELECT tagesnummer, name, created_at, users.username FROM new_sample LEFT JOIN users ON new_sample.created_by = users.user_id;"
	err := database.Instance.QueryRow(query, sample.Tagesnummer, sample.Name, sample.CreatedBy).Scan(&sample.Tagesnummer, &sample.Name, &sample.CreatedAt, &sample.CreatedBy)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			ctx.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": "sample already exists"})
			return
		}
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusCreated, &sample)
}
