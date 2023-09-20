package controller

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/kallepan/pcr-backend/app/service"
)

type SampleController interface {
	AddSample(ctx *gin.Context)
	UpdateSample(ctx *gin.Context)
	DeleteSample(ctx *gin.Context)
	GetSamples(ctx *gin.Context)
}

type SampleControllerImpl struct {
	svc service.SampleService
}

func SampleControllerInit(sampleService service.SampleService) *SampleControllerImpl {
	return &SampleControllerImpl{
		svc: sampleService,
	}
}

func (s SampleControllerImpl) GetSamples(ctx *gin.Context) {
	s.svc.GetSamples(ctx)
}

func (s SampleControllerImpl) AddSample(ctx *gin.Context) {
	s.svc.AddSample(ctx)
}

func (s SampleControllerImpl) UpdateSample(ctx *gin.Context) {
	s.svc.UpdateSample(ctx)
}

func (s SampleControllerImpl) DeleteSample(ctx *gin.Context) {
	s.svc.DeleteSample(ctx)
}
