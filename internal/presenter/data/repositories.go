package data

import (
	"github.com/jmoiron/sqlx"
	"github.com/mitchellh/mapstructure"
	"gitlab.com/oraksil/azumma/internal/domain/models"
	"gitlab.com/oraksil/azumma/internal/presenter/data/dto"
)

type GameRepositoryMySqlImpl struct {
	DB *sqlx.DB
}

func (r *GameRepositoryMySqlImpl) FindAvailableGames(offset, limit int) []*models.Game {
	result := []dto.GameData{}
	r.DB.Select(&result, "select * from game limit ? offset ?", limit, offset)

	var games []*models.Game
	mapstructure.Decode(result, &games)

	for _, g := range games {
		g.Maker = "unknown"
	}

	return games
}

func (r *GameRepositoryMySqlImpl) FindRunningGames(offset, limit int) []*models.RunningGame {
	return nil
}
