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

	IsControl bool
}

func createCopy(templatePath string) (string, error) {
	// Creates a copy of the template file to the tmp folder renaming it with a timestamp

	outputPath := fmt.Sprintf("tmp/%s.xlsm", time.Now().Format("20060102150405"))

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

func UpdateSampleAnalysisInDatabase(sampleID string, analysisID string, run string, device string) error {
	query := `
		UPDATE samplesanalyses
		SET
			run = $1,
			device = $2
		WHERE
			sample_id = $3 AND
			analysis_id = $4;
		`

	_, err := database.Instance.Exec(query, run, device, sampleID, analysisID)
	return err
}

func GetPositionForSampleAnalysis(tx *sql.Tx, sampleID string, analysisID string) (*int, error) {
	// Fetch position for sample analysis
	var position int
	if err := tx.QueryRow(`
		SELECT position 
		FROM samplesanalyses 
		WHERE 
			sample_id = $1 AND
			analysis_id = $2
		`, sampleID, analysisID).Scan(&position); err != nil {
		return nil, err
	}

	return &position, nil
}

func CheckIfSampleAnalysisIsInRun(sampleID string, analysisID string) error {
	query := `
		SELECT *
		FROM samplesanalyses
		WHERE
			sample_id = $1 AND
			analysis_id = $2 AND (
				run IS NOT NULL AND
				device IS NOT NULL
			)
	`
	err := database.Instance.QueryRow(query, sampleID, analysisID).Scan()
	switch err {
	case sql.ErrNoRows:
		return nil
	default:
		return fmt.Errorf("sample analysis was already used")
	}
}

func CreateRun(ctx *gin.Context) {
	// Create transaction to rollback if something goes wrong
	tx, err := database.Instance.Begin()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": fmt.Sprintf("%s", r)})
		}
	}()

	var request CreateRunRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	// Load template
	templatePath := "/app/templates/v1.xlsm"

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
	// Close file and remove it from disk
	defer func() {
		file.Close()
		os.Remove(outputPath)
	}()

	var exportData []ExportData
	// Validate data in a separate loop to avoid partial data being inserted
	for _, postElement := range request.PostElements {
		if postElement.SampleID != nil && postElement.AnalysisID != nil {
			// SampleAnalysis
			// Check if SampleAnalysis was already used
			err := CheckIfSampleAnalysisIsInRun(*postElement.SampleID, *postElement.AnalysisID)
			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("Sample: %s, Analysis: %s was already used", *postElement.SampleID, *postElement.AnalysisID)})
				return
			}

			// Create new element for exportData
			var exportDataElement ExportData

			// Fetch data from database
			sample, err := samples.FetchSampleInformationFromDatabase(*postElement.SampleID)
			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				tx.Rollback()
				return
			}
			analysis, err := analyses.FetchAnalysisInformationFromDatabase(*postElement.AnalysisID)
			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				tx.Rollback()
				return
			}
			exportDataElement.sample = sample
			exportDataElement.analysis = analysis
			exportDataElement.IsControl = false
			// Append description --> last occurence of sample in a run
			exportData = append(exportData, exportDataElement)
		} else if postElement.ControlID != nil && postElement.Description != nil {
			// Control
			var exportDataElement ExportData
			exportDataElement.Description = postElement.Description
			exportDataElement.IsControl = true
			exportData = append(exportData, exportDataElement)
		} else {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("%v is not valid data", postElement)})
			tx.Rollback()
			return
		}
	}

	// Insert data into database and excel file
	for idx, exportDataElement := range exportData {
		if !exportDataElement.IsControl {
			// SampleAnalysis
			// Insert data into database -> position is auto incremented in the database
			if err := UpdateSampleAnalysisInDatabase(exportDataElement.sample.SampleID, exportDataElement.analysis.AnalysisID, request.Run, request.Device); err != nil {
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				tx.Rollback()
				return
			}

			// Get position from database
			position, err := GetPositionForSampleAnalysis(tx, exportDataElement.sample.SampleID, exportDataElement.analysis.AnalysisID)
			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				tx.Rollback()
				return
			}

			// Insert data into excel file
			file.SetCellValue(
				"Lauf",
				fmt.Sprintf("A%d", idx+12),
				*position,
			)
			file.SetCellValue(
				"Lauf",
				fmt.Sprintf("B%d", idx+12),
				fmt.Sprintf("%s, %s", exportDataElement.sample.SampleID, exportDataElement.sample.FullName),
			)
			file.SetCellValue(
				"Lauf",
				fmt.Sprintf("C%d", idx+12),
				fmt.Sprintf("%s-%s-%s", exportDataElement.analysis.Analyt, exportDataElement.analysis.Material, exportDataElement.analysis.Assay),
			)
			file.SetCellValue(
				"Lauf",
				fmt.Sprintf("E%d", idx+12),
				exportDataElement.sample.Comment,
			)
			if exportDataElement.sample.Sputalysed {
				file.SetCellValue(
					"Lauf",
					fmt.Sprintf("F%d", idx+12),
					"X",
				)
			}
		} else {
			// Control
			// Insert data into excel file
			file.SetCellValue(
				"Lauf",
				fmt.Sprintf("A%d", idx+12),
				"NA",
			)
			file.SetCellValue(
				"Lauf",
				fmt.Sprintf("D%d", idx+12),
				*exportDataElement.Description,
			)
		}
	}

	// Insert metadata into excel file
	file.SetCellValue("Lauf", "B9", time.Now().Format("02.01.2006"))
	file.SetCellValue("Lauf", "C9", request.Device)
	file.SetCellValue("Lauf", "D9", request.Run)

	// Set Headers and status
	ctx.Header("Content-Type", "application/vnd.ms-excel.sheet.macroEnabled.12")
	ctx.Header("Content-Disposition", "attachment; filename=run.xlsm")

	if err := file.Write(ctx.Writer); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		tx.Rollback()
		return
	}
	if err := file.Save(); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		tx.Rollback()
		return
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	ctx.File(outputPath)
}
