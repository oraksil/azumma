package usecases

import (
	"time"

	"github.com/oraksil/azumma/internal/domain/models"
	"github.com/stretchr/testify/mock"
)

type MockPackRepository struct {
	mock.Mock
}

func (r *MockPackRepository) GetById(id int) (*models.Pack, error) {
	args := r.Called(id)
	return args.Get(0).(*models.Pack), args.Error(1)
}

func (r *MockPackRepository) FindAll(offset, limit int) []*models.Pack {
	args := r.Called(offset, limit)
	return args.Get(0).([]*models.Pack)
}

func (r *MockPackRepository) FindByStatus(status, offset, limit int) []*models.Pack {
	args := r.Called(status, offset, limit)
	return args.Get(0).([]*models.Pack)
}

type MockGameRepository struct {
	mock.Mock
}

func (r *MockGameRepository) GetById(id int64) (*models.Game, error) {
	args := r.Called(id)
	return args.Get(0).(*models.Game), args.Error(1)
}

func (r *MockGameRepository) Save(game *models.Game) (*models.Game, error) {
	args := r.Called(game)
	game.Id = 1
	return game, args.Error(1)
}

type MockSignalingRepository struct {
	mock.Mock
}

func (m *MockSignalingRepository) FindByToken(token string, sinceId int64) ([]*models.Signaling, error) {
	args := m.Called(token, sinceId)
	return args.Get(0).([]*models.Signaling), args.Error(1)
}

func (m *MockSignalingRepository) Save(signaling *models.Signaling) (*models.Signaling, error) {
	args := m.Called(signaling)
	return args.Get(0).(*models.Signaling), args.Error(1)
}

type MockSessionContext struct {
	mock.Mock
}

func (m *MockSessionContext) GetSession() (*models.Session, error) {
	args := m.Called()
	return args.Get(0).(*models.Session), args.Error(1)
}

func (m *MockSessionContext) SetSession(session *models.Session) error {
	args := m.Called(session)
	return args.Error(0)
}

func (m *MockSessionContext) Validate() error {
	args := m.Called()
	return args.Error(0)
}

type MockK8SOrakkiDriver struct {
	mock.Mock
}

func (d *MockK8SOrakkiDriver) RunInstance() (string, error) {
	args := d.Called()
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
