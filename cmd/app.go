package main

import (
	"github.com/joho/godotenv"
	"github.com/oraksil/azumma/internal/presenter/di"
)

func main() {
	godotenv.Load(".env")

	di.InitContainer()

	mqSvc := di.InjectMqService()
	mqSvc.AddHandler(di.InjectGameHandler())
	mqSvc.AddHandler(di.InjectSignalingHandler())

	conf := di.InjectServiceConfig()
	go func() { mqSvc.Run(conf.MqRpcIdentifier, true) }()

	webSvc := di.InjectWebService()
	webSvc.AddController(di.InjectPlayerController())
	webSvc.AddController(di.InjectGameController())
	webSvc.AddController(di.InjectSignalingController())
	webSvc.Run("8000")
}
