package usecases

import (
	"time"

	"github.com/mitchellh/mapstructure"
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
	ServiceConfig  services.ServiceConfig
}

func (uc *GameCtrlUseCase) CreateNewGame(gameId int, firstPlayer *models.Player) (*models.RunningGame, error) {
	// validate game
	game, err := uc.GameRepository.GetGameById(gameId)
	if err != nil {
		return nil, err
	}

	// provision orakki
	newOrakki, err := uc.provisionOrakki()
	if err != nil {
		return nil, err
	}

	// persist orakki context
	runningGame := models.RunningGame{
		Orakki:  newOrakki,
		Game:    game,
		Players: []*models.Player{firstPlayer},
	}
	saved, err := uc.GameRepository.SaveRunningGame(&runningGame)
	if err != nil {
		uc.OrakkiDriver.DeleteInstance(newOrakki.Id)
		return nil, err
	}

	// healthcheck if orakki instance is ready
	go uc.postProvisionHandler(&runningGame)

	return saved, nil
}

func (uc *GameCtrlUseCase) provisionOrakki() (*models.Orakki, error) {
	var newOrakkiId, newPeerName string

	if uc.ServiceConfig.UseStaticOrakki {
		newOrakkiId = uc.ServiceConfig.StaticOrakkiId
		newPeerName = newOrakkiId
	} else {
		orakkiId, err := uc.OrakkiDriver.RunInstance(newPeerName)
		if err != nil {
			return nil, err
		}
		newOrakkiId = orakkiId
		newPeerName = newOrakkiId
	}

	return &models.Orakki{
		Id:       newOrakkiId,
		PeerName: newPeerName,
		State:    models.ORAKKI_STATE_INIT,
	}, nil
}

func (uc *GameCtrlUseCase) postProvisionHandler(runningGame *models.RunningGame) {
	newOrakki := runningGame.Orakki

	maxWaitTime := 30 * time.Second
	startTime := time.Now()

	for {
		resp, _ := uc.MessageService.Request(
			newOrakki.PeerName,
			models.MSG_FETCH_ORAKKI_STATE,
			"",
			5*time.Second,
		)

		var orakkiState models.OrakkiState
		mapstructure.Decode(resp, &orakkiState)

		if orakkiState.State == models.ORAKKI_STATE_READY {
			runningGame.Orakki.State = models.ORAKKI_STATE_READY
			uc.GameRepository.SaveRunningGame(runningGame)
			break
		}

		elapsedTime := time.Since(startTime)
		if elapsedTime > maxWaitTime {
			uc.OrakkiDriver.DeleteInstance(newOrakki.Id)

			runningGame.Orakki.State = models.ORAKKI_STATE_PANIC
			uc.GameRepository.SaveRunningGame(runningGame)
			break
		}

		time.Sleep(5 * time.Second)
	}
}

func (uc *GameCtrlUseCase) JoinGame() {

}
