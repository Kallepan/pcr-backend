package samples

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/kaka/pcr-backend/common/database"
	"gitlab.com/kaka/pcr-backend/common/models"
)

type AddSampleRequest struct {
	SampleId   string `json:"sample_id" binding:"required"`
	FullName   string `json:"full_name" binding:"required"`
	Sputalysed bool   `json:"sputalysed"`
	Comment    string `json:"comment,omitempty"`
	Birthdate  string `json:"birthdate" binding:"required"`
}

func AddSample(ctx *gin.Context) {
	user_id := ctx.MustGet("user_id").(string)

	var request AddSampleRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	sputalysed := request.Sputalysed || false

	sample := models.Sample{
		SampleId:   request.SampleId,
		FullName:   request.FullName,
		Sputalysed: sputalysed,
		Comment:    &request.Comment,
		Birthdate:  &request.Birthdate,
	}

	// Check if sample already exists
	if SampleExists(sample.SampleId) {
		error_message := fmt.Sprintf("Sample %s already exists", sample.SampleId)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": error_message})
		return
	}

	// Insert sample
	query := `
		WITH new_sample AS 
		(
			INSERT INTO samples (sample_id,full_name,sputalysed,comment,birthdate,created_by)
			VALUES ($1, $2, $3, $4, $5, $6) RETURNING created_at, created_by)
			SELECT created_at, users.username
			FROM new_sample
			LEFT JOIN users ON new_sample.created_by = users.user_id`
	err := database.Instance.QueryRow(query, sample.SampleId, sample.FullName, sample.Sputalysed, sample.Comment, sample.Birthdate, user_id).Scan(&sample.CreatedAt, &sample.CreatedBy)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, &sample)
}
