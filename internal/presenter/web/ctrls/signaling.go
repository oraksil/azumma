package ctrls

import (
	"net/http"
	"strconv"

	"clevergo.tech/jsend"
	"github.com/gin-gonic/gin"
	"gitlab.com/oraksil/azumma/internal/domain/usecases"
	"gitlab.com/oraksil/azumma/internal/presenter/web"
)

type SignalingController struct {
	SignalingUseCase *usecases.SignalingUseCase
}

func (ctrl *SignalingController) offerHandler(c *gin.Context) {
	gameid, _ := strconv.ParseInt(c.Query("gameid"), 10, 64)
	playerid, _ := strconv.ParseInt(c.PostForm("playerid"), 10, 64)
	offer := c.PostForm("offer")
	connectionInfo, err := ctrl.SignalingUseCase.NewOffer(gameid, playerid, offer)

	var result map[string]string
	if err != nil {
		result = map[string]string{
			"error": err.Error(),
		}
	} else {
		result = map[string]string{
			"state": strconv.Itoa(connectionInfo.State),
		}
	}
	c.JSON(http.StatusOK, jsend.New(result))
}

func (ctrl *SignalingController) getAnswer(c *gin.Context) {
	gameid, _ := strconv.ParseInt(c.Query("gameid"), 10, 64)
	playerid, _ := strconv.ParseInt(c.Query("playerid"), 10, 64)

	hasAnswer, answer, err := ctrl.SignalingUseCase.TryFetchAnswer(gameid, playerid)

	if err != nil {
		result := map[string]string{
			"error": err.Error(),
		}
		c.JSON(http.StatusNotAcceptable, jsend.New(result))
	} else {
		result := map[string]interface{}{
			"answer":    answer,
			"hasAnswer": hasAnswer,
		}

		c.JSON(http.StatusOK, jsend.New(result))
	}
}

func (ctrl *SignalingController) Routes() []web.Route {
	return []web.Route{
		{Spec: "POST /api/v1/signaling/offer", Handler: ctrl.offerHandler},
		{Spec: "GET /api/v1/signaling/answer", Handler: ctrl.getAnswer},
	}
}
