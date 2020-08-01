package usecases

import (
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

	// healthcheck if orakki instance is ready
	go func() {
		// for {
		// resp, _ := uc.MessageService.Request(newPeerName, "msg-fetch-state", "")
		// }
	}()

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

	return saved, nil
}

func (uc *GameCtrlUseCase) provisionOrakki() (*models.Orakki, error) {
	// TODO: generate unique and not too long orakki peer name
	newPeerName := generateNewPeerName("orakki")

	// provisioning new orakki
	newOrakkiId, err := uc.OrakkiDriver.RunInstance(newPeerName)
	if err != nil {
		return nil, err
	}

	return &models.Orakki{
		Id:       newOrakkiId,
		State:    models.ORAKKI_STATE_INIT,
		PeerName: newPeerName,
	}, nil
}

func generateNewPeerName(prefix string) string {
	return "orakki-cfea2hea"
}

func (uc *GameCtrlUseCase) JoinGame() {

}
