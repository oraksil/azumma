package ctrls

import (
	"net/http"

	"clevergo.tech/jsend"
	"github.com/gin-gonic/gin"
	"github.com/oraksil/azumma/internal/domain/models"
	"github.com/oraksil/azumma/internal/domain/usecases"
	"github.com/oraksil/azumma/internal/presenter/web"
)

type SignalingController struct {
	SignalingUseCase *usecases.SignalingUseCase
}

func (ctrl *SignalingController) handleSdpExchange(c *gin.Context) {
	// TODO: get player id from session
	// player := ctrl.SessionUseCase.GetPlayerFromSession(...)
	player := models.Player{Id: 1, Name: "eddy"}

	type JsonParams struct {
		gameId   int64  `json:"game_id"`
		sdpOffer string `json:"sdp_offer"`
	}

	var params JsonParams
	err := c.BindJSON(&params)

	signaling, err := ctrl.SignalingUseCase.NewOffer(params.gameId, player.Id, params.sdpOffer)

	if err != nil {
		c.JSON(http.StatusOK, jsend.NewFail(map[string]interface{}{
			"code":    400,
			"message": "invalid game or player id",
		}))
	} else {
		c.JSON(http.StatusOK, jsend.New(signaling))
	}
}

func (ctrl *SignalingController) getOrakkiIceCandidates(c *gin.Context) {
	// TODO: get player id from session
	// player := ctrl.SessionUseCase.GetPlayerFromSession(...)
	// player := models.Player{Id: 1, Name: "eddy"}

	type QueryParams struct {
		gameId  int64 `json:"game_id"`
		sinceId int64 `form:"since_id"`
	}

	var params QueryParams
	err := c.Bind(&params)

	signaling, err := ctrl.SignalingUseCase.GetIceCandidate(params.gameId, params.sinceId)

	if err != nil {
		c.JSON(http.StatusOK, jsend.NewFail(map[string]interface{}{
			"code":    400,
			"message": "no ice candidates availble with given id, seq",
		}))
		return
	}

	c.JSON(http.StatusOK, jsend.New(signaling))
}

func (ctrl *SignalingController) postPlayerIceCandidate(c *gin.Context) {
	// TODO: get player id from session
	// player := ctrl.SessionUseCase.GetPlayerFromSession(...)
	player := models.Player{Id: 1, Name: "eddy"}

	type JsonParams struct {
		gameId    int64  `json:"game_id"`
		candidate string `json:"candidate"`
	}

	var params JsonParams
	err := c.BindJSON(&params)

	result, err := ctrl.SignalingUseCase.AddIceCandidate(params.gameId, player.Id, params.candidate)

	if err != nil {
		c.JSON(http.StatusOK, jsend.NewFail(map[string]interface{}{
			"code":    400,
			"message": "invalid game id",
		}))
	} else {
		response := map[string]interface{}{
			"gameid":   result.Game.Id,
			"playerid": player.Id,
		}

		c.JSON(http.StatusOK, jsend.New(response))
	}
}

func (ctrl *SignalingController) Routes() []web.Route {
	return []web.Route{
		{Spec: "POST /api/v1/signaling/:running_game_id/offer", Handler: ctrl.handleSdpExchange},
		{Spec: "GET /api/v1/signaling/:running_game_id/ice", Handler: ctrl.getOrakkiIceCandidates},
		{Spec: "POST /api/v1/signaling/:running_game_id/ice", Handler: ctrl.postPlayerIceCandidate},
	}
}
