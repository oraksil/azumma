package di

import (
	"github.com/golobby/container"
	"github.com/oraksil/azumma/internal/domain/services"
	"github.com/oraksil/azumma/internal/presenter/mq/handlers"
	"github.com/oraksil/azumma/internal/presenter/web"
	"github.com/oraksil/azumma/internal/presenter/web/ctrls"
	"github.com/sangwonl/mqrpc"
)

func InitContainer() {
	container.Singleton(newServiceConfig)
	container.Singleton(newOrakkiDriver)
	container.Singleton(newWebService)
	container.Singleton(newMqService)
	container.Singleton(newMessageService)
	container.Singleton(newMySqlDb)
	container.Singleton(newPlayerRepository)
	container.Singleton(newPackRepository)
	container.Singleton(newGameRepository)
	container.Singleton(newSignalingRepository)
	container.Singleton(newPlayerUseCase)
	container.Singleton(newPlayerController)
	container.Singleton(newGameFetchUseCase)
	container.Singleton(newGameCtrlUseCase)
	container.Singleton(newGameController)
	container.Singleton(newSignalingUseCases)
	container.Singleton(newSignalingHandler)
	container.Singleton(newSignalingController)
}

func InjectServiceConfig() *services.ServiceConfig {
	var serviceConf *services.ServiceConfig
	container.Make(&serviceConf)
	return serviceConf
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

func InjectSignalingHandler() *handlers.SignalingHandler {
	var handler *handlers.SignalingHandler
	container.Make(&handler)
	return handler
}

func InjectPlayerController() *ctrls.PlayerController {
	var ctrl *ctrls.PlayerController
	container.Make(&ctrl)
	return ctrl
}

func InjectGameController() *ctrls.GameController {
	var ctrl *ctrls.GameController
	container.Make(&ctrl)
	return ctrl
}

func InjectSignalingController() *ctrls.SignalingController {
	var ctrl *ctrls.SignalingController
	container.Make(&ctrl)
	return ctrl
}
