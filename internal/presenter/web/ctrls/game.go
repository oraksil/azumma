package ctrls

import (
	"net/http"

	"clevergo.tech/jsend"
	"github.com/gin-gonic/gin"
	"github.com/oraksil/azumma/internal/domain/models"
	"github.com/oraksil/azumma/internal/domain/usecases"
	"github.com/oraksil/azumma/internal/presenter/web"
	"github.com/oraksil/azumma/internal/presenter/web/ctrls/dto"
	"github.com/oraksil/azumma/internal/presenter/web/ctrls/helpers"
)

type GameController struct {
	GameFetchUseCase *usecases.GameFetchUseCase
	GameCtrlUseCase  *usecases.GameCtrlUseCase
}

func (ctrl *GameController) getPacks(c *gin.Context) {
	p := dto.Pagination{Page: 0, Size: 10}
	c.BindQuery(&p)

	type QueryParams struct {
		Status string `form:"status"`
	}

	var queryParams QueryParams
	c.BindQuery(&queryParams)

	var packs []*models.Pack
	switch queryParams.Status {
	case "ready":
		packs = ctrl.GameFetchUseCase.GetPacksByStatus(models.PackStatusReady, p.Page, p.Size)
		break
	case "prepare":
		packs = ctrl.GameFetchUseCase.GetPacksByStatus(models.PackStatusPreparing, p.Page, p.Size)
		break
	default:
		packs = ctrl.GameFetchUseCase.GetAllPacks(p.Page, p.Size)
	}

	c.JSON(http.StatusOK, jsend.New(dto.PacksToDto(packs)))
}

func (ctrl *GameController) createNewGame(c *gin.Context) {
	sessionCtx := helpers.NewSessionCtx(c)
	if sessionCtx.Validate() != nil {
		c.JSON(http.StatusOK, jsend.NewFail(map[string]interface{}{
			"code":    400,
			"message": "invalid session",
		}))
		return
	}

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

	game, err := ctrl.GameCtrlUseCase.CreateNewGame(uriParams.PackId, sessionCtx)
	if err != nil {
		c.JSON(http.StatusOK, jsend.NewFail(map[string]interface{}{
			"code":    400,
			"message": err.Error(),
		}))
		return
	}

	c.JSON(http.StatusOK, jsend.New(dto.GameToDto(game)))
}

func (ctrl *GameController) Routes() []web.Route {
	return []web.Route{
		{Spec: "GET /api/v1/packs", Handler: ctrl.getPacks},
		{Spec: "POST /api/v1/packs/:pack_id/new", Handler: ctrl.createNewGame},
	}
}
