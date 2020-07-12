package main

import (
	"fmt"

	"oraksil.com/sil/internal/domain/usecases"
	"oraksil.com/sil/internal/presenter/data"
	"oraksil.com/sil/internal/presenter/web"
	"oraksil.com/sil/internal/presenter/web/ctrls"
)

func main() {
	fmt.Print("Hello World")

	// repositories
	gameRepositoryImpl := data.GameRepository{}

	// usecases
	gameFetchUseCase := usecases.GameFetchUseCase{
		GameRepository: &gameRepositoryImpl,
	}

	// controllers
	helloCtrl := ctrls.HelloController{}

	gameCtrl := ctrls.GameController{
		GameFetchUseCase: &gameFetchUseCase,
	}

	w := web.NewWebService()
	w.AddController(&helloCtrl)
	w.AddController(&gameCtrl)
	w.Run("8000")
}
