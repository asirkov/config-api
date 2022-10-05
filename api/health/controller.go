package health

import (
	"io/ioutil"
	"net/http"

	"github.com/minelytix/config-api/api"
	"github.com/minelytix/config-api/log"
	"github.com/minelytix/config-api/router"
	"github.com/savsgio/atreugo/v11"
)

type HealthController struct {
	api HealthApi
}

func NewHealthController(api HealthApi) *HealthController {
	return &HealthController{
		api: api,
	}
}

func (r *HealthController) GetHandlers() []router.Handler {
	return []router.Handler{
		{
			Pattern: "/healthcheck",
			Handler: r.HealthCheck,
			Method:  router.GET,
		},
		{
			Pattern: "/healthpage",
			Handler: r.HealthPage,
			Method:  router.GET,
		},
	}
}

func (c *HealthController) HealthCheck(ctx *atreugo.RequestCtx) error {
	_, failed := c.api.HealthCheck(ctx)
	if failed {
		ctx.Response.Header.Set("Content-Type", "text/plain; charset=utf-8")
		ctx.Response.Header.Set("X-Content-Type-Options", "nosniff")
		ctx.Response.Header.SetStatusCode(http.StatusServiceUnavailable)
		if _, err := ctx.Write([]byte("DOWN")); err != nil {
			log.WarnCtx(ctx, "Utils, HealthCheck, Failed to write: ", err)
			api.HandleError(ctx, err)
			return nil
		}
	} else {

		ctx.Response.Header.Set("Content-Type", "text/plain; charset=utf-8")
		ctx.Response.Header.Set("X-Content-Type-Options", "nosniff")
		ctx.Response.Header.SetStatusCode(http.StatusOK)
		if _, err := ctx.Write([]byte("UP")); err != nil {
			log.WarnCtx(ctx, "Utils, HealthCheck, Failed to write: ", err)
			api.HandleError(ctx, err)
			return nil
		}
	}
	return nil
}

func (c *HealthController) HealthPage(ctx *atreugo.RequestCtx) error {
	statuses, failed := c.api.HealthCheck(ctx)
	version, err := ioutil.ReadFile("VERSION")
	if err != nil {
		version = []byte("?")
	}

	body := "<!DOCTYPE html><html><HEAD><TITLE>Prediction API health page</TITLE></HEAD>"
	body = body + "<body><H1>Prediction API health page</H1><table>"
	body = body + "<tr><td>Version</td><td>" + string(version) + "</td></tr>"

	for _, status := range statuses {
		body = body + "<tr><td>" + status.Service + ":</td><td>"
		if status.Failed {
			body = body + `<span style="color:red;">DOWN</span>`
		} else {
			body = body + `<span style="color:green;">UP</span>`
		}
		body = body + "</td><td>" + status.Message + "</td></tr>"
	}
	body = body + "</table></body></html>"

	if failed {
		ctx.Response.Header.SetStatusCode(http.StatusServiceUnavailable)
	} else {
		ctx.Response.Header.SetStatusCode(http.StatusOK)
	}

	ctx.Response.Header.Set("Content-Type", "text/html; charset=utf-8")
	if _, err := ctx.Write([]byte(body)); err != nil {
		log.WarnCtx(ctx, "Utils, HealthPage, Failed to write health page: ", err)
		api.HandleError(ctx, err)
		return nil
	}

	return nil
}
