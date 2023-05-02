package samples

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/kaka/pcr-backend/common/database"
	"gitlab.com/kaka/pcr-backend/common/models"
)

type UpdateSampleRequest struct {
	FirstName  string `json:"firstname" binding:"required"`
	LastName   string `json:"lastname" binding:"required"`
	Sputalysed *bool  `json:"sputalysed" binding:"required"`
}

func UpdateSample(ctx *gin.Context) {
	sample_id := ctx.Param("sample_id")
	body := UpdateSampleRequest{}

	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// Check if sample exists
	if !SampleExists(sample_id) {
		error_message := fmt.Sprintf("sample %s does not exist", sample_id)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": error_message})
		return
	}

	query := `
		WITH updated_sample as (UPDATE samples SET firstname = $1, lastname = $2, sputalysed = $3 WHERE sample_id = $4 returning *) 
		SELECT sample_id, updated_sample.firstname, updated_sample.lastname, updated_sample.created_at, updated_sample.sputalysed , users.username 
		FROM updated_sample 
		LEFT JOIN users ON updated_sample.created_by = users.user_id;`

	result := database.Instance.QueryRow(query, body.FirstName, body.LastName, body.Sputalysed, sample_id)

	var sample models.Sample

	switch err := result.Scan(&sample.SampleID, &sample.FirstName, &sample.LastName, &sample.CreatedAt, &sample.Sputalysed, &sample.CreatedBy); err {
	case nil:
		break
	default:
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, &sample)
}
