package dco

import (
	"fmt"
	"regexp"
)

var GlobalTemplate = `q256
N
A150,40,0,5,1,1,N,"%s"
A30,0,0,2,1,1,N,"%s"
A30,20,0,2,1,1,N,"%s"
A30,50,0,2,1,1,N,"%s"
A30,70,0,2,1,1,N,"%s%s"
A30,100,0,1,1,1,N,"%s"
P1
`

type PrintRequestElement struct {
	SampleID string `json:"sample_id" binding:"required"`
	PanelID  string `json:"panel_id" binding:"required"`
}

type PrintRequest struct {
	Elements []PrintRequestElement `json:"elements" binding:"required"`
}

type PrintData struct {
	Position string
	Name     string
	SampleID string
	PanelID  string
	Device   string
	Run      string
	Date     string
}

func (pd PrintData) CreateLabel(template string) (string, error) {
	label := fmt.Sprintf(template, pd.Position, pd.SampleID, pd.Name, pd.PanelID, pd.Device, pd.Run, pd.Date)

	regex, err := regexp.Compile("[[:^ascii:]]")
	if err != nil {
		return "", err
	}

	label = regex.ReplaceAllString(label, "?")

	// Return the formatted label
	return label, nil
}
