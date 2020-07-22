package handlers

import (
	"encoding/json"
	"fmt"

	"gitlab.com/oraksil/azumma/internal/domain/models"
	"gitlab.com/oraksil/azumma/internal/domain/usecases"
	"github.com/sangwonl/go-mq-rpc"
)

type HelloHandler struct {
	GameCtrlUseCase *usecases.GameCtrlUseCase
}

func (h *HelloHandler) handleHello(ctx *mq.Context) interface{} {
	var temp map[string]string
	json.Unmarshal(ctx.GetMessage().Payload, &temp)
	fmt.Println(temp)

	return nil
}

func (h *HelloHandler) Routes() []mq.Route {
	return []mq.Route{
		{MsgType: models.MSG_HELLO, Handler: h.handleHello},
	}
}
