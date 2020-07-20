package handlers

import (
	"encoding/json"
	"fmt"

	"gitlab.com/oraksil/sil/backend/internal/domain/models"
	"gitlab.com/oraksil/sil/backend/internal/domain/usecases"
	"gitlab.com/oraksil/sil/backend/pkg/mq"
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
