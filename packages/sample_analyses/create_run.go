package samplesanalyses

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"gitlab.com/kaka/pcr-backend/common/database"
)

type SampleAnalysisBundle struct {
	SampleID   string `json:"sample_id"`
	AnalysisID string `json:"analysis_id"`
}

type CreateRunRequest struct {
	Device          string                 `json:"device" binding:"required"`
	Run             string                 `json:"run" binding:"required"`
	SamplesAnalyses []SampleAnalysisBundle `json:"samplesAnalyses" binding:"required"`
}

func CreateRunAlternative(ctx *gin.Context) {
	// Dont crate an excel file, just return update the database
	var request CreateRunRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// Iterate over samples and analysis and update in database
	for idx, sampleAnalysis := range request.SamplesAnalyses {
		query := `
			UPDATE samplesanalyses device = $1, run = $2, position = $3
			WHERE sample_id = $4 AND analysis_id = $5;`
		_, err := database.Instance.Exec(query, request.Device, request.Run, idx, sampleAnalysis.SampleID, sampleAnalysis.AnalysisID)

		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
	}

	// Return success
	ctx.Status(http.StatusOK)
}

func CreateRun(ctx *gin.Context) {
	var request CreateRunRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// Load template
	templatePath := "templates/v1.xlsx"
	file, err := excelize.OpenFile(templatePath)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	file.SetCellValue("Lauf", "C12", request.Device)
	file.SetCellValue("Lauf", "B12", request.Run)

	// Set Headers and status
	ctx.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	ctx.Header("Content-Disposition", "attachment; filename=run.xlsx")
	ctx.Status(http.StatusOK)

	err = file.Write(ctx.Writer)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
}
