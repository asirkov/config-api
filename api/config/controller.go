package config

import (
	"encoding/json"
	"net/http"

	"github.com/minelytix/config-api/api"
	"github.com/minelytix/config-api/router"
	"github.com/savsgio/atreugo/v11"
)

type ConfigController struct {
	api ConfigApi
}

func NewConfigController(api ConfigApi) *ConfigController {
	return &ConfigController{
		api: api,
	}
}

func (r *ConfigController) GetHandlers() []router.Handler {
	return []router.Handler{
		{
			Pattern: "/api/v1/config",
			Handler: r.ListConfigs,
			Method:  router.GET,
		},
		{
			Pattern: "/api/v1/config/{id}",
			Handler: r.GetConfig,
			Method:  router.GET,
		},
		{
			Pattern: "/api/v1/config",
			Handler: r.CreateConfig,
			Method:  router.POST,
		},
		{
			Pattern: "/api/v1/config/{id}",
			Handler: r.UpdateConfig,
			Method:  router.PUT,
		},
		{
			Pattern: "/api/v1/config/{id}",
			Handler: r.DeleteConfig,
			Method:  router.DELETE,
		},
	}
}

func (c *ConfigController) ListConfigs(ctx *atreugo.RequestCtx) error {
	pagination := api.GetPagination(ctx)

	configs, err := c.api.ListConfigs(pagination)
	if err != nil {
		return api.HandleError(ctx, err)
	}

	if bytes, err := json.Marshal(configs); err == nil {
		ctx.Response.Header.SetContentType("application/json")
		ctx.Response.Header.SetStatusCode(http.StatusOK)
		ctx.Response.SetBody(bytes)
	} else {
		return api.HandleError(ctx, err)
	}
	return nil
}

func (c *ConfigController) GetConfig(ctx *atreugo.RequestCtx) error {
	id := ctx.UserValue("id").(string)
	if id == "" {
		return api.HandleError(ctx, &api.WarningValidation{Msg: "missed id in request path"})
	}

	config, err := c.api.GetConfig(id)
	if err != nil {
		return api.HandleError(ctx, err)
	}

	if bytes, err := json.Marshal(config); err == nil {
		ctx.Response.Header.SetContentType("application/json")
		ctx.Response.Header.SetStatusCode(http.StatusOK)
		ctx.Response.SetBody(bytes)
	} else {
		return api.HandleError(ctx, err)
	}
	return nil
}

func (c *ConfigController) CreateConfig(ctx *atreugo.RequestCtx) error {
	var body ConfigDto
	if err := json.Unmarshal(ctx.Request.Body(), &body); err != nil {
		return api.HandleError(ctx, err)
	}

	config, err := c.api.CreateConfig(&body)
	if err != nil {
		return api.HandleError(ctx, err)
	}

	if bytes, err := json.Marshal(config); err == nil {
		ctx.Response.Header.SetContentType("application/json")
		ctx.Response.Header.SetStatusCode(http.StatusOK)
		ctx.Response.SetBody(bytes)
	} else {
		return api.HandleError(ctx, err)
	}
	return nil
}

func (c *ConfigController) UpdateConfig(ctx *atreugo.RequestCtx) error {
	id := ctx.UserValue("id").(string)
	if id == "" {
		return api.HandleError(ctx, &api.WarningValidation{Msg: "missed id in request path"})
	}

	var body ConfigDto
	if err := json.Unmarshal(ctx.Request.Body(), &body); err != nil {
		return api.HandleError(ctx, err)
	}

	config, err := c.api.UpdateConfig(id, &body)
	if err != nil {
		return api.HandleError(ctx, err)
	}

	if bytes, err := json.Marshal(config); err == nil {
		ctx.Response.Header.SetContentType("application/json")
		ctx.Response.Header.SetStatusCode(http.StatusOK)
		ctx.Response.SetBody(bytes)
	} else {
		return api.HandleError(ctx, err)
	}
	return nil
}

func (c *ConfigController) DeleteConfig(ctx *atreugo.RequestCtx) error {
	id := ctx.UserValue("id").(string)
	if id == "" {
		return api.HandleError(ctx, &api.WarningValidation{Msg: "missed id in request path"})
	}

	if err := c.api.DeleteConfig(id); err != nil {
		return api.HandleError(ctx, err)
	}
	ctx.Response.Header.SetContentType("application/json")
	ctx.Response.Header.SetStatusCode(http.StatusNoContent)

	return nil
}
