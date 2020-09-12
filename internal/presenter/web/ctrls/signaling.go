package ctrls

import (
	"net/http"
	"strconv"

	"clevergo.tech/jsend"
	"github.com/gin-gonic/gin"
	"github.com/oraksil/azumma/internal/domain/usecases"
	"github.com/oraksil/azumma/internal/presenter/web"
)

type SignalingController struct {
	SignalingUseCase *usecases.SignalingUseCase
}

func (ctrl *SignalingController) offerHandler(c *gin.Context) {
	orakkiId := c.PostForm("orakki_id")
	playerId, _ := strconv.ParseInt(c.PostForm("player_id"), 10, 64)
	offer := c.PostForm("offer")
	signalingInfo, err := ctrl.SignalingUseCase.NewOffer(orakkiId, playerId, offer)

	if err != nil {
		c.JSON(http.StatusNotAcceptable, jsend.NewFail(map[string]interface{}{
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
	err := c.BindQuery(&params)

	SignalingInfo, err := ctrl.SignalingUseCase.GetIceCandidate(params.OrakkiId, params.SeqAfter, 1)

	if err != nil {
		c.JSON(http.StatusNotAcceptable, jsend.NewFail(map[string]interface{}{
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
	orakkiId := c.PostForm("orakki_id")
	playerId, _ := strconv.ParseInt(c.PostForm("player_id"), 10, 64)
	candidates := c.PostForm("candidates")

	result, err := ctrl.SignalingUseCase.AddIceCandidate(orakkiId, playerId, candidates)

	if err != nil {
		c.JSON(http.StatusNotAcceptable, jsend.NewFail(map[string]interface{}{
			"code":    400,
			"message": "invalid game id",
		}))
	} else {
		response := map[string]interface{}{
			"gameid":   result.Game.Id,
			"playerid": playerId,
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
