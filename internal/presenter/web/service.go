package web

import (
	"fmt"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Controller interface {
	Routes() []Route
}

type WebService struct {
	routes      *gin.Engine
	controllers []Controller
}

func NewWebService() *WebService {
	routes := gin.Default()
	routes.Use(cors.Default())

	return &WebService{routes: routes, controllers: nil}
}

func (w *WebService) Run(port string) {
	if port == "" {
		port = os.Getenv("PORT")
	}

	if port == "" {
		port = "8000"
	}

	w.routes.Run(fmt.Sprintf(":%s", port))
}

func (w *WebService) AddController(ctrl Controller) {
	w.controllers = append(w.controllers, ctrl)
	for _, r := range ctrl.Routes() {
		splits := strings.Split(r.Spec, " ")
		w.routes.Handle(splits[0], splits[1], r.Handler)
	}
}