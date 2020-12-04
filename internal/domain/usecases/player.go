package usecases

import (
	"errors"
	"time"

	"github.com/oraksil/azumma/internal/domain/models"
	"github.com/oraksil/azumma/internal/domain/services"
)

const INITIAL_COINS = 10

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
		LastCoins:       INITIAL_COINS,
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

	session, err := sessionCtx.GetSession()
	if err != nil {
		return nil, err
	}

	player, err := uc.PlayerRepo.GetById(session.Player.Id)
	if err != nil {
		return nil, err
	}

	return player, nil
}
