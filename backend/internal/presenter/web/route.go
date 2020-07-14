package web

import (
	"github.com/gin-gonic/gin"
)

type Route struct {
	Spec    string
	Handler gin.HandlerFunc
}
