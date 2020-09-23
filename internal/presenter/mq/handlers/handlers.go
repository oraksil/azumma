package handlers

import (

	// "fmt"

	"encoding/json"

	"github.com/oraksil/azumma/internal/domain/models"
	"github.com/oraksil/azumma/internal/domain/usecases"
	"github.com/sangwonl/mqrpc"
)

type SignalingHandler struct {
	SignalingUseCase *usecases.SignalingUseCase
}

func (h *SignalingHandler) handleOrakkiIceCandidate(ctx *mqrpc.Context) interface{} {
	var orakkiIce models.IceCandidate
	json.Unmarshal(ctx.GetMessage().Payload, &orakkiIce)

	h.SignalingUseCase.OnOrakkiIceCandidate(orakkiIce.PeerId, orakkiIce.IceBase64Encoded)

	return nil
}

func (h *SignalingHandler) Routes() []mqrpc.Route {
	return []mqrpc.Route{
		{MsgType: models.MsgRemoteIceCandidate, Handler: h.handleOrakkiIceCandidate},
	}
}
