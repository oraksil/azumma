package ctrls

import (
	"net/http"
	"reflect"
	"time"

	"clevergo.tech/jsend"
	"github.com/gin-gonic/gin"
	"gitlab.com/oraksil/azumma/internal/domain/models"
	"gitlab.com/oraksil/azumma/internal/domain/usecases"
	"gitlab.com/oraksil/azumma/internal/presenter/web"
	"gitlab.com/oraksil/azumma/internal/presenter/web/ctrls/dto"
)

type GameController struct {
	GameFetchUseCase *usecases.GameFetchUseCase
	GameCtrlUseCase  *usecases.GameCtrlUseCase
}

func (ctrl *GameController) getAvailableGames(c *gin.Context) {
	p := dto.Pagination{Page: 0, Size: 10}
	c.Bind(&p)

	games := ctrl.GameFetchUseCase.GetAvailableGames(p.Page, p.Size)

	c.JSON(http.StatusOK, jsend.New(dto.GamesToDto(games)))
}

func timeToIntDecodeHook(
	f reflect.Kind,
	t reflect.Kind,
	data interface{}) (interface{}, error) {
	return data.(time.Time).Second(), nil
}

func (ctrl *GameController) createNewGame(c *gin.Context) {
	// TODO: get player id from session
	// player := ctrl.SessionUseCase.GetPlayerFromSession(...)
	player := models.Player{Id: 1, Name: "eddy"}

	type QueryParams struct {
		GameId int `form:"game_id"`
	}

	var params QueryParams
	err := c.BindQuery(&params)

	if err != nil {
		c.JSON(http.StatusOK, jsend.NewFail(map[string]interface{}{
			"code":    400,
			"message": "invalid game id",
		}))
		return
	}

	runningGame, err := ctrl.GameCtrlUseCase.CreateNewGame(params.GameId, &player)
	if err != nil {
		c.JSON(http.StatusOK, jsend.NewFail(map[string]interface{}{
			"code":    400,
			"message": "failed to create a new game",
		}))
		return
	}

	c.JSON(http.StatusOK, jsend.New(dto.RunningGameToDto(runningGame)))
}

func (ctrl *GameController) Routes() []web.Route {
	return []web.Route{
		{Spec: "GET /api/v1/games/available", Handler: ctrl.getAvailableGames},
		{Spec: "POST /api/v1/games/new", Handler: ctrl.createNewGame},
	}
}
