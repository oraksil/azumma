package ctrls

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"oraksil.com/sil/internal/presenter/web"
)

type HelloController struct {
}

func (ctrl *HelloController) hello(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"say": "hello"})
}

func (ctrl *HelloController) Routes() []web.Route {
	return []web.Route{
		{Method: web.GET, Url: "/hello", Handler: ctrl.hello},
	}
}
