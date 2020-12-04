package ctrls

import (
	"net/http"

	"clevergo.tech/jsend"
	"github.com/gin-gonic/gin"
	"github.com/oraksil/azumma/internal/domain/usecases"
	"github.com/oraksil/azumma/internal/presenter/web"
	"github.com/oraksil/azumma/internal/presenter/web/ctrls/dto"
	"github.com/oraksil/azumma/internal/presenter/web/ctrls/helpers"
)

type PlayerController struct {
	PlayerUseCase *usecases.PlayerUseCase
}

func (ctrl *PlayerController) createNewPlayer(c *gin.Context) {
	sessionCtx := helpers.NewSessionCtx(c)

	type JsonParams struct {
		NickName string `json:"name"`
	}

	var jsonParams JsonParams
	c.BindJSON(&jsonParams)

	player, err := ctrl.PlayerUseCase.CreateNewPlayer(jsonParams.NickName, sessionCtx)
	if err != nil {
		c.JSON(http.StatusOK, jsend.NewFail(map[string]interface{}{
			"code":    400,
			"message": err.Error(),
		}))
		return
	}

	c.JSON(http.StatusOK, jsend.New(dto.PlayerToDto(player)))
}

func (ctrl *PlayerController) getPlayerFromSession(c *gin.Context) {
	sessionCtx := helpers.NewSessionCtx(c)

	player, err := ctrl.PlayerUseCase.GetPlayerFromSession(sessionCtx)
	if err != nil {
		c.JSON(http.StatusOK, jsend.NewFail(map[string]interface{}{
			"code":    400,
			"message": err.Error(),
		}))
		return
	}

	c.JSON(http.StatusOK, jsend.New(dto.PlayerToDto(player)))
}

func (ctrl *PlayerController) useCoin(c *gin.Context) {
	sessionCtx := helpers.NewSessionCtx(c)

	ctrl.PlayerUseCase.UseCoin(1, sessionCtx)

	c.JSON(http.StatusOK, jsend.New(dto.Empty()))
}

func (ctrl *PlayerController) Routes() []web.Route {
	return []web.Route{
		{Spec: "POST /api/v1/players/new", Handler: ctrl.createNewPlayer},
		{Spec: "GET /api/v1/players/me", Handler: ctrl.getPlayerFromSession},
		{Spec: "POST /api/v1/players/coins/use", Handler: ctrl.useCoin},
	}
}
