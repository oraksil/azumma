package data

import "oraksil.com/sil/internal/domain/models"

type GameRepository struct {
}

func (r *GameRepository) GetAllAvailableGames(offset, limit int) []*models.Game {
	return nil
}

func (r *GameRepository) GetAllRunningGames(offset, limit int) []*models.RunningGame {
	return nil
}
