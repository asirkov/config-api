package router

import (
	"fmt"
	"strconv"
	"time"

	"github.com/minelytix/config-api/log"
	"github.com/savsgio/atreugo/v11"
)

type Router struct {
	router *atreugo.Atreugo
}

func NewRouter(port int) *Router {
	config := atreugo.Config{
		Addr:             ":" + strconv.Itoa(port),
		GracefulShutdown: true,
		Logger:           Logger{},
		ReadTimeout:      30 * time.Second,
		WriteTimeout:     30 * time.Second,
	}
	router := atreugo.New(config)

	return &Router{
		router: router,
	}
}

func (r *Router) RegisterController(controllers ...Controller) {
	for _, controller := range controllers {
		getHandlers := controller.GetHandlers()

		for _, handler := range getHandlers {
			if handler.Method == GET {
				r.router.GET(handler.Pattern, handler.Handler)
			} else if handler.Method == PUT {
				r.router.PUT(handler.Pattern, handler.Handler)
			} else if handler.Method == POST {
				r.router.POST(handler.Pattern, handler.Handler)
			} else if handler.Method == DELETE {
				r.router.DELETE(handler.Pattern, handler.Handler)
			}
		}
	}

}

func (r *Router) Listen() error {
	if err := r.router.ListenAndServe(); err != nil {
		return err
	}
	return nil
}

type Logger struct{}

func (l Logger) Print(args ...interface{}) {
	log.Info(fmt.Sprint(args...))
}

func (l Logger) Printf(format string, args ...interface{}) {
	log.Infof(format, args...)
}

func (r *Router) Shutdown() error {
	return nil
}
