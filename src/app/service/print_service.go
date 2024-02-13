package service

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net"

	"github.com/gin-gonic/gin"
	"gitlab.com/kallepan/pcr-backend/app/constant"
	"gitlab.com/kallepan/pcr-backend/app/domain/dco"
	"gitlab.com/kallepan/pcr-backend/app/pkg"
	"gitlab.com/kallepan/pcr-backend/app/repository"
)

var printerAddress = "bc-pcr2.labmed.de:9100"

type PrintService interface {
	PrintSample(ctx *gin.Context)
}

type PrintServiceImpl struct {
	printRepository repository.PrintRepository
}

func PrintServiceInit(printRepository repository.PrintRepository) *PrintServiceImpl {
	return &PrintServiceImpl{
		printRepository: printRepository,
	}
}

func (p PrintServiceImpl) PrintSample(ctx *gin.Context) {
	defer pkg.PanicHandler(ctx)
	// Connect to the printer
	conn, err := net.Dial("tcp", printerAddress)
	if err != nil {
		slog.Error("Error connecting to printer", err)
		pkg.PanicException(constant.UnknownError)
	}
	defer conn.Close()

	// Extract the request body into a struct
	var request dco.PrintRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		slog.Error("Error binding JSON", err)
		pkg.PanicException(constant.InvalidRequest)
	}

	// Query Sample to retrieve the inge name
	var labels []string
	for _, element := range request.Elements {
		printData, err := p.printRepository.QuerySample(element.SampleID, element.PanelID)
		if err == sql.ErrNoRows {
			message := fmt.Sprintf("Sample %s with panel %s not found", element.SampleID, element.PanelID)
			slog.Error(message)
			pkg.PanicExceptionWithMessage(constant.InvalidRequest, message)
		}
		if err != nil {
			slog.Error("Error querying sample", err)
			pkg.PanicException(constant.UnknownError)
		}

		// Set the sample ID and panel ID
		printData.SampleID = element.SampleID
		printData.PanelID = element.PanelID

		// Generate the label string
		label, err := printData.CreateLabel(dco.GlobalTemplate)
		if err != nil {
			slog.Error("Error creating label", err)
			pkg.PanicException(constant.UnknownError)
		}

		// Append the label to the slice
		labels = append(labels, label)
	}

	// If no error occurred, send the labels to the printer
	for i, label := range labels {
		slog.Info("Printing label", "id", i+1)
		if _, err := fmt.Fprint(conn, label); err != nil {
			slog.Error("Error sending label to printer", err)
			// Error does not need to be handled, we want to print as many labels as possible
		}
	}

	ctx.JSON(200, pkg.BuildResponse(constant.Success, pkg.Null()))
}
