package samples

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/kaka/pcr-backend/common/database"
	"gitlab.com/kaka/pcr-backend/common/models"
)

type AddSampleRequest struct {
	SampleID  string `json:"sample_id" binding:"required"`
	FirstName string `json:"firstname,omitempty" binding:"required"`
	LastName  string `json:"lastname,omitempty" binding:"required"`
}

func AddSample(ctx *gin.Context) {
	user_id := ctx.MustGet("user_id").(string)

	var request AddSampleRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	sample := models.Sample{
		SampleID:  request.SampleID,
		FirstName: request.FirstName,
		LastName:  request.LastName,
	}

	// Check if sample already exists
	if SampleExists(sample.SampleID) {
		error_message := fmt.Sprintf("sample %s already exists", sample.SampleID)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": error_message})
		return
	}

	// Insert sample
	query := `
		WITH new_sample AS 
		(
			INSERT INTO samples (sample_id,firstname,lastname,created_by)
			VALUES ($1, $2, $3, $4) RETURNING *)
			SELECT created_at, users.username
			FROM new_sample
			LEFT JOIN users ON new_sample.created_by = users.user_id`
	err := database.Instance.QueryRow(query, sample.SampleID, sample.FirstName, sample.LastName, user_id).Scan(&sample.CreatedAt, &sample.CreatedBy)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, &sample)
}
