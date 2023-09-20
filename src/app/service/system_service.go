package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gitlab.com/kallepan/pcr-backend/app/constant"
	"gitlab.com/kallepan/pcr-backend/app/pkg"
	"gitlab.com/kallepan/pcr-backend/app/repository"
)

type SystemService interface {
	GetPing(c *gin.Context)
}

type SystemServiceImpl struct {
	systemRepository repository.SystemRepository
}

func (s SystemServiceImpl) GetPing(c *gin.Context) {
	defer pkg.PanicHandler(c)

	data, err := s.systemRepository.GetVersion()
	if err != nil {
		log.Error("Got error when get version: ", err)
		pkg.PanicException(constant.DataNotFound)
	}

	c.JSON(http.StatusOK, pkg.BuildResponse(constant.Success, data))
}

func SystemServiceInit(systemRepository repository.SystemRepository) *SystemServiceImpl {
	return &SystemServiceImpl{
		systemRepository: systemRepository,
	}
}
