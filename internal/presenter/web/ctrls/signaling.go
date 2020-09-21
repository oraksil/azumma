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

type SignalingController struct {
	SignalingUseCase *usecases.SignalingUseCase
}

func (ctrl *SignalingController) handleSdpExchange(c *gin.Context) {
	// TODO: get player id from session
	// player := ctrl.SessionUseCase.GetPlayerFromSession(...)
	player := models.Player{Id: 1, Name: "gamz"}

	type UriParams struct {
		GameId int64 `uri:"game_id"`
	}

	type JsonParams struct {
		SdpOffer string `json:"sdp_offer"`
	}

	var uriParams UriParams
	c.BindUri(&uriParams)

	var jsonParams JsonParams
	c.BindJSON(&jsonParams)

	sdpInfo, err := ctrl.SignalingUseCase.NewOffer(uriParams.GameId, player.Id, jsonParams.SdpOffer)
	if err != nil {
		c.JSON(http.StatusOK, jsend.NewFail(map[string]interface{}{
			"code":    400,
			"message": "invalid game or player id",
		}))
		return
	}

	c.JSON(http.StatusOK, jsend.New(dto.SdpToDto(sdpInfo)))
}

func (ctrl *SignalingController) getOrakkiIceCandidates(c *gin.Context) {
	// TODO: get player id from session
	// player := ctrl.SessionUseCase.GetPlayerFromSession(...)
	// player := models.Player{Id: 1, Name: "gamz"}

	type UriParams struct {
		GameId int64 `uri:"game_id"`
	}

	type QueryParams struct {
		LastSeq int64 `form:"last_seq"`
	}

	var uriParams UriParams
	c.BindUri(&uriParams)

	var queryParams QueryParams
	c.BindQuery(&queryParams)

	iceCandidate, err := ctrl.SignalingUseCase.GetOrakkiIceCandidate(
		uriParams.GameId,
		queryParams.LastSeq,
	)

	if err != nil {
		c.JSON(http.StatusOK, jsend.New(nil))
		return
	}

	c.JSON(http.StatusOK, jsend.New(dto.IceToDto(iceCandidate)))
}

func (ctrl *SignalingController) postPlayerIceCandidate(c *gin.Context) {
	// TODO: get player id from session
	// player := ctrl.SessionUseCase.GetPlayerFromSession(...)
	player := models.Player{Id: 1, Name: "gamz"}

	type UriParams struct {
		GameId int64 `uri:"game_id"`
	}

	type JsonParams struct {
		IceCandidate string `json:"ice_candidate"`
	}

	var uriParams UriParams
	c.BindUri(&uriParams)

	var jsonParams JsonParams
	c.BindJSON(&jsonParams)

	err := ctrl.SignalingUseCase.OnPlayerIceCandidate(
		uriParams.GameId,
		player.Id,
		jsonParams.IceCandidate,
	)

	if err != nil {
		c.JSON(http.StatusOK, jsend.NewFail(map[string]interface{}{
			"code":    400,
			"message": "failed to post ice candidate",
		}))
		return
	}

	response := map[string]interface{}{}

	c.JSON(http.StatusOK, jsend.New(response))
}

func (ctrl *SignalingController) Routes() []web.Route {
	return []web.Route{
		{Spec: "POST /api/v1/games/:game_id/signaling/sdp", Handler: ctrl.handleSdpExchange},
		{Spec: "GET /api/v1/games/:game_id/signaling/ice", Handler: ctrl.getOrakkiIceCandidates},
		{Spec: "POST /api/v1/games/:game_id/signaling/ice", Handler: ctrl.postPlayerIceCandidate},
	}
}
