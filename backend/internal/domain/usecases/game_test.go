package usecases

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"oraksil.com/sil/internal/domain/models"
)

type MockGameRepository struct {
	mock.Mock
}

func (m *MockGameRepository) FindAvailableGames(offset, limit int) []*models.Game {
	m.Called(offset, limit)
	return []*models.Game{
		{Id: 1, Title: "Cadilacs", Description: "Game", MaxPlayers: 3},
		{Id: 2, Title: "Bobl Bubl", Description: "Game", MaxPlayers: 2},
	}
}

func (m *MockGameRepository) FindRunningGames(offset, limit int) []*models.RunningGame {
	return nil
}

func TestGameFetchUseCase(t *testing.T) {
	mockRepo := new(MockGameRepository)

	useCase := GameFetchUseCase{GameRepository: mockRepo}

	mockRepo.On("FindAvailableGames", 0, 2).Return(mock.Anything)
	games := useCase.GetAvailableGames(0, 2)
	assert.Equal(t, len(games), 2)
	mockRepo.AssertExpectations(t)

	mockRepo.On("FindAvailableGames", 2, 2).Return(mock.Anything)
	games = useCase.GetAvailableGames(1, 2)
	assert.Equal(t, len(games), 2)
	mockRepo.AssertExpectations(t)
}
