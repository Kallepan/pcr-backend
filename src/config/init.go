package config

import (
	"gitlab.com/kallepan/pcr-backend/app/controller"
	"gitlab.com/kallepan/pcr-backend/app/repository"
	"gitlab.com/kallepan/pcr-backend/app/service"
)

type Initialization struct {
	sysRepo         repository.SystemRepository
	sysSvc          service.SystemService
	SysCtrl         controller.SystemController
	userRepo        repository.UserRepository
	userSvc         service.UserService
	UserCtrl        controller.UserController
	importRepo      repository.ImportRepository
	importSvc       service.ImportService
	ImportCtrl      controller.ImportController
	printRepo       repository.PrintRepository
	printSvc        service.PrintService
	PrintCtrl       controller.PrintController
	panelRepo       repository.PanelRepository
	panelSvc        service.PanelService
	PanelCtrl       controller.PanelController
	sampleRepo      repository.SampleRepository
	sampleSvc       service.SampleService
	SampleCtrl      controller.SampleController
	SynchroRepo     repository.SynchronizeRepository
	samplePanelRepo repository.SamplePanelRepository
	samplePanelSvc  service.SamplePanelService
	SamplePanelCtrl controller.SamplePanelController
	runRepo         repository.RunRepository
	runSvc          service.RunService
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
	printRepo repository.PrintRepository,
	printSvc service.PrintService,
	printCtrl controller.PrintController,
	panelRepo repository.PanelRepository,
	panelSvc service.PanelService,
	panelCtrl controller.PanelController,
	sampleRepo repository.SampleRepository,
	sampleSvc service.SampleService,
	sampleCtrl controller.SampleController,
	synchroRepo repository.SynchronizeRepository,
	samplePanelRepo repository.SamplePanelRepository,
	samplePanelSvc service.SamplePanelService,
	samplePanelCtrl controller.SamplePanelController,
	runRepo repository.RunRepository,
	runSvc service.RunService,
) *Initialization {
	return &Initialization{
		sysRepo:         sysRepo,
		sysSvc:          sysSvc,
		SysCtrl:         sysCtrl,
		userRepo:        userRepo,
		userSvc:         userSvc,
		UserCtrl:        userCtrl,
		importRepo:      importRepo,
		importSvc:       importSvc,
		ImportCtrl:      importCtrl,
		printRepo:       printRepo,
		printSvc:        printSvc,
		PrintCtrl:       printCtrl,
		panelRepo:       panelRepo,
		panelSvc:        panelSvc,
		PanelCtrl:       panelCtrl,
		sampleRepo:      sampleRepo,
		sampleSvc:       sampleSvc,
		SampleCtrl:      sampleCtrl,
		SynchroRepo:     synchroRepo,
		samplePanelRepo: samplePanelRepo,
		samplePanelSvc:  samplePanelSvc,
		SamplePanelCtrl: samplePanelCtrl,
		runRepo:         runRepo,
		runSvc:          runSvc,
	}
}
