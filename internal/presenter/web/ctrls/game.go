package ctrls

import (
	"net/http"

	"clevergo.tech/jsend"
	"github.com/gin-gonic/gin"
	"github.com/oraksil/azumma/internal/domain/models"
	"github.com/oraksil/azumma/internal/domain/usecases"
	"github.com/oraksil/azumma/internal/presenter/web"
	"github.com/oraksil/azumma/internal/presenter/web/ctrls/dto"
)

type GameController struct {
	GameFetchUseCase *usecases.GameFetchUseCase
	GameCtrlUseCase  *usecases.GameCtrlUseCase
}

func (ctrl *GameController) getAvailablePacks(c *gin.Context) {
	p := dto.Pagination{Page: 0, Size: 10}
	c.BindQuery(&p)

	packs := ctrl.GameFetchUseCase.GetPacks(p.Page, p.Size)

	c.JSON(http.StatusOK, jsend.New(dto.PackToDto(packs)))
}

func (ctrl *GameController) createNewGame(c *gin.Context) {
	// TODO: get player id from session
	// player := ctrl.SessionUseCase.GetPlayerFromSession(...)
	player := models.Player{Id: 1, Name: "gamz"}

	type UriParams struct {
		PackId int `uri:"pack_id"`
	}

	var uriParams UriParams
	err := c.BindUri(&uriParams)
	if err != nil {
		c.JSON(http.StatusOK, jsend.NewFail(map[string]interface{}{
			"code":    400,
			"message": "invalid pack id",
		}))
		return
	}

	game, err := ctrl.GameCtrlUseCase.CreateNewGame(uriParams.PackId, &player)
	if err != nil {
		c.JSON(http.StatusOK, jsend.NewFail(map[string]interface{}{
			"code":    400,
			"message": "failed to create a new game",
		}))
		return
	}

	c.JSON(http.StatusOK, jsend.New(dto.GameToDto(game)))
}

func (ctrl *GameController) joinGame(c *gin.Context) {

}

func (ctrl *GameController) Routes() []web.Route {
	return []web.Route{
		{Spec: "GET /api/v1/packs", Handler: ctrl.getAvailablePacks},
		{Spec: "POST /api/v1/packs/:pack_id/new", Handler: ctrl.createNewGame},
		{Spec: "POST /api/v1/games/:game_id/join", Handler: ctrl.joinGame},
	}
}
