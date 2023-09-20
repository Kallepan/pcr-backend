package dco

import (
	"time"
)

type SamplePanelRequest struct {
	SamplePanel []SamplePanel `json:"samples" binding:"required"`
}

type SamplePanel struct {
	Material *string `json:"material" binding:"omitempty"`

	Name      string `json:"name" binding:"required"`
	SampleID  string `json:"sample_id" binding:"required,min=8,max=12,numeric"`
	Birthdate string `json:"birthdate" binding:"required"`
	PanelID   string `json:"panel_id" binding:"required,alpha"`
}

func (s *SamplePanel) Validate() error {
	/* Validate the sample panel */

	// Validate birthdate
	if _, err := time.Parse("2006-01-02", s.Birthdate); err != nil {
		return err
	}

	return nil
}
