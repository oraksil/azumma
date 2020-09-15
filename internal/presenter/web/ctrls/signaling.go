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
		RunningGameId int64  `json:"running_game_id"`
		SdpOffer      string `json:"sdp_offer"`
	}

	var params JsonParams
	err := c.BindJSON(&params)

	signalingInfo, err := ctrl.SignalingUseCase.NewOffer(params.RunningGameId, player.Id, params.SdpOffer)

	if err != nil {
		c.JSON(http.StatusOK, jsend.NewFail(map[string]interface{}{
			"code":    400,
			"message": "invalid game or player id",
		}))
	} else {
		c.JSON(http.StatusOK, jsend.New(signalingInfo))
	}
}

func (ctrl *SignalingController) getOrakkiIceCandidates(c *gin.Context) {
	// TODO: get player id from session
	// player := ctrl.SessionUseCase.GetPlayerFromSession(...)
	// player := models.Player{Id: 1, Name: "eddy"}

	type QueryParams struct {
		RunningGameId int64 `json:"running_game_id"`
		SinceId       int64 `form:"since_id"`
	}

	var params QueryParams
	err := c.Bind(&params)

	signalingInfo, err := ctrl.SignalingUseCase.GetIceCandidate(params.RunningGameId, params.SinceId)

	if err != nil {
		c.JSON(http.StatusOK, jsend.NewFail(map[string]interface{}{
			"code":    400,
			"message": "no ice candidates availble with given id, seq",
		}))
		return
	}

	c.JSON(http.StatusOK, jsend.New(signalingInfo))
}

func (ctrl *SignalingController) postPlayerIceCandidate(c *gin.Context) {
	// TODO: get player id from session
	// player := ctrl.SessionUseCase.GetPlayerFromSession(...)
	player := models.Player{Id: 1, Name: "eddy"}

	type BodyParams struct {
		RunningGameId int64  `json:"running_game_id"`
		Candidate     string `json:"candidate"`
	}

	var params BodyParams
	err := c.BindJSON(&params)

	result, err := ctrl.SignalingUseCase.AddIceCandidate(params.RunningGameId, player.Id, params.Candidate)

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
