package main

import (
	"oraksil.com/sil/internal/presenter/di"
	"oraksil.com/sil/internal/presenter/web"
)

func main() {
	di.InitContainer()

	w := web.NewWebService()
	w.AddController(di.InjectGameController())

	w.Run("8000")
}
