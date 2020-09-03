package main

import (
	"github.com/joho/godotenv"
	"github.com/oraksil/azumma/internal/presenter/di"
)

func main() {
	godotenv.Load(".env")

	di.InitContainer()

	mqSvc := di.InjectMqService()
	mqSvc.AddHandler(di.InjectHelloHandler())

	conf := di.InjectServiceConfig()
	go func() { mqSvc.Run(conf.PeerName) }()

	webSvc := di.InjectWebService()
	webSvc.AddController(di.InjectGameController())
	webSvc.Run("8000")
}
