package main

import (
	"gitlab.com/oraksil/sil/backend/internal/presenter/di"
)

func main() {
	di.InitContainer()

	mqSvc := di.InjectMqService()
	mqSvc.AddHandler(di.InjectHelloHandler())
	go func() { mqSvc.Run("sil", "") }()

	webSvc := di.InjectWebService()
	webSvc.AddController(di.InjectGameController())
	webSvc.Run("8000")
}
