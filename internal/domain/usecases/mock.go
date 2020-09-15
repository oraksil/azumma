package usecases

import (
	"time"

	"github.com/oraksil/azumma/internal/domain/models"
	"github.com/stretchr/testify/mock"
)

type MockGameRepository struct {
	mock.Mock
}

func (r *MockGameRepository) GetById(id int) (*models.Game, error) {
	args := r.Called(id)
	return args.Get(0).(*models.Game), args.Error(1)
}

func (r *MockGameRepository) Find(offset, limit int) []*models.Game {
	args := r.Called(offset, limit)
	return args.Get(0).([]*models.Game)
}

type MockRunningGameRepository struct {
	mock.Mock
}

func (r *MockRunningGameRepository) Find(offset, limit int) []*models.RunningGame {
	args := r.Called(offset, limit)
	return args.Get(0).([]*models.RunningGame)
}

func (r *MockRunningGameRepository) FindById(id int64) (*models.RunningGame, error) {
	args := r.Called(id)
	return args.Get(0).(*models.RunningGame), args.Error(1)
}

func (r *MockRunningGameRepository) Save(game *models.RunningGame) (*models.RunningGame, error) {
	args := r.Called(game)
	game.Id = 1
	return game, args.Error(1)
}

type MockSignalingRepository struct {
	mock.Mock
}

func (m *MockSignalingRepository) Save(signalingInfo *models.SignalingInfo) (*models.SignalingInfo, error) {
	args := m.Called(signalingInfo)
	return args.Get(0).(*models.SignalingInfo), args.Error(1)
}

func (m *MockSignalingRepository) FindByRunningGameId(runningGameId int64, sinceId int64) (*models.SignalingInfo, error) {
	args := m.Called(runningGameId, sinceId)
	return args.Get(0).(*models.SignalingInfo), args.Error(1)
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
