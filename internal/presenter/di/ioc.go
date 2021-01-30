package di

import (
	"github.com/golobby/container"
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
	container.Singleton(newUserFeedbackRepository)
	container.Singleton(newPlayerUseCase)
	container.Singleton(newPlayerController)
	container.Singleton(newGameFetchUseCase)
	container.Singleton(newGameCtrlUseCase)
	container.Singleton(newGameController)
	container.Singleton(newGameHandler)
	container.Singleton(newSignalingUseCases)
	container.Singleton(newSignalingHandler)
	container.Singleton(newSignalingController)
	container.Singleton(newUserFeedbackUseCase)
	container.Singleton(newUserFeedbackController)
}

func Resolve(receiver interface{}) {
	container.Make(receiver)
}
