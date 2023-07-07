package samplespanels

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"gitlab.com/kaka/pcr-backend/common/database"
	"gitlab.com/kaka/pcr-backend/common/models"
	"gitlab.com/kaka/pcr-backend/packages/panels"
	"gitlab.com/kaka/pcr-backend/packages/samples"
)

type PostElementData struct {
	SampleId    *string `json:"sample_id" binding:"required"`
	PanelId     *string `json:"panel_id" binding:"required"`
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
	panel *models.Panel
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

func UpdateSampleAnalysisInDatabase(tx *sql.Tx, sampleId string, panelId string, run string, device string) error {
	query := `
		UPDATE samplespanels
		SET
			run = $1,
			device = $2
		WHERE
			sample_id = $3 AND
			panel_id = $4;
		`

	_, err := tx.Exec(query, run, device, sampleId, panelId)
	return err
}

func GetPositionForSampleAnalysis(tx *sql.Tx, sampleId string, panelId string) (*int, error) {
	// Fetch position for sample analysis
	var position int
	if err := tx.QueryRow(`
		SELECT position 
		FROM samplespanels 
		WHERE 
			sample_id = $1 AND
			panel_id = $2
		`, sampleId, panelId).Scan(&position); err != nil {
		return nil, err
	}

	return &position, nil
}

func CheckIfSampleAnalysisIsInRun(sampleId string, panelId string) error {
	query := `
		SELECT *
		FROM samplespanels
		WHERE
			sample_id = $1 AND
			panel_id = $2 AND (
				run IS NOT NULL AND
				device IS NOT NULL
			)
	`
	err := database.Instance.QueryRow(query, sampleId, panelId).Scan()
	switch err {
	case sql.ErrNoRows:
		return nil
	default:
		return fmt.Errorf("sample panel combination is already present")
	}
}

var createRunLock sync.Mutex

func CreateRun(ctx *gin.Context) {
	createRunLock.Lock()
	defer createRunLock.Unlock()
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
	// Load template TODO: move to config
	//templatePath := "/app/templates/v1.xlsm"
	templatePath := "templates/v1.xlsm"
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
		if postElement.SampleId != nil && postElement.PanelId != nil {
			// SampleAnalysis
			// Check if SampleAnalysis was already used
			err := CheckIfSampleAnalysisIsInRun(*postElement.SampleId, *postElement.PanelId)
			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("Sample: %s, Panel/Analysis: %s was already used", *postElement.SampleId, *postElement.PanelId)})
				tx.Rollback()
				return
			}

			// Create new element for exportData
			var exportDataElement ExportData

			// Fetch data from database
			sample, err := samples.FetchSampleInformationFromDatabase(*postElement.SampleId)
			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				tx.Rollback()
				return
			}
			panel, err := panels.FetchPanelInformationFromDatabase(*postElement.PanelId)
			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				tx.Rollback()
				return
			}
			exportDataElement.sample = sample
			exportDataElement.panel = panel
			exportDataElement.IsControl = false
			// Append description --> last occurence of sample in a run
			exportData = append(exportData, exportDataElement)
		} else if postElement.ControlID != nil && postElement.Description != nil {
			// Handle Controls
			var exportDataElement ExportData
			exportDataElement.Description = postElement.Description
			exportDataElement.IsControl = true
			exportData = append(exportData, exportDataElement)
		} else {
			// Invalid data
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
			if err := UpdateSampleAnalysisInDatabase(tx, exportDataElement.sample.SampleId, exportDataElement.panel.PanelId, request.Run, request.Device); err != nil {
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				tx.Rollback()
				return
			}

			// Get position from database
			position, err := GetPositionForSampleAnalysis(tx, exportDataElement.sample.SampleId, exportDataElement.panel.PanelId)
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
				fmt.Sprintf("%s, %s - %s", exportDataElement.sample.SampleId, exportDataElement.sample.FullName, exportDataElement.sample.Birthdate),
			)
			file.SetCellValue(
				"Lauf",
				fmt.Sprintf("C%d", idx+12),
				exportDataElement.panel.DisplayName,
			)
			// Check if comment is not nil
			if exportDataElement.sample.Comment != nil {
				file.SetCellValue(
					"Lauf",
					fmt.Sprintf("D%d", idx+12),
					*exportDataElement.sample.Comment,
				)
			}
			// Check if sample is sputalysed
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
