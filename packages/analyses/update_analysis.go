package analyses

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/kaka/pcr-backend/common/database"
)

type UpdateAnalysisRequest struct {
	ReadyMix *bool `json:"ready_mix"`
	IsActive *bool `json:"is_active"`
}

func UpdateAnalysis(ctx *gin.Context) {
	analysis_id := ctx.Param("analysis_id")

	var request UpdateAnalysisRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if request.ReadyMix != nil {
		err := updateReadyMix(ctx, analysis_id, request)

		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
	}

	if request.IsActive != nil {
		err := updateIsActive(ctx, analysis_id, request)

		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
	}

	ctx.Status(http.StatusOK)
}

func updateReadyMix(ctx *gin.Context, analysis_id string, request UpdateAnalysisRequest) error {
	_, err := database.Instance.Exec("UPDATE analyses SET ready_mix = $1 WHERE analysis_id = $2 RETURNING *", *request.ReadyMix, analysis_id)

	if err != nil {
		return err
	}

	return nil
}

func updateIsActive(ctx *gin.Context, analysis_id string, request UpdateAnalysisRequest) error {
	_, err := database.Instance.Exec("UPDATE analyses SET is_active = $1 WHERE analysis_id = $2", *request.IsActive, analysis_id)

	if err != nil {
		return err
	}

	return nil
}
