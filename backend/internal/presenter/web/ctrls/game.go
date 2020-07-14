package ctrls

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"oraksil.com/sil/internal/domain/usecases"
	"oraksil.com/sil/internal/presenter/web"
)

type Pagination struct {
	Page int `form:"page"`
	Size int `form:"size"`
}

type GameController struct {
	GameFetchUseCase *usecases.GameFetchUseCase
}

func (ctrl *GameController) getAvailableGames(c *gin.Context) {
	p := Pagination{Page: 0, Size: 10}
	c.Bind(&p)

	games := ctrl.GameFetchUseCase.GetAvailableGames(p.Page, p.Size)
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   games,
	})
}

func (ctrl *GameController) Routes() []web.Route {
	return []web.Route{
		{Method: web.GET, Url: "/api/v1/games/available", Handler: ctrl.getAvailableGames},
	}
}
