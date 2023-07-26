package printer

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"gitlab.com/kaka/pcr-backend/common/database"
)

var printerAddress = "bc-pcr2.labmed.de:9100"

var template = `q256
N
A150,40,0,5,1,1,N,"%s"
A30,0,0,2,1,1,N,"%s"
A30,20,0,2,1,1,N,"%s"
A30,50,0,2,1,1,N,"%s"
A30,70,0,2,1,1,N,"%s"
A30,100,0,1,1,1,N,"%s"
P1`

type PrintRequestElement struct {
	SampleId string `json:"sample_id" binding:"required"`
	PanelId  string `json:"panel_id" binding:"required"`
}

type PrintRequest struct {
	Elements []PrintRequestElement `json:"elements" binding:"required"`
}

type PrintData struct {
	Position string
	Name     string
	SampleId string
	Panel    string
	Device   string
	Run      string
	Date     string
}

func (pd PrintData) createLabel() string {
	// Truncate the date to 10 characters
	pd.Date = pd.Date[0:10]
	label := fmt.Sprintf(template, pd.Position, pd.SampleId, pd.Name, pd.Panel, pd.Device+pd.Run, pd.Date)
	regex := regexp.MustCompile("[[:^ascii:]]")
	label = regex.ReplaceAllString(label, "?")

	// Return the formatted label
	return label
}

func queryElement(sampleID string, panelID string) (*PrintData, error) {
	var printData PrintData

	query := `
		SELECT position, run_date, samples.full_name, device, run
		FROM samplespanels
		LEFT JOIN samples ON samplespanels.sample_id = samples.sample_id
		WHERE samplespanels.sample_id = $1 AND panel_id = $2
	`
	err := database.Instance.QueryRow(query, sampleID, panelID).Scan(&printData.Position, &printData.Date, &printData.Name, &printData.Device, &printData.Run)

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
		printData, err := queryElement(element.SampleId, element.PanelId)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		// Generate the label string
		label := printData.createLabel()

		// Append the label to the slice
		labels = append(labels, label)
	}

	// If no error occurred, send the labels to the printer
	for _, label := range labels {
		_, err = conn.Write([]byte(label))
		if err != nil {
			log.Println(err.Error())
			// Do not return here, because we want to print as many labels as possible
			continue
		}
	}

	ctx.Status(http.StatusOK)
}
