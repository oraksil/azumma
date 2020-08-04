package main

import (
	"gitlab.com/oraksil/azumma/internal/presenter/di"
)

func main() {
	di.InitContainer()

	mqSvc := di.InjectMqService()
	mqSvc.AddHandler(di.InjectHelloHandler())
	go func() { mqSvc.Run("sil-temp") }()

	webSvc := di.InjectWebService()
	webSvc.AddController(di.InjectGameController())
	webSvc.AddController(di.InjectSignalingController())
	webSvc.Run("8000")
}
