package main

import (
	"fmt"

	"oraksil.com/sil/internal/presenter/di"
	"oraksil.com/sil/internal/presenter/web"
)

func main() {
	fmt.Print("Hello World")

	di.InitContainer()

	w := web.NewWebService()
	w.AddController(di.InjectGameController())

	w.Run("8000")
}
