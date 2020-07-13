package di

import (
	"github.com/golobby/container"
	"oraksil.com/sil/internal/domain/models"
	"oraksil.com/sil/internal/domain/usecases"
	"oraksil.com/sil/internal/presenter/data"
	"oraksil.com/sil/internal/presenter/web/ctrls"
)

func newGameRepository() models.GameRepository {
	return &data.GameRepositoryImpl{}
}

func newGameFetchUseCase() *usecases.GameFetchUseCase {
	var repo models.GameRepository
	container.Make(&repo)

	return &usecases.GameFetchUseCase{GameRepository: repo}
}

func newGameController() *ctrls.GameController {
	var gameFetchUseCase *usecases.GameFetchUseCase
	container.Make(&gameFetchUseCase)

	return &ctrls.GameController{GameFetchUseCase: gameFetchUseCase}
}
