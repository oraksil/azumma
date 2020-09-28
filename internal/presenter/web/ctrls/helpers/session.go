package helpers

import (
	"errors"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/oraksil/azumma/internal/domain/models"
	"github.com/oraksil/azumma/internal/domain/services"
)

type GinSessionContext struct {
	ginContext *gin.Context
}

func extractPlayer(sessionStore sessions.Session) (*models.Player, error) {
	playerId := sessionStore.Get("player_id")
	playerName := sessionStore.Get("player_name")

	if playerId == nil || playerName == nil {
		return nil, errors.New("no player info found in session")
	}

	return &models.Player{
		Id:   playerId.(int64),
		Name: playerName.(string),
	}, nil
}

func storePlayer(sessionStore sessions.Session, player *models.Player) error {
	sessionStore.Set("player_id", player.Id)
	sessionStore.Set("player_name", player.Name)

	return sessionStore.Save()
}

func (c *GinSessionContext) GetSession() (*models.Session, error) {
	sessionStore := sessions.Default(c.ginContext)

	player, err := extractPlayer(sessionStore)
	if err != nil {
		return nil, err
	}

	return &models.Session{Player: player}, nil
}

func (c *GinSessionContext) SetSession(session *models.Session) error {
	sessionStore := sessions.Default(c.ginContext)

	return storePlayer(sessionStore, session.Player)
}

func (c *GinSessionContext) Validate() error {
	sessionStore := sessions.Default(c.ginContext)

	_, err := extractPlayer(sessionStore)

	return err
}

func NewSessionCtx(c *gin.Context) services.SessionContext {
	newCtx := GinSessionContext{ginContext: c}
	return &newCtx
}
