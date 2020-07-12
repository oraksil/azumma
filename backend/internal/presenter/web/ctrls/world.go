package ctrls

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"oraksil.com/sil/internal/presenter/web"
)

type WorldController struct {
}

func (ctrl *WorldController) world(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"say": "world"})
}

func (ctrl *WorldController) Routes() []web.Route {
	return []web.Route{
		{Method: web.GET, Url: "/world", Handler: ctrl.world},
	}
}
