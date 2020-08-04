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

func (ctrl *SignalingController) sdpHandler(c *gin.Context) {
	gameid := c.Param("gameid")
	playerid, _ := strconv.ParseInt(c.PostForm("playerid"), 10, 64)
	offer := c.PostForm("offer")
	connectionInfo, _ := ctrl.SignalingUseCase.NewOffer(gameid, playerid, offer)

	empty := map[string]string{
		"status": strconv.Itoa(connectionInfo.State),
	}
	c.JSON(http.StatusOK, jsend.New(empty))
}

func (ctrl *SignalingController) Routes() []web.Route {
	return []web.Route{
		{Spec: "POST /api/v1/games/:gameid/signaling/sdp", Handler: ctrl.sdpHandler},
	}
}
