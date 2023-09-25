package controller

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/kallepan/pcr-backend/app/service"
)

type SamplePanelController interface {
	GetSamplePanels(ctx *gin.Context)
	ResetSamplePanel(ctx *gin.Context)
	UpdateSamplePanel(ctx *gin.Context)
	CreateSamplePanel(ctx *gin.Context)

	GetStatistics(ctx *gin.Context)

	CreateRun(ctx *gin.Context)
}

type SamplePanelControllerImpl struct {
	svc    service.SamplePanelService
	runSVc service.RunService
}

func SamplePanelControllerInit(samplePanelService service.SamplePanelService, runService service.RunService) *SamplePanelControllerImpl {
	return &SamplePanelControllerImpl{
		svc:    samplePanelService,
		runSVc: runService,
	}
}

func (s SamplePanelControllerImpl) GetSamplePanels(ctx *gin.Context) {
	s.svc.GetSamplePanels(ctx)
}

func (s SamplePanelControllerImpl) ResetSamplePanel(ctx *gin.Context) {
	s.svc.ResetSamplePanel(ctx)
}

func (s SamplePanelControllerImpl) UpdateSamplePanel(ctx *gin.Context) {
	s.svc.UpdateSamplePanel(ctx)
}

func (s SamplePanelControllerImpl) CreateSamplePanel(ctx *gin.Context) {
	s.svc.CreateSamplePanel(ctx)
}

func (s SamplePanelControllerImpl) GetStatistics(ctx *gin.Context) {
	s.svc.GetStatistics(ctx)
}

func (s SamplePanelControllerImpl) CreateRun(ctx *gin.Context) {
	s.runSVc.CreateRun(ctx)
}
