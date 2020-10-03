package handlers

import (

	// "fmt"

	"encoding/json"

	"github.com/oraksil/azumma/internal/domain/models"
	"github.com/oraksil/azumma/internal/domain/usecases"
	"github.com/sangwonl/mqrpc"
)

type GameHandler struct {
	GameCtrlUseCase *usecases.GameCtrlUseCase
}

func (h *GameHandler) handlePlayerJoined(ctx *mqrpc.Context) interface{} {
	var playerPart models.PlayerParticipation
	json.Unmarshal(ctx.GetMessage().Payload, &playerPart)

	h.GameCtrlUseCase.JoinGame(playerPart.GameId, playerPart.PlayerId)

	return nil
}

func (h *GameHandler) handlePlayerJoinFailed(ctx *mqrpc.Context) interface{} {
	return nil
}

func (h *GameHandler) handlePlayerLeft(ctx *mqrpc.Context) interface{} {
	var playerPart models.PlayerParticipation
	json.Unmarshal(ctx.GetMessage().Payload, &playerPart)

	h.GameCtrlUseCase.LeaveGame(playerPart.GameId, playerPart.PlayerId)

	return nil
}

func (h *GameHandler) Routes() []mqrpc.Route {
	return []mqrpc.Route{
		{MsgType: models.MsgPlayerJoined, Handler: h.handlePlayerJoined},
		{MsgType: models.MsgPlayerLeft, Handler: h.handlePlayerLeft},
	}
}

type SignalingHandler struct {
	SignalingUseCase *usecases.SignalingUseCase
}

func (h *SignalingHandler) handleOrakkiIceCandidate(ctx *mqrpc.Context) interface{} {
	var orakkiIce models.IceCandidate
	json.Unmarshal(ctx.GetMessage().Payload, &orakkiIce)

	h.SignalingUseCase.OnOrakkiIceCandidate(orakkiIce)

	return nil
}

func (h *SignalingHandler) Routes() []mqrpc.Route {
	return []mqrpc.Route{
		{MsgType: models.MsgRemoteIceCandidate, Handler: h.handleOrakkiIceCandidate},
	}
}
