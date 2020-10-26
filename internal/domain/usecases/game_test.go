package usecases

import (
	"testing"
	"time"

	"github.com/oraksil/azumma/internal/domain/models"
	"github.com/oraksil/azumma/internal/domain/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGameFetchUseCaseFindAvailablePacks(t *testing.T) {
	// given
	mockPackRepo := new(MockPackRepository)
	useCase := GameFetchUseCase{
		PackRepo: mockPackRepo,
	}

	mockPacks := []*models.Pack{
		{Id: 1, Title: "Cadilacs", Description: "Game", MaxPlayers: 3},
		{Id: 2, Title: "Bobl Bubl", Description: "Game", MaxPlayers: 2},
	}
	mockPackRepo.On("FindAll", 0, 2).Return(mockPacks)

	// when
	packs := useCase.GetAllPacks(0, 2)

	// then
	assert.Equal(t, len(packs), 2)
	mockPackRepo.AssertExpectations(t)

	// given
	mockPackRepo.On("FindAll", 2, 2).Return(mockPacks)

	// when
	packs = useCase.GetAllPacks(1, 2)

	// then
	assert.Equal(t, len(packs), 2)
	mockPackRepo.AssertExpectations(t)
}

func TestGameCtrlUseCaseCreateNewGame(t *testing.T) {
	// given
	mockPackRepo := new(MockPackRepository)
	mockGameRepo := new(MockGameRepository)
	mockDriver := new(MockK8SOrakkiDriver)
	mockMsgSvc := new(MockMessageService)
	mockSessionCtx := new(MockSessionContext)
	serviceConf := newServiceConfig()

	mockPlayer := models.Player{
		Id:         1,
		Name:       "player0123",
		TotalCoins: 10,
	}

	mockPack := models.Pack{
		Id:          1,
		Title:       "Bubl Boble",
		Maker:       "TAITO",
		Description: "",
		MaxPlayers:  2,
	}

	mockSession := models.Session{
		Player: &mockPlayer,
	}

	mockPackRepo.On("GetById", 1).Return(&mockPack, nil)
	mockGameRepo.On("Save", mock.Anything).Return(mock.Anything, nil)
	mockDriver.On("RunInstance").Return(serviceConf.StaticOrakkiId, nil)
	mockSessionCtx.On("GetSession").Return(&mockSession, nil)

	// when
	useCase := GameCtrlUseCase{
		PackRepo:       mockPackRepo,
		GameRepo:       mockGameRepo,
		OrakkiDriver:   mockDriver,
		MessageService: mockMsgSvc,
		ServiceConfig:  serviceConf,
	}
	game, err := useCase.CreateNewGame(1, mockSessionCtx)

	// then
	assert.NotNil(t, game)
	assert.Nil(t, err)
	assert.Equal(t, serviceConf.StaticOrakkiId, game.Orakki.Id)
	assert.Equal(t, models.OrakkiStateInit, game.Orakki.State)
	assert.Equal(t, 1, len(game.Players))
	assert.Equal(t, &mockPlayer, game.Players[0])

	mockGameRepo.AssertExpectations(t)
	mockDriver.AssertExpectations(t)

	// given
	mockOrakki := models.Orakki{
		Id:    game.Orakki.Id,
		State: models.OrakkiStateReady,
	}
	mockMsgSvc.
		On("Request", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(mockOrakki, nil)

	mockMsgSvc.
		On("Send", mock.Anything, mock.Anything, mock.Anything).
		Return(nil)

	// when
	useCase.postProvisionHandler(game)

	// then
	assert.Equal(t, models.OrakkiStateReady, game.Orakki.State)
}

func newServiceConfig() *services.ServiceConfig {
	return &services.ServiceConfig{
		StaticOrakkiId:   "",
		ProvisionMaxWait: time.Duration(5),
	}
}
