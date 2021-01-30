package main

import (
	"github.com/joho/godotenv"
	"github.com/oraksil/azumma/internal/domain/services"
	"github.com/oraksil/azumma/internal/presenter/di"
	"github.com/oraksil/azumma/internal/presenter/mq/handlers"
	"github.com/oraksil/azumma/internal/presenter/web"
	"github.com/oraksil/azumma/internal/presenter/web/ctrls"
	"github.com/sangwonl/mqrpc"
)

func main() {
	godotenv.Load(".env")

	di.InitContainer()

	di.Resolve(func(
		mqSvc *mqrpc.MqService,
		gameHandler *handlers.GameHandler,
		signalingHandler *handlers.SignalingHandler,
		serviceConf *services.ServiceConfig) {

		mqSvc.AddHandler(gameHandler)
		mqSvc.AddHandler(signalingHandler)

		go func() { mqSvc.Run(serviceConf.MqRpcIdentifier, true) }()
	})

	di.Resolve(func(
		webSvc *web.WebService,
		playerCtrl *ctrls.PlayerController,
		gameCtrl *ctrls.GameController,
		signalingCtrl *ctrls.SignalingController,
		userFeedbackCtrl *ctrls.UserFeedbackController) {

		webSvc.AddController(playerCtrl)
		webSvc.AddController(gameCtrl)
		webSvc.AddController(signalingCtrl)
		webSvc.AddController(userFeedbackCtrl)

		webSvc.Run("8000")
	})
}
