package printer

import (
	"fmt"
	"net"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"gitlab.com/kaka/pcr-backend/common/database"
)

var printerAddress = "bc-pcr2.labmed.de:9100"

var label = `q256
N
A150,40,0,5,1,1,N,"%s"
A30,0,0,2,1,1,N,"%s"
A30,20,0,2,1,1,N,"%s"
A30,50,0,2,1,1,N,"%s"
A30,70,0,2,1,1,N,"%s"
A30,100,0,1,1,1,N,"%s"
P1
`

type PrintRequest struct {
	SampleId string `json:"sample_id" binding:"required"`
	Panel    string `json:"panel" binding:"required"`
}

type PrintData struct {
	Position string `json:"position"`
	Name     string `json:"name"`
	SampleId string `json:"sample_id"`
	Panel    string `json:"panel"`
	Device   string `json:"device"`
	Run      string `json:"run"`
	Date     string `json:"date"`
}

func Print(ctx *gin.Context) {
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

	var printData PrintData
	printData.SampleId = request.SampleId
	printData.Panel = request.Panel

	// Query Sample to retrieve the inge name
	query := `
		SELECT position, run_date, samples.full_name, device, run
		FROM samplespanels
		LEFT JOIN samples ON samplespanels.sample_id = samples.sample_id
		WHERE samplespanels.sample_id = $1 AND panel_id = $2
		`
	err = database.Instance.QueryRow(query, request.SampleId, request.Panel).Scan(&printData.Position, &printData.Date, &printData.Name, &printData.Device, &printData.Run)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	// Format Date and Sample
	printData.Date = printData.Date[0:10]

	// Replace all non-ASCII characters with a question mark
	printString := fmt.Sprintf(label, printData.Position, printData.SampleId, printData.Name, printData.Panel, printData.Device+printData.Run, printData.Date)
	regex := regexp.MustCompile("[[:^ascii:]]")
	printString = regex.ReplaceAllString(printString, "?")

	// Send the label to the printer
	_, err = conn.Write([]byte(printString))

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}
