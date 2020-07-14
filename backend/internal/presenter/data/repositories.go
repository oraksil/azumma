package data

import (
	"github.com/jmoiron/sqlx"
	"oraksil.com/sil/internal/domain/models"
)

type GameRepositoryMySqlImpl struct {
	DB *sqlx.DB
}

func (r *GameRepositoryMySqlImpl) FindAvailableGames(offset, limit int) []*models.Game {
	games := []*models.Game{}
	r.DB.Select(&games, "select * from game limit ? offset ?", limit, offset)
	return games
}

func (r *GameRepositoryMySqlImpl) FindRunningGames(offset, limit int) []*models.RunningGame {
	return nil
}
