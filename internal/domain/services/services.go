package services

import (
	"time"

	"github.com/oraksil/azumma/internal/domain/models"
)

type SessionContext interface {
	GetSession() (*models.Session, error)
	SetSession(session *models.Session) error
	Validate() error
}

type MessageService interface {
	Identifier() string
	Send(to, msgType string, payload interface{}) error
	SendToAny(msgType string, payload interface{}) error
	Broadcast(msgType string, payload interface{}) error
	Request(to, msgType string, payload interface{}, timeout time.Duration) (interface{}, error)
}

type OrakkiDriver interface {
	RunInstance(id string) (string, error)
	DeleteInstance(id string) error
}
