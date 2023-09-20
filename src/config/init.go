package config

import (
	"gitlab.com/kallepan/pcr-backend/app/controller"
	"gitlab.com/kallepan/pcr-backend/app/repository"
	"gitlab.com/kallepan/pcr-backend/app/service"
)

type Initialization struct {
	sysRepo    repository.SystemRepository
	sysSvc     service.SystemService
	SysCtrl    controller.SystemController
	userRepo   repository.UserRepository
	userSvc    service.UserService
	UserCtrl   controller.UserController
	importRepo repository.ImportRepository
	importSvc  service.ImportService
	ImportCtrl controller.ImportController
}

func NewInitialization(
	sysRepo repository.SystemRepository,
	sysSvc service.SystemService,
	sysCtrl controller.SystemController,
	userRepo repository.UserRepository,
	userSvc service.UserService,
	userCtrl controller.UserController,
	importRepo repository.ImportRepository,
	importSvc service.ImportService,
	importCtrl controller.ImportController,
) *Initialization {
	return &Initialization{
		sysRepo:    sysRepo,
		sysSvc:     sysSvc,
		SysCtrl:    sysCtrl,
		userRepo:   userRepo,
		userSvc:    userSvc,
		UserCtrl:   userCtrl,
		importRepo: importRepo,
		importSvc:  importSvc,
		ImportCtrl: importCtrl,
	}
}
