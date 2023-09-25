package dco

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"gitlab.com/kallepan/pcr-backend/app/domain/dao"
)

type ExportData struct {
	// Sample info
	Sample dao.Sample
	// Analysis info
	Panel dao.Panel

	// Control info
	Description *string
	LastRunId   string

	IsControl bool
}

func (e *ExportData) GetFormattedComment() string {
	var comment string
	if e.Sample.Comment == nil {
		comment = "NA"
	} else {
		comment = *e.Sample.Comment
	}

	if e.Sample.Sputalysed {
		comment = fmt.Sprintf("Mit Sputasol; %s", comment)
	}

	return comment
}

func (e *ExportData) GetFormattedMaterial() string {
	if e.Sample.Material == nil {
		return "NA"
	}

	return *e.Sample.Material
}

func (e *ExportData) GetFormattedSampleID() string {
	// Format sample id
	if len(e.Sample.SampleID) == 8 {
		// XXXX XXXX
		return fmt.Sprintf("%s %s", e.Sample.SampleID[:4], e.Sample.SampleID[4:])
	} else if len(e.Sample.SampleID) == 12 {
		// XXXX XXXXXX XX
		return fmt.Sprintf("%s %s %s", e.Sample.SampleID[:4], e.Sample.SampleID[4:10], e.Sample.SampleID[10:])
	}

	return e.Sample.SampleID
}

func (e *ExportData) GetFormattedBirthdate() string {
	if e.Sample.Birthdate == nil {
		return "NA"
	}

	birthdate := *e.Sample.Birthdate
	birthdate = birthdate[:10]

	return birthdate
}

func (e *ExportData) GetFormattedName() string {
	// Get name, split by comma, take the last name and first letter of the first name
	// e.g. "Doe,John" -> "Doe J"
	// e.g. "Doe" -> "Doe"

	if e.Sample.FullName == "" {
		return ""
	}

	// Split by comma
	split := strings.Split(e.Sample.FullName, ",")
	if len(split) == 1 {
		return e.Sample.FullName
	}

	// Get last name
	lastName := split[0]

	// Get first name
	firstName := split[1]

	// Get first letter of first name
	firstLetter := string(firstName[0])

	return lastName + " " + firstLetter
}

type Element struct {
	SampleID    *string `json:"sample_id" binding:"required"`
	PanelID     *string `json:"panel_id" binding:"required"`
	ControlID   *string `json:"control_id" binding:"required"`
	Description *string `json:"description" binding:"required"`
}

type CreatRunRequest struct {
	Device string `json:"device" binding:"required"`
	Run    string `json:"run" binding:"required"`
	Date   string `json:"date" binding:"required"`

	Elements []Element `json:"elements" binding:"required"`
}

func (c *CreatRunRequest) Validate() error {
	_, err := time.Parse("2006-01-02", c.Date)
	if err != nil {
		return err
	}

	if len(c.Elements) == 0 {
		return errors.New("elements must not be empty")
	}

	return nil
}
