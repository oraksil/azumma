package usecases

import (
	"errors"
	"time"

	"github.com/oraksil/azumma/internal/domain/models"
	"github.com/oraksil/azumma/internal/domain/services"
)

type PlayerUseCase struct {
	PlayerRepo models.PlayerRepository
}

func (uc *PlayerUseCase) CreateNewPlayer(
	nickName string, sessionCtx services.SessionContext) (*models.Player, error) {

	session, err := sessionCtx.GetSession()
	if session != nil {
		return nil, errors.New("a player already exists in session")
	}

	newPlayer, err := uc.PlayerRepo.Save(&models.Player{
		Name:            nickName,
		LastCoins:       models.INITIAL_COINS,
		LastCoinsUsedAt: time.Now().UTC(),
	})
	if err != nil {
		return nil, err
	}

	// new session
	session = &models.Session{Player: newPlayer}
	sessionCtx.SetSession(session)

	return newPlayer, nil
}

func (uc *PlayerUseCase) GetPlayerFromSession(
	sessionCtx services.SessionContext) (*models.Player, error) {
	return uc.playerFromSession(sessionCtx)
}

func (uc *PlayerUseCase) UseCoin(numCoins int, sessionCtx services.SessionContext) (*models.Player, error) {
	player, err := uc.playerFromSession(sessionCtx)
	if err != nil {
		return nil, err
	}

	if player.LastCoins < numCoins {
		return nil, errors.New("not enough coins")
	}

	player.UseCoins(numCoins)

	player, err = uc.PlayerRepo.Save(player)

	return player, err
}

func (uc *PlayerUseCase) playerFromSession(sessionCtx services.SessionContext) (*models.Player, error) {
	session, err := sessionCtx.GetSession()
	if err != nil {
		return nil, err
	}

	player, err := uc.PlayerRepo.GetById(session.Player.Id)
	if err != nil {
		return nil, err
	}

	player.UpdateCoins()

	return player, nil
}
