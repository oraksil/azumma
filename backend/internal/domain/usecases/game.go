package usecases

import "oraksil.com/sil/internal/domain/models"

type GameFetchUseCase struct {
	GameRepository models.GameRepository
}

func (uc *GameFetchUseCase) GetAvailableGames(page, size int) []*models.Game {
	return uc.GameRepository.GetAllAvailableGames(page*size, size)
}
