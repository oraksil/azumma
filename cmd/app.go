package main

import (
	"gitlab.com/oraksil/azumma/internal/presenter/di"
)

func main() {
	di.InitContainer()

	mqSvc := di.InjectMqService()
	mqSvc.AddHandler(di.InjectHelloHandler())

	conf := di.InjectServiceConfig()
	go func() { mqSvc.Run(conf.PeerName) }()

	webSvc := di.InjectWebService()
	webSvc.AddController(di.InjectGameController())
	webSvc.Run("8000")
}
