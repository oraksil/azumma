package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/sangwonl/mqrpc"
	"gitlab.com/oraksil/azumma/internal/domain/models"
	"gitlab.com/oraksil/azumma/internal/domain/usecases"
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
