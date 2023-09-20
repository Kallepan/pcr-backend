// go:build wireinject
//go:build wireinject
// +build wireinject

package config

import (
	"github.com/google/wire"
	"gitlab.com/kallepan/pcr-backend/app/controller"
	"gitlab.com/kallepan/pcr-backend/app/repository"
	"gitlab.com/kallepan/pcr-backend/app/service"
)

/* System */
var (
	systemRepoSet = wire.NewSet(repository.SystemRepositoryInit,
		wire.Bind(new(repository.SystemRepository), new(*repository.SystemRepositoryImpl)),
	)
	systemSvcSet = wire.NewSet(service.SystemServiceInit,
		wire.Bind(new(service.SystemService), new(*service.SystemServiceImpl)),
	)
	systemCtrlrSet = wire.NewSet(controller.SystemControllerInit,
		wire.Bind(new(controller.SystemController), new(*controller.SystemControllerImpl)),
	)
)

/* User */
var (
	userRepoSet = wire.NewSet(repository.UserRepositoryInit,
		wire.Bind(new(repository.UserRepository), new(*repository.UserRepositoryImpl)),
	)
	userSvcSet = wire.NewSet(service.UserServiceInit,
		wire.Bind(new(service.UserService), new(*service.UserServiceImpl)),
	)
	userCtrlrSet = wire.NewSet(controller.UserControllerInit,
		wire.Bind(new(controller.UserController), new(*controller.UserControllerImpl)),
	)
)

/* Import */
var (
	importRepoSet = wire.NewSet(repository.ImportRepositoryInit,
		wire.Bind(new(repository.ImportRepository), new(*repository.ImportRepositoryImpl)),
	)
	importSvcSet = wire.NewSet(service.ImportServiceInit,
		wire.Bind(new(service.ImportService), new(*service.ImportServiceImpl)),
	)
	importCtrlrSet = wire.NewSet(controller.ImportControllerInit,
		wire.Bind(new(controller.ImportController), new(*controller.ImportControllerImpl)),
	)
)

/* Print */
var (
	printRepoSet = wire.NewSet(repository.PrintRepositoryInit,
		wire.Bind(new(repository.PrintRepository), new(*repository.PrintRepositoryImpl)),
	)
	printSvcSet = wire.NewSet(service.PrintServiceInit,
		wire.Bind(new(service.PrintService), new(*service.PrintServiceImpl)),
	)
	printCtrlrSet = wire.NewSet(controller.PrintControllerInit,
		wire.Bind(new(controller.PrintController), new(*controller.PrintControllerImpl)),
	)
)

func Init() *Initialization {
	wire.Build(
		NewInitialization,
		systemRepoSet,
		systemSvcSet,
		systemCtrlrSet,
		userRepoSet,
		userSvcSet,
		userCtrlrSet,
		importRepoSet,
		importSvcSet,
		importCtrlrSet,
		printRepoSet,
		printSvcSet,
		printCtrlrSet,
	)
	return nil
}
