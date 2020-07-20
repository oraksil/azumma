package di

import (
	"github.com/golobby/container"
	"oraksil.com/sil/internal/presenter/mq/handlers"
	"oraksil.com/sil/internal/presenter/web"
	"oraksil.com/sil/internal/presenter/web/ctrls"
	"oraksil.com/sil/pkg/mq"
)

func InitContainer() {
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

func InjectMqService() *mq.MqService {
	var svc *mq.MqService
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
