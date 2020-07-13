package data

import "oraksil.com/sil/internal/domain/models"

type GameRepositoryImpl struct {
}

func (r *GameRepositoryImpl) GetAllAvailableGames(offset, limit int) []*models.Game {
	return nil
}

func (r *GameRepositoryImpl) GetAllRunningGames(offset, limit int) []*models.RunningGame {
	return nil
}
