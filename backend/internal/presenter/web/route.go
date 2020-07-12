package web

import (
	"github.com/gin-gonic/gin"
)

type Method int

const (
	GET Method = iota
	POST
	PUT
	DELETE
)

func (m Method) toString() string {
	return [...]string{"GET", "POST", "PUT", "DELETE"}[m]
}

type Route struct {
	Method  Method
	Url     string
	Handler gin.HandlerFunc
}
