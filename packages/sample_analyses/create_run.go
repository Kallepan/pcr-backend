package samplesanalyses

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"gitlab.com/kaka/pcr-backend/common/database"
	"gitlab.com/kaka/pcr-backend/common/models"
	"gitlab.com/kaka/pcr-backend/packages/analyses"
	"gitlab.com/kaka/pcr-backend/packages/samples"
)

type PostElementData struct {
	SampleID    *string `json:"sample_id" binding:"required"`
	AnalysisID  *string `json:"analysis_id" binding:"required"`
	ControlID   *string `json:"control_id" binding:"required"`
	Description *string `json:"description" binding:"required"`
}

type CreateRunRequest struct {
	Device       string            `json:"device" binding:"required"`
	Run          string            `json:"run" binding:"required"`
	PostElements []PostElementData `json:"elements" binding:"required"`
}

type ExportData struct {
	// Sample info
	sample *models.Sample
	// Analysis info
	analysis *models.Analysis
	// Control info
	Description *string

	// Run info
	Position *int
}

func createCopy(templatePath string) (string, error) {
	// Creates a copy of the template file to the tmp folder renaming it with a timestamp

	outputPath := fmt.Sprintf("tmp/%s.xlsx", time.Now().Format("20060102150405"))

	src, err := os.Open(templatePath)
	if err != nil {
		return "", err
	}
	defer src.Close()

	dst, err := os.Create(outputPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		return "", err
	}

	return outputPath, nil
}

func getFreePositionFromDatabase() (int, error) {
	// Get latest position from today from database
	freePosition := 0

	query := `
		SELECT position
		FROM samplesanalyses
		WHERE 
			DATE(created_at) = CURRENT_DATE AND
			position IS NOT NULL
		ORDER BY position DESC
		LIMIT 1;
		`

	row := database.Instance.QueryRow(query)
	err := row.Scan(&freePosition)
	switch {
	case err == sql.ErrNoRows:
		return 1, nil
	case err != nil:
		return 0, err
	default:
		return freePosition + 1, nil
	}
}

func UpdateSampleAnalysisInDatabase(sampleID string, analysisID string, position int, run string, device string) error {
	query := `
		UPDATE samplesanalyses
		SET
			position = $1,
			run = $2,
			device = $3
		WHERE
			sample_id = $4 AND
			analysis_id = $5;
		`

	_, err := database.Instance.Exec(query, position, run, device, sampleID, analysisID)
	return err
}

func CheckIfSampleAnalysisWasAlreadyUsed(sampleID string, analysisID string) error {
	query := `
		SELECT *
		FROM samplesanalyses
		WHERE
			sample_id = $1 AND
			analysis_id = $2 AND (
				position IS NOT NULL OR
				run IS NOT NULL OR
				device IS NOT NULL
			)
		`

	row := database.Instance.QueryRow(query, sampleID, analysisID)
	err := row.Scan()
	switch {
	case err == sql.ErrNoRows:
		return nil
	case err != nil:
		return err
	default:
		return fmt.Errorf("sample analysis was already used")
	}
}

func CreateRun(ctx *gin.Context) {
	var request CreateRunRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	// Load template
	templatePath := "templates/v1.xlsx"

	// Create copy of template
	outputPath, err := createCopy(templatePath)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	// Open copy of template
	file, err := excelize.OpenFile(outputPath)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	defer file.Close()
	defer os.Remove(outputPath)

	freePosition, err := getFreePositionFromDatabase()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	var exportData []ExportData
	// Fetch data from database for SampleAnalyses such as final position, last run, etc.
	for idx, postElement := range request.PostElements {
		if postElement.SampleID != nil && postElement.AnalysisID != nil {
			// SampleAnalysis
			// Check if SampleAnalysis was already used
			err := CheckIfSampleAnalysisWasAlreadyUsed(*postElement.SampleID, *postElement.AnalysisID)
			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("Sample: %s, Analysis: %s was already used", *postElement.SampleID, *postElement.AnalysisID)})
				return
			}

			newPosition := freePosition + idx
			// Create new element for exportData
			var exportDataElement ExportData

			// Fetch data from database
			sample, err := samples.FetchSampleInformationFromDatabase(*postElement.SampleID)
			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				return
			}
			analysis, err := analyses.FetchAnalysisInformationFromDatabase(*postElement.AnalysisID)
			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				return
			}
			exportDataElement.sample = sample
			exportDataElement.analysis = analysis
			exportDataElement.Description = nil
			// Update position
			exportDataElement.Position = &newPosition

			exportData = append(exportData, exportDataElement)
		} else if postElement.ControlID != nil && postElement.Description != nil {
			// Control
			// Create new element for exportData
			var exportDataElement ExportData
			exportDataElement.Description = postElement.Description
			exportDataElement.Position = nil

			exportData = append(exportData, exportDataElement)
		} else {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("%v is not valid data", postElement)})
			return
		}
	}

	// Update SampleAnalysis in the database
	for idx, element := range exportData {
		if element.Position == nil {
			// Control, we don't need to insert it into the database
			continue
		}
		// SampleAnalysis
		err = UpdateSampleAnalysisInDatabase(element.sample.SampleID, element.analysis.AnalysisID, *element.Position, request.Run, request.Device)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		// Insert data into excel file
		file.SetCellValue("Lauf", fmt.Sprintf("A%d", 12+idx), element.Position)
		// TODO: Insert rest of letters
	}

	// Insert data into excel file
	file.SetCellValue("Lauf", "C9", request.Device)
	file.SetCellValue("Lauf", "D9", request.Run)

	// Set Headers and status
	ctx.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	ctx.Header("Content-Disposition", "attachment; filename=run.xlsx")

	err = file.Write(ctx.Writer)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	err = file.Save()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	ctx.File(outputPath)
}
