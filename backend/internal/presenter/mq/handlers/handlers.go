package handlers

import (
	"encoding/json"
	"fmt"

	"oraksil.com/sil/internal/domain/models"
	"oraksil.com/sil/internal/domain/usecases"
	"oraksil.com/sil/pkg/mq"
)

type HelloHandler struct {
	GameCtrlUseCase *usecases.GameCtrlUseCase
}

func (h *HelloHandler) handleHello(ctx *mq.Context) {
	var temp map[string]string
	json.Unmarshal(ctx.GetMessage().Payload, &temp)
	fmt.Println(temp)
}

func (h *HelloHandler) Routes() []mq.Route {
	return []mq.Route{
		{MsgType: models.MSG_HELLO, Handler: h.handleHello},
	}
}
