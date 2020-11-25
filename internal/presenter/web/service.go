package web

import (
	"fmt"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/oraksil/azumma/pkg/utils"
)

type Controller interface {
	Routes() []Route
}

type WebService struct {
	routes      *gin.Engine
	controllers []Controller
}

func getAllowOrigins() []string {
	ginMode := utils.GetStrEnv("GIN_MODE", "")
	if ginMode == "release" {
		return []string{
			"https://oraksil.fun",
			"https://beta.oraksil.fun",
		}
	}

	return []string{"*"}
}

func NewWebService() *WebService {
	routes := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowCredentials = true
	corsConfig.AllowOrigins = getAllowOrigins()
	routes.Use(cors.New(corsConfig))

	// cookie-based session
	store := cookie.NewStore([]byte("423F4528482B4D62"))
	routes.Use(sessions.Sessions("session", store))

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
