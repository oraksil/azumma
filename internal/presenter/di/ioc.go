package di

import (
	"github.com/golobby/container"
	"github.com/sangwonl/mqrpc"
	"gitlab.com/oraksil/azumma/internal/presenter/mq/handlers"
	"gitlab.com/oraksil/azumma/internal/presenter/web"
	"gitlab.com/oraksil/azumma/internal/presenter/web/ctrls"
)

func InitContainer() {
	container.Singleton(newServiceConfig)
	container.Singleton(newOrakkiDriver)
	container.Singleton(newWebService)
	container.Singleton(newMqService)
	container.Singleton(newMessageService)
	container.Singleton(newMySqlDb)
	container.Singleton(newGameRepository)
	container.Singleton(newGameFetchUseCase)
	container.Singleton(newGameCtrlUseCase)
	container.Singleton(newGameController)
	container.Singleton(newHelloHandler)
}

func InjectWebService() *web.WebService {
	var svc *web.WebService
	container.Make(&svc)
	return svc
}

func InjectMqService() *mqrpc.MqService {
	var svc *mqrpc.MqService
	container.Make(&svc)
	return svc
}

func InjectHelloHandler() *handlers.HelloHandler {
	var handler *handlers.HelloHandler
	container.Make(&handler)
	return handler
}

func InjectGameController() *ctrls.GameController {
	var ctrl *ctrls.GameController
	container.Make(&ctrl)
	return ctrl
}
