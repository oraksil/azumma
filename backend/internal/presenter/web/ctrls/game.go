package ctrls

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"oraksil.com/sil/internal/domain/models"
	"oraksil.com/sil/internal/domain/usecases"
	"oraksil.com/sil/internal/presenter/web"
)

type GameFetchUseCase interface {
	GetAvailableGames(page, size int) []*models.Game
	GetRunningGames(page, size int) []*models.RunningGame
}

type GameController struct {
	GameFetchUseCase *usecases.GameFetchUseCase
}

func (ctrl *GameController) getAvailableGames(c *gin.Context) {
	ctrl.GameFetchUseCase.GetAvailableGames(0, 10)
	c.JSON(http.StatusOK, gin.H{"say": "world"})
}

func (ctrl *GameController) Routes() []web.Route {
	return []web.Route{
		{Method: web.GET, Url: "/api/v1/games/available", Handler: ctrl.getAvailableGames},
	}
}
