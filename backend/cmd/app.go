package main

import (
	"fmt"

	"oraksil.com/sil/internal/presenter/web"
	"oraksil.com/sil/internal/presenter/web/ctrls"
)

func main() {
	fmt.Print("Hello World")

	w := web.NewWebService()

	helloCtrl := ctrls.HelloController{}
	w.AddController(&helloCtrl)

	worldCtrl := ctrls.WorldController{}
	w.AddController(&worldCtrl)

	w.Run("8000")
}
