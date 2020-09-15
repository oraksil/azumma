package usecases

import (
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/oraksil/azumma/internal/domain/models"
	"github.com/oraksil/azumma/internal/domain/services"
	"github.com/oraksil/azumma/pkg/utils"
)

type GameFetchUseCase struct {
	PackRepo models.PackRepository
	GameRepo models.GameRepository
}

func (uc *GameFetchUseCase) GetPacks(page, size int) []*models.Pack {
	return uc.PackRepo.Find(page*size, size)
}

func (uc *GameFetchUseCase) GetGames(page, size int) []*models.Game {
	return uc.GameRepo.Find(page*size, size)
}

type GameCtrlUseCase struct {
	PackRepo       models.PackRepository
	GameRepo       models.GameRepository
	OrakkiDriver   services.OrakkiDriver
	MessageService services.MessageService
	ServiceConfig  *services.ServiceConfig
}

func (uc *GameCtrlUseCase) CreateNewGame(packId int, firstPlayer *models.Player) (*models.Game, error) {
	// validate game
	pack, err := uc.PackRepo.GetById(packId)
	if err != nil {
		return nil, err
	}

	// provision orakki
	newOrakki, err := uc.provisionOrakki()
	if err != nil {
		return nil, err
	}

	// persist orakki context
	game := models.Game{
		Orakki:  newOrakki,
		Pack:    pack,
		Players: []*models.Player{firstPlayer},
	}
	saved, err := uc.GameRepo.Save(&game)
	if err != nil {
		uc.OrakkiDriver.DeleteInstance(newOrakki.Id)
		return nil, err
	}

	// healthcheck if orakki instance is ready
	go uc.postProvisionHandler(&game)

	return saved, nil
}

func (uc *GameCtrlUseCase) provisionOrakki() (*models.Orakki, error) {
	var newOrakkiId, newPeerName string

	if uc.ServiceConfig.UseStaticOrakki {
		newPeerName = uc.ServiceConfig.StaticOrakkiPeerName
		newOrakkiId = uc.ServiceConfig.StaticOrakkiId
	} else {
		newPeerName = utils.NewId("orakki")
		orakkiId, err := uc.OrakkiDriver.RunInstance(newPeerName)
		if err != nil {
			return nil, err
		}
		newOrakkiId = orakkiId
	}

	return &models.Orakki{
		Id:       newOrakkiId,
		PeerName: newPeerName,
		State:    models.ORAKKI_STATE_INIT,
	}, nil
}

func (uc *GameCtrlUseCase) postProvisionHandler(game *models.Game) {
	newOrakki := game.Orakki

	maxWaitTime := uc.ServiceConfig.ProvisionMaxWait * time.Second
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

		if orakkiState.OrakkiId == game.Orakki.Id &&
			orakkiState.State == models.ORAKKI_STATE_READY {
			game.Orakki.State = models.ORAKKI_STATE_READY
			uc.GameRepo.Save(game)
			break
		}

		elapsedTime := time.Since(startTime)
		if elapsedTime > maxWaitTime {
			if !uc.ServiceConfig.UseStaticOrakki {
				uc.OrakkiDriver.DeleteInstance(newOrakki.Id)
			}

			game.Orakki.State = models.ORAKKI_STATE_PANIC
			uc.GameRepo.Save(game)
			break
		}

		time.Sleep(5 * time.Second)
	}
}

func (uc *GameCtrlUseCase) JoinGame() {

}
