package usecases

import (
	"errors"
	"fmt"
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

func (uc *GameFetchUseCase) GetAllPacks(page, size int) []*models.Pack {
	return uc.PackRepo.FindAll(page*size, size)
}

func (uc *GameFetchUseCase) GetPacksByStatus(status, page, size int) []*models.Pack {
	return uc.PackRepo.FindByStatus(status, page*size, size)
}

type GameCtrlUseCase struct {
	ServiceConfig *services.ServiceConfig

	PackRepo   models.PackRepository
	GameRepo   models.GameRepository
	PlayerRepo models.PlayerRepository

	OrakkiDriver   services.OrakkiDriver
	MessageService services.MessageService
}

func (uc *GameCtrlUseCase) CreateNewGame(packId int, sessionCtx services.SessionContext) (*models.Game, error) {
	// validate game
	pack, err := uc.PackRepo.GetById(packId)
	if err != nil {
		return nil, err
	}

	// try to create new game
	session, _ := sessionCtx.GetSession()
	game := models.Game{
		Pack:    pack,
		Players: []*models.Player{session.Player},
	}
	saved, err := uc.GameRepo.Save(&game)
	if err != nil {
		return nil, err
	}

	// provision orakki
	provisionInfo := &models.ProvisionInfo{
		OrakkiId: fmt.Sprintf("orakki-%d", saved.Id),
		Pack:     pack,
	}
	newOrakki, err := uc.provisionOrakki(provisionInfo)
	if err != nil {
		return nil, err
	}

	// persist orakki context
	saved.Orakki = newOrakki
	saved, err = uc.GameRepo.Save(&game)
	if err != nil {
		uc.OrakkiDriver.DeleteInstance(newOrakki.Id)
		return nil, err
	}

	// healthcheck if orakki instance is ready
	go uc.postProvisionHandler(&game)

	return saved, nil
}

func (uc *GameCtrlUseCase) provisionOrakki(info *models.ProvisionInfo) (*models.Orakki, error) {
	var err error
	var orakkiId string
	if uc.ServiceConfig.StaticOrakkiId != "" {
		orakkiId = uc.ServiceConfig.StaticOrakkiId
	} else {
		orakkiId, err = uc.OrakkiDriver.RunInstance(info.OrakkiId, info.Pack.RomName)
		if err != nil {
			return nil, err
		}
	}

	return &models.Orakki{
		Id:    orakkiId,
		State: models.OrakkiStateInit,
	}, nil
}

func (uc *GameCtrlUseCase) postProvisionHandler(game *models.Game) {
	newOrakki := game.Orakki

	maxWaitTime := uc.ServiceConfig.ProvisionMaxWait * time.Second
	startTime := time.Now()

	for {
		resp, _ := uc.MessageService.Request(
			newOrakki.Id,
			models.MsgPrepareOrakki,
			&models.PrepareOrakki{GameId: game.Id},
			5*time.Second,
		)

		var orakki models.Orakki
		mapstructure.Decode(resp, &orakki)

		if orakki.Id == game.Orakki.Id &&
			orakki.State == models.OrakkiStateReady {
			game.Orakki.State = models.OrakkiStateReady
			uc.GameRepo.Save(game)
			break
		}

		elapsedTime := time.Since(startTime)
		if elapsedTime > maxWaitTime {
			if uc.ServiceConfig.StaticOrakkiId == "" {
				uc.OrakkiDriver.DeleteInstance(newOrakki.Id)
			}

			game.Orakki.State = models.OrakkiStatePanic
			uc.GameRepo.Save(game)
			break
		}

		time.Sleep(2 * time.Second)
	}

	if game.Orakki.State == models.OrakkiStateReady {
		uc.MessageService.Send(
			newOrakki.Id,
			models.MsgStartGame,
			&models.GameInfo{GameId: game.Id, MaxPlayers: game.Pack.MaxPlayers},
		)
	}
}

func (uc *GameCtrlUseCase) CanJoinGame(gameId int64, sessionCtx services.SessionContext) (string, error) {
	game, err := uc.GameRepo.GetById(gameId)
	if err != nil {
		return "", errors.New("game not found")
	}

	session, err := sessionCtx.GetSession()
	if err != nil || session.Player == nil {
		return "", errors.New("invalid player")
	}

	if len(game.Players) >= game.Pack.MaxPlayers {
		return "", errors.New("game is full")
	}

	joinToken := utils.NewId("")

	return joinToken, nil
}

func (uc *GameCtrlUseCase) JoinGame(gameId int64, playerId int64) (*models.Game, error) {
	game, err := uc.GameRepo.GetById(gameId)
	if err != nil {
		return nil, err
	}

	player, err := uc.PlayerRepo.GetById(playerId)
	if err != nil {
		return nil, err
	}

	game.Join(player)

	_, err = uc.GameRepo.Save(game)
	if err != nil {
		return nil, err
	}

	return game, nil
}

func (uc *GameCtrlUseCase) LeaveGame(gameId int64, playerId int64) error {
	game, err := uc.GameRepo.GetById(gameId)
	if err != nil {
		return err
	}

	player, err := uc.PlayerRepo.GetById(playerId)
	if err != nil {
		return err
	}

	game.Leave(player)

	_, err = uc.GameRepo.Save(game)
	if err != nil {
		return err
	}

	return nil
}
