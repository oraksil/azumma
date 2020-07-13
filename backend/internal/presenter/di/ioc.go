package di

import (
	"github.com/golobby/container"
	"oraksil.com/sil/internal/presenter/web/ctrls"
)

func InitContainer() {
	container.Singleton(newGameRepository)
	container.Singleton(newGameFetchUseCase)
	container.Singleton(newGameController)
}

func InjectGameController() *ctrls.GameController {
	var ctrl *ctrls.GameController
	container.Make(&ctrl)
	return ctrl
}
