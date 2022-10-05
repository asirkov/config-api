package router

import (
	"github.com/savsgio/atreugo/v11"
)

type Method int

const (
	GET Method = iota
	PUT
	POST
	DELETE
)

type Handler struct {
	Pattern string
	Handler atreugo.View
	Method  Method
}

type Controller interface {
	GetHandlers() []Handler
}
