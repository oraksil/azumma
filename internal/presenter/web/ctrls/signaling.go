package ctrls

import (
	"net/http"

	"clevergo.tech/jsend"
	"github.com/gin-gonic/gin"
	"github.com/oraksil/azumma/internal/domain/usecases"
	"github.com/oraksil/azumma/internal/presenter/web"
)

type SignalingController struct {
	SignalingUseCase *usecases.SignalingUseCase
}

func (ctrl *SignalingController) offerHandler(c *gin.Context) {

	type BodyParams struct {
		OrakkiId string `json:"orakki_id" form:"orakki_id"`
		PlayerId int64  `json:"player_id" form:"player_id"`
		Offer    string `json:"offer" form:"offer"`
	}

	var params BodyParams
	err := c.Bind(&params)

	signalingInfo, err := ctrl.SignalingUseCase.NewOffer(params.OrakkiId, params.PlayerId, params.Offer)

	if err != nil {
		c.JSON(http.StatusOK, jsend.NewFail(map[string]interface{}{
			"code":    400,
			"message": "invalid game or player id",
		}))
	} else {
		response := map[string]string{
			"data": signalingInfo.Data,
		}
		c.JSON(http.StatusOK, jsend.New(response))
	}
}

func (ctrl *SignalingController) getIceCandidate(c *gin.Context) {
	type QueryParams struct {
		OrakkiId string `form:"orakki_id"`
		PlayerId int64  `form:"player_id"`
		SeqAfter int    `form:"seq_after"`
	}

	var params QueryParams
	err := c.Bind(&params)

	SignalingInfo, err := ctrl.SignalingUseCase.GetIceCandidate(params.OrakkiId, params.SeqAfter)

	if err != nil {
		c.JSON(http.StatusOK, jsend.NewFail(map[string]interface{}{
			"code":    400,
			"message": "no ice candidates availble with given id, seq",
		}))
	} else {
		result := map[string]interface{}{
			"seq":          SignalingInfo.Id,
			"icecandidate": SignalingInfo.Data,
			"isLast":       SignalingInfo.IsLast,
		}

		c.JSON(http.StatusOK, jsend.New(result))
	}
}

func (ctrl *SignalingController) postIceCandidate(c *gin.Context) {
	type BodyParams struct {
		OrakkiId  string `json:"orakki_id" form:"orakki_id"`
		PlayerId  int64  `json:"player_id" form:"player_id"`
		Candidate string `json:"candidate" form:"candidate"`
	}

	var params BodyParams
	err := c.Bind(&params)

	result, err := ctrl.SignalingUseCase.AddIceCandidate(params.OrakkiId, params.PlayerId, params.Candidate)

	if err != nil {
		c.JSON(http.StatusOK, jsend.NewFail(map[string]interface{}{
			"code":    400,
			"message": "invalid game id",
		}))
	} else {
		response := map[string]interface{}{
			"gameid":   result.Game.Id,
			"playerid": params.PlayerId,
		}

		c.JSON(http.StatusOK, jsend.New(response))
	}
}

func (ctrl *SignalingController) Routes() []web.Route {
	return []web.Route{
		{Spec: "POST /api/v1/signaling/offer", Handler: ctrl.offerHandler},
		{Spec: "GET /api/v1/signaling/ice", Handler: ctrl.getIceCandidate},
		{Spec: "POST /api/v1/signaling/ice", Handler: ctrl.postIceCandidate},
	}
}
