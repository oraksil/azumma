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

func setupRoutes(mqSvc *mqrpc.MqService, routes []handlers.Route) {
	for _, r := range routes {
		err := mqSvc.AddHandler(r.MsgType, r.Handler)
		panic(err)
	}
}

func main() {
	godotenv.Load(".env")

	di.InitContainer()

	di.Resolve(func(
		mqSvc *mqrpc.MqService,
		gameHandler *handlers.GameHandler,
		signalingHandler *handlers.SignalingHandler,
		serviceConf *services.ServiceConfig) {

		setupRoutes(mqSvc, gameHandler.Routes())
		setupRoutes(mqSvc, signalingHandler.Routes())

		go func() { mqSvc.Run(true) }()
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
