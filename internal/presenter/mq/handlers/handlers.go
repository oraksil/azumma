package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/oraksil/azumma/internal/domain/models"
	"github.com/oraksil/azumma/internal/domain/usecases"
	"github.com/sangwonl/mqrpc"
)

type HelloHandler struct {
	GameCtrlUseCase *usecases.GameCtrlUseCase
}

func (h *HelloHandler) handleHello(ctx *mqrpc.Context) interface{} {
	var temp map[string]string
	json.Unmarshal(ctx.GetMessage().Payload, &temp)
	fmt.Println(temp)

	return nil
}

func (h *HelloHandler) Routes() []mqrpc.Route {
	return []mqrpc.Route{
		{MsgType: models.MSG_HELLO, Handler: h.handleHello},
	}
}

type SignalingHandler struct {
	SignalingUseCase *usecases.SignalingUseCase
}

func (h *SignalingHandler) handleIceCandidate(ctx *mqrpc.Context) interface{} {
	var temp models.Icecandidate
	json.Unmarshal(ctx.GetMessage().Payload, &temp)

	h.SignalingUseCase.AddServerIceCandidate(temp.OrakkiId, temp.IceString)
	return nil
}

func (h *SignalingHandler) Routes() []mqrpc.Route {
	return []mqrpc.Route{
		{MsgType: models.MSG_HANDLE_SETUP_ICECANDIDATE, Handler: h.handleIceCandidate},
	}
}
