package importer

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/kaka/pcr-backend/common/database"
	"gitlab.com/kaka/pcr-backend/packages/panels"
)

type SamplePanel struct {
	Material *string `json:"material" binding:"omitempty"`

	Name      string `json:"name" binding:"required"`
	SampleID  string `json:"sample_id" binding:"required,min=8,max=12,numeric"`
	Birthdate string `json:"birthdate" binding:"required"`
	PanelID   string `json:"panel_id" binding:"required,alpha"`
}

type SamplePanelRequest struct {
	SamplePanel []SamplePanel `json:"samples" binding:"required"`
}

func validateDate(s string) error {
	/* Validate date */

	_, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}

	return nil
}

func (s *SamplePanel) Validate() error {
	/* Validate the sample panel */

	// Validate birthdate
	if err := validateDate(s.Birthdate); err != nil {
		return err
	}

	if !panels.PanelExists(s.PanelID) {
		return ErrPanelNotFound
	}

	return nil
}

func (s *SamplePanelRequest) Validate() error {
	/* Validate the sample panel request */

	for _, samplePanel := range s.SamplePanel {
		if err := samplePanel.Validate(); err != nil {
			message := "Sample ID: " + samplePanel.SampleID + " - " + err.Error()
			return errors.New(message)
		}
	}

	return nil
}

func (s *SamplePanelRequest) Save() error {
	/* Save the sample panel request in the database */
	tx, err := database.Instance.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Println("Recovered in error:", r)
		}
	}()

	for _, samplePanel := range s.SamplePanel {
		query := `
		WITH inserted_sample AS (
			INSERT INTO samples (sample_id, full_name, birthdate, material, created_by)
			VALUES ($1, $2, $3, $4, $5)
			ON CONFLICT (sample_id) DO UPDATE SET
				full_name = $2,
				birthdate = $3,
				material = $4
			RETURNING sample_id
		)
		INSERT INTO samplespanels (sample_id, panel_id, created_by)
		VALUES ((SELECT sample_id FROM inserted_sample), $6, $5)
		ON CONFLICT (sample_id, panel_id) DO NOTHING
		`
		var material = "NA"
		if samplePanel.Material != nil {
			material = *samplePanel.Material
		}
		if _, err := tx.Exec(query, samplePanel.SampleID, samplePanel.Name, samplePanel.Birthdate, material, 1, samplePanel.PanelID); err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func PostSampleMaterial(ctx *gin.Context) {
	/* Post a sample material */

	var samplePanelRequest SamplePanelRequest
	if err := ctx.ShouldBindJSON(&samplePanelRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := samplePanelRequest.Validate(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := samplePanelRequest.Save(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	message := fmt.Sprintf("Saved %d sample(s)", len(samplePanelRequest.SamplePanel))
	ctx.JSON(http.StatusOK, gin.H{"message": message})
}
