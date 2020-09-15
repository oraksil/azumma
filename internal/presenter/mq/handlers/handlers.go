package handlers

import (
	"encoding/json"
	// "fmt"

	"github.com/oraksil/azumma/internal/domain/models"
	"github.com/oraksil/azumma/internal/domain/usecases"
	"github.com/sangwonl/mqrpc"
)

type SignalingHandler struct {
	SignalingUseCase *usecases.SignalingUseCase
}

func (h *SignalingHandler) handleIceCandidate(ctx *mqrpc.Context) interface{} {
	var temp models.Icecandidate
	json.Unmarshal(ctx.GetMessage().Payload, &temp)

	h.SignalingUseCase.AddServerIceCandidate(1, temp.IceString)
	return nil
}

func (h *SignalingHandler) Routes() []mqrpc.Route {
	return []mqrpc.Route{
		{MsgType: models.MSG_REMOTE_ICE_CANDIDATE, Handler: h.handleIceCandidate},
	}
}
