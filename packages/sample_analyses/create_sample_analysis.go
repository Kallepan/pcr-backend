package samplesanalyses

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/kaka/pcr-backend/common/database"
	"gitlab.com/kaka/pcr-backend/common/models"
	"gitlab.com/kaka/pcr-backend/packages/analyses"
	"gitlab.com/kaka/pcr-backend/packages/samples"
)

type AddAnalysisToSampleRequest struct {
	SampleID   string            `json:"sample_id" binding:"required"`
	AnalysisID string            `json:"analysis_id" binding:"required"`
	Run        models.NullString `json:"run,omitempty"`
	Device     models.NullString `json:"device,omitempty"`
	Position   models.NullInt64  `json:"position,omitempty"`
}

func AddAnalysisToSample(ctx *gin.Context) {
	user_id := ctx.MustGet("user_id").(string)

	body := AddAnalysisToSampleRequest{}

	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// Check if analysis exists
	analysis, anlysis_err := analyses.AnalysisExistsByID(body.AnalysisID)
	if anlysis_err != nil {
		error_message := fmt.Sprintf("analysis %s does not exist", body.AnalysisID)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": error_message})
		return
	}

	// Check if sample exists
	if !samples.SampleExists(body.SampleID) {
		error_message := fmt.Sprintf("sample %s does not exist", body.SampleID)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": error_message})
		return
	}

	// Check if sample analysis already exists
	if SampleAnalysisExists(body.SampleID, body.AnalysisID) {
		error_message := fmt.Sprintf("Probe mit Analyse %s %s-%s-%s existiert bereits", body.SampleID, analysis.Analyt, analysis.Material, analysis.Assay)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": error_message})
		return
	}

	// Create sample analysis
	sample_analysis := models.SampleAnalysis{}
	sample_analysis.Sample = models.Sample{}
	sample_analysis.Analysis = models.Analysis{}
	sample_analysis.Sample.SampleID = body.SampleID
	sample_analysis.Analysis.AnalysisID = body.AnalysisID
	sample_analysis.Run = body.Run
	sample_analysis.Device = body.Device
	sample_analysis.Position = body.Position

	// Run query
	query := `
		WITH sample_query AS (
		SELECT samples.sample_id, samples.firstname, samples.lastname, samples.created_at, samples.created_by, samples.sputalysed, users.username AS sample_created_by
		FROM samples
		LEFT JOIN users ON samples.created_by = users.user_id
		WHERE samples.sample_id = $1
		GROUP BY samples.sample_id, samples.firstname, samples.lastname, samples.created_at, samples.sputalysed, users.username
		), new_sample_analysis AS ( 
			INSERT INTO samplesanalyses (sample_id, analysis_id, run, device, created_by) 
			VALUES ($1, $2, $3, $4, $5)
			RETURNING sample_id, created_at, created_by, analysis_id
		)
			SELECT new_sample_analysis.created_at, users.username, analyses.analyt, analyses.material, analyses.assay, analyses.ready_mix, sample_query.firstname, sample_query.lastname, sample_query.sputalysed, sample_query.created_at, sample_query.sample_created_by
			FROM new_sample_analysis
			LEFT JOIN sample_query ON new_sample_analysis.sample_id = sample_query.sample_id
			LEFT JOIN users ON new_sample_analysis.created_by = users.user_id
			LEFT JOIN analyses ON new_sample_analysis.analysis_id = analyses.analysis_id
		`
	err := database.Instance.QueryRow(
		query,
		&sample_analysis.Sample.SampleID,
		&sample_analysis.Analysis.AnalysisID,
		&sample_analysis.Run,
		&sample_analysis.Device,
		user_id).Scan(
		&sample_analysis.CreatedAt,
		&sample_analysis.CreatedBy,
		&sample_analysis.Analysis.Analyt,
		&sample_analysis.Analysis.Material,
		&sample_analysis.Analysis.Assay,
		&sample_analysis.Analysis.ReadyMix,
		&sample_analysis.Sample.FirstName,
		&sample_analysis.Sample.LastName,
		&sample_analysis.Sample.Sputalysed,
		&sample_analysis.Sample.CreatedAt,
		&sample_analysis.Sample.CreatedBy)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, &sample_analysis)
}
