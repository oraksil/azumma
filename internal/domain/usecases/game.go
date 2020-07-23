package usecases

import (
	"fmt"

	"gitlab.com/oraksil/azumma/internal/domain/models"
	"gitlab.com/oraksil/azumma/internal/domain/services"
)

type GameFetchUseCase struct {
	GameRepository models.GameRepository
}

func (uc *GameFetchUseCase) GetAvailableGames(page, size int) []*models.Game {
	return uc.GameRepository.FindAvailableGames(page*size, size)
}

func (uc *GameFetchUseCase) GetRunningGames(page, size int) []*models.RunningGame {
	return uc.GameRepository.FindRunningGames(page*size, size)
}

type GameCtrlUseCase struct {
	GameRepository models.GameRepository
	OrakkiDriver   services.OrakkiDriver
	MessageService services.MessageService
}

func (uc *GameCtrlUseCase) CreateNewGame() {
	temp := map[string]string{
		"hello": "world",
	}
	resp := uc.MessageService.Request("orakki-temp", models.MSG_HELLO, temp)
	fmt.Println(resp)
}

func (uc *GameCtrlUseCase) JoinGame() {

}
