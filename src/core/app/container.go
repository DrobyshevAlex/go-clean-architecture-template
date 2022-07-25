package app

import (
	"main/src/core/config"
	"main/src/core/db"
	"main/src/core/http"
	v1_0 "main/src/core/http/controllers/v1.0"
	"main/src/core/http/ws"
	userRepo "main/src/features/user/data/repositories"
	userUC "main/src/features/user/domain/usecases"
	userCtrl "main/src/features/user/presentation/controllers"

	"go.uber.org/dig"
)

func BuildContainer() *dig.Container {
	container := dig.New()

	processError(container.Provide(NewApplication))
	processError(container.Provide(ws.NewWSHub))
	processError(container.Provide(http.NewServer))
	processError(container.Provide(config.NewConfig))
	processError(container.Provide(db.NewDb))

	// Controllers
	// v1.0
	processError(container.Provide(v1_0.NewApiController))
	processError(container.Provide(userCtrl.NewUserController))

	// Repositories
	processError(container.Provide(userRepo.NewUsersRepository))

	// UserCases
	processError(container.Provide(userUC.NewUserUsecase))

	return container
}

func processError(err error) {
	if err != nil {
		panic(err)
	}
}
