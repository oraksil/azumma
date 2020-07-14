package ctrls

import (
	"net/http"

	"clevergo.tech/jsend"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"oraksil.com/sil/internal/domain/usecases"
	"oraksil.com/sil/internal/presenter/web"
	"oraksil.com/sil/internal/presenter/web/ctrls/dto"
)

type GameController struct {
	GameFetchUseCase *usecases.GameFetchUseCase
}

func (ctrl *GameController) getAvailableGames(c *gin.Context) {
	p := dto.Pagination{Page: 0, Size: 10}
	c.Bind(&p)

	games := ctrl.GameFetchUseCase.GetAvailableGames(p.Page, p.Size)

	var gamesDto []dto.AvailableGame
	mapstructure.Decode(games, &gamesDto)

	c.JSON(http.StatusOK, jsend.New(gamesDto))
}

func (ctrl *GameController) Routes() []web.Route {
	return []web.Route{
		{Spec: "GET /api/v1/games/available", Handler: ctrl.getAvailableGames},
	}
}
