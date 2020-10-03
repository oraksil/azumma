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

type SignalingController struct {
	SignalingUseCase *usecases.SignalingUseCase
}

func (ctrl *SignalingController) handleSdpExchange(c *gin.Context) {
	sessionCtx := helpers.NewSessionCtx(c)
	if sessionCtx.Validate() != nil {
		c.JSON(http.StatusOK, jsend.NewFail(map[string]interface{}{
			"code":    400,
			"message": "invalid session",
		}))
		return
	}

	type UriParams struct {
		GameId int64 `uri:"game_id"`
	}

	type JsonParams struct {
		Token    string `json:"token"`
		SdpOffer string `json:"sdp_offer"`
	}

	var uriParams UriParams
	c.BindUri(&uriParams)

	var jsonParams JsonParams
	c.BindJSON(&jsonParams)

	sdpInfo, err := ctrl.SignalingUseCase.NewOffer(
		uriParams.GameId, jsonParams.Token, jsonParams.SdpOffer, sessionCtx)
	if err != nil {
		c.JSON(http.StatusOK, jsend.NewFail(map[string]interface{}{
			"code":    400,
			"message": err.Error(),
		}))
		return
	}

	c.JSON(http.StatusOK, jsend.New(dto.SdpToDto(sdpInfo)))
}

func (ctrl *SignalingController) getOrakkiIceCandidates(c *gin.Context) {
	sessionCtx := helpers.NewSessionCtx(c)
	if sessionCtx.Validate() != nil {
		c.JSON(http.StatusOK, jsend.NewFail(map[string]interface{}{
			"code":    400,
			"message": "invalid session",
		}))
		return
	}

	type UriParams struct {
		GameId int64 `uri:"game_id"`
	}

	type QueryParams struct {
		Token   string `form:"token"`
		LastSeq int64  `form:"last_seq"`
	}

	var uriParams UriParams
	c.BindUri(&uriParams)

	var queryParams QueryParams
	c.BindQuery(&queryParams)

	iceCandidates, err := ctrl.SignalingUseCase.GetOrakkiIceCandidates(
		queryParams.Token,
		queryParams.LastSeq,
		sessionCtx,
	)

	if err != nil {
		c.JSON(http.StatusOK, jsend.New(dto.IcesToDto([]*models.IceCandidate{})))
		return
	}

	c.JSON(http.StatusOK, jsend.New(dto.IcesToDto(iceCandidates)))
}

func (ctrl *SignalingController) postPlayerIceCandidate(c *gin.Context) {
	sessionCtx := helpers.NewSessionCtx(c)
	if sessionCtx.Validate() != nil {
		c.JSON(http.StatusOK, jsend.NewFail(map[string]interface{}{
			"code":    400,
			"message": "invalid session",
		}))
		return
	}

	type UriParams struct {
		GameId int64 `uri:"game_id"`
	}

	type JsonParams struct {
		Token        string `json:"token"`
		IceCandidate string `json:"ice_candidate"`
	}

	var uriParams UriParams
	c.BindUri(&uriParams)

	var jsonParams JsonParams
	c.BindJSON(&jsonParams)

	err := ctrl.SignalingUseCase.OnPlayerIceCandidate(
		uriParams.GameId,
		jsonParams.Token,
		jsonParams.IceCandidate,
		sessionCtx,
	)

	if err != nil {
		c.JSON(http.StatusOK, jsend.NewFail(map[string]interface{}{
			"code":    400,
			"message": "failed to post ice candidate",
		}))
		return
	}

	c.JSON(http.StatusOK, jsend.New(dto.Empty()))
}

func (ctrl *SignalingController) Routes() []web.Route {
	return []web.Route{
		{Spec: "POST /api/v1/games/:game_id/signaling/sdp", Handler: ctrl.handleSdpExchange},
		{Spec: "GET /api/v1/games/:game_id/signaling/ice", Handler: ctrl.getOrakkiIceCandidates},
		{Spec: "POST /api/v1/games/:game_id/signaling/ice", Handler: ctrl.postPlayerIceCandidate},
	}
}
