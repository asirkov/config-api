package validation

import (
	"encoding/json"
	"net/http"

	"github.com/minelytix/config-api/api"
	"github.com/minelytix/config-api/router"
	"github.com/savsgio/atreugo/v11"
)

type ValidationController struct {
	api ValidationApi
}

func NewValidationController(api ValidationApi) *ValidationController {
	return &ValidationController{
		api: api,
	}
}

func (r *ValidationController) GetHandlers() []router.Handler {
	return []router.Handler{
		{
			Pattern: "/api/v1/config/{id}/validate",
			Handler: r.ValidateConfig,
			Method:  router.POST,
		},
	}
}

func (c *ValidationController) ValidateConfig(ctx *atreugo.RequestCtx) error {
	id := ctx.UserValue("id").(string)
	if id == "" {
		return api.HandleError(ctx, &api.WarningValidation{Msg: "missed id in request path"})
	}

	var body map[string]interface{}
	if err := json.Unmarshal(ctx.Request.Body(), &body); err != nil {
		return api.HandleError(ctx, err)
	}

	data, err := c.api.ValidateConfig(id, &body)
	if err != nil {
		return api.HandleError(ctx, err)
	}

	if bytes, err := json.Marshal(data); err == nil {
		ctx.Response.Header.SetContentType("application/json")
		ctx.Response.Header.SetStatusCode(http.StatusOK)
		ctx.Response.SetBody(bytes)
	} else {
		return api.HandleError(ctx, err)
	}
	return nil
}
