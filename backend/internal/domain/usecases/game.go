package usecases

import (
	"fmt"

	"gitlab.com/oraksil/sil/backend/internal/domain/models"
	"gitlab.com/oraksil/sil/backend/internal/domain/services"
	"gitlab.com/oraksil/sil/backend/pkg/mq"
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
	MessageService mq.MessageService
}

func (uc *GameCtrlUseCase) CreateNewGame() {
	temp1 := map[string]string{
		"abcdcdcd": "broadcast",
	}

	temp2 := map[string]string{
		"abcdcdcd": "p2p",
	}
	uc.MessageService.Broadcast(models.MSG_HELLO, temp1)
	uc.MessageService.Send("generated", models.MSG_HELLO, temp2)
	resp := uc.MessageService.Request("generated", models.MSG_HELLO, temp2)
	fmt.Println(resp)
}

func (uc *GameCtrlUseCase) JoinGame() {

}
