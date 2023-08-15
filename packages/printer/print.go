package printer

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"gitlab.com/kaka/pcr-backend/common/database"
)

var printerAddress = "bc-pcr2.labmed.de:9100"

var globalTemplate = `q256
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

func (pd PrintData) createLabel(template string) (string, error) {
	label := fmt.Sprintf(template, pd.Position, pd.SampleID, pd.Name, pd.PanelID, pd.Device, pd.Run, pd.Date)

	regex, err := regexp.Compile("[[:^ascii:]]")
	if err != nil {
		return "", err
	}

	label = regex.ReplaceAllString(label, "?")

	// Return the formatted label
	return label, nil
}

func queryElement(sampleID string, panelID string) (*PrintData, error) {
	var printData PrintData

	var runDate sql.NullTime
	var run sql.NullString
	var device sql.NullString
	var full_name sql.NullString
	var position sql.NullString

	query := `
		SELECT position, run_date, samples.full_name, device, run
		FROM samplespanels
		LEFT JOIN samples ON samplespanels.sample_id = samples.sample_id
		WHERE samplespanels.sample_id = $1 AND panel_id = $2
	`

	err := database.Instance.QueryRow(query, sampleID, panelID).Scan(&position, &runDate, &full_name, &device, &run)

	// Parse the attributes
	if runDate.Valid {
		printData.Date = runDate.Time.Format("2006-01-02")
	} else {
		printData.Date = ""
	}
	if run.Valid {
		printData.Run = run.String
	} else {
		printData.Run = ""
	}
	if device.Valid {
		printData.Device = device.String
	} else {
		printData.Device = ""
	}
	if full_name.Valid {
		printData.Name = full_name.String
	} else {
		printData.Name = ""
	}
	if position.Valid {
		printData.Position = position.String
	} else {
		printData.Position = ""
	}

	return &printData, err
}

func Print(ctx *gin.Context) {
	// Connect to the printer
	conn, err := net.Dial("tcp", printerAddress)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	defer conn.Close()

	// Extract the request body into a struct
	var request PrintRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// Query Sample to retrieve the inge name
	var labels []string

	for _, element := range request.Elements {
		printData, err := queryElement(element.SampleID, element.PanelID)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		// Set the sample ID and panel ID
		printData.SampleID = element.SampleID
		printData.PanelID = element.PanelID

		// Generate the label string
		label, err := printData.createLabel(globalTemplate)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		// Append the label to the slice
		labels = append(labels, label)
	}

	// If no error occurred, send the labels to the printer
	for i, label := range labels {
		log.Println("Printing label", i+1)

		_, err = conn.Write([]byte(label))
		if err != nil {
			log.Println(err.Error())
			// Do not return here, because we want to print as many labels as possible
			continue
		}
	}

	ctx.Status(http.StatusOK)
}
