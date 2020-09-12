package usecases

import (
	"testing"
	"time"

	"github.com/oraksil/azumma/internal/domain/models"
	"github.com/oraksil/azumma/internal/domain/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockGameRepository struct {
	mock.Mock
}

func (r *MockGameRepository) GetGameById(id int) (*models.Game, error) {
	args := r.Called(id)
	return args.Get(0).(*models.Game), args.Error(1)
}

func (r *MockGameRepository) FindAvailableGames(offset, limit int) []*models.Game {
	args := r.Called(offset, limit)
	return args.Get(0).([]*models.Game)
}

func (r *MockGameRepository) FindRunningGames(offset, limit int) []*models.RunningGame {
	args := r.Called(offset, limit)
	return args.Get(0).([]*models.RunningGame)
}

func (r *MockGameRepository) SaveRunningGame(game *models.RunningGame) (*models.RunningGame, error) {
	args := r.Called(game)
	game.Id = 1
	return game, args.Error(1)
}

func (r *MockGameRepository) GetSignalingInfo(orakkiId string, playerId int64) (*models.SignalingInfo, error) {
	var signalingInfo *models.SignalingInfo

	return signalingInfo, nil
}

func (r *MockGameRepository) SaveSignalingInfo(signalingInfo *models.SignalingInfo) (*models.SignalingInfo, error) {
	return nil, nil
}

type MockK8SOrakkiDriver struct {
	mock.Mock
}

func (d *MockK8SOrakkiDriver) RunInstance(peerName string) (string, error) {
	args := d.Called(peerName)
	return args.String(0), args.Error(1)
}

func (d *MockK8SOrakkiDriver) DeleteInstance(id string) error {
	args := d.Called(id)
	return args.Error(0)
}

type MockMessageService struct {
	mock.Mock
}

func (m *MockMessageService) Identifier() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockMessageService) Send(to, msgType string, payload interface{}) error {
	args := m.Called(to, msgType, payload)
	return args.Error(0)
}

func (m *MockMessageService) SendToAny(msgType string, payload interface{}) error {
	args := m.Called(msgType, payload)
	return args.Error(0)
}

func (m *MockMessageService) Broadcast(msgType string, payload interface{}) error {
	args := m.Called(msgType, payload)
	return args.Error(0)
}

func (m *MockMessageService) Request(to, msgType string, payload interface{}, timeout time.Duration) (interface{}, error) {
	args := m.Called(to, msgType, payload, timeout)
	return args.Get(0), args.Error(1)
}

func TestGameFetchUseCaseFindAvailableGames(t *testing.T) {
	// given
	mockRepo := new(MockGameRepository)
	useCase := GameFetchUseCase{GameRepository: mockRepo}

	mockGames := []*models.Game{
		{Id: 1, Title: "Cadilacs", Description: "Game", MaxPlayers: 3},
		{Id: 2, Title: "Bobl Bubl", Description: "Game", MaxPlayers: 2},
	}
	mockRepo.On("FindAvailableGames", 0, 2).Return(mockGames)

	// when
	games := useCase.GetAvailableGames(0, 2)

	// then
	assert.Equal(t, len(games), 2)
	mockRepo.AssertExpectations(t)

	// given
	mockRepo.On("FindAvailableGames", 2, 2).Return(mockGames)

	// when
	games = useCase.GetAvailableGames(1, 2)

	// then
	assert.Equal(t, len(games), 2)
	mockRepo.AssertExpectations(t)
}

func TestGameCtrlUseCaseCreateNewGame(t *testing.T) {
	// given
	mockRepo := new(MockGameRepository)
	mockDriver := new(MockK8SOrakkiDriver)
	mockMsgSvc := new(MockMessageService)
	serviceConf := newServiceConfig()

	mockPlayer := models.Player{
		Id:         1,
		Name:       "player0123",
		TotalCoins: 10,
	}

	mockGame := models.Game{
		Id:          1,
		Title:       "Bubl Boble",
		Maker:       "TAITO",
		Description: "",
		MaxPlayers:  2,
	}

	mockRepo.On("GetGameById", 1).Return(&mockGame, nil)
	mockRepo.On("SaveRunningGame", mock.Anything).Return(mock.Anything, nil)
	mockDriver.On("RunInstance", mock.Anything).Return(serviceConf.StaticOrakkiId, nil)

	// when
	useCase := GameCtrlUseCase{
		GameRepository: mockRepo,
		OrakkiDriver:   mockDriver,
		MessageService: mockMsgSvc,
		ServiceConfig:  serviceConf,
	}
	runningGame, err := useCase.CreateNewGame(1, &mockPlayer)

	// then
	assert.NotNil(t, runningGame)
	assert.Nil(t, err)
	assert.Equal(t, serviceConf.StaticOrakkiId, runningGame.Orakki.Id)
	assert.Equal(t, models.ORAKKI_STATE_INIT, runningGame.Orakki.State)
	assert.Equal(t, 1, len(runningGame.Players))
	assert.Equal(t, &mockPlayer, runningGame.Players[0])

	mockRepo.AssertExpectations(t)
	mockDriver.AssertExpectations(t)

	// given
	mockState := models.OrakkiState{
		OrakkiId: runningGame.Orakki.Id,
		State:    models.ORAKKI_STATE_READY,
	}
	mockMsgSvc.
		On("Request", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(mockState, nil)

	// when
	useCase.postProvisionHandler(runningGame)

	// then
	assert.Equal(t, models.ORAKKI_STATE_READY, runningGame.Orakki.State)
}

func newServiceConfig() *services.ServiceConfig {
	return &services.ServiceConfig{
		UseStaticOrakki:  false,
		StaticOrakkiId:   "test-orakki-id",
		ProvisionMaxWait: time.Duration(5),
	}
}
