package static

import (
	"github.com/minelytix/config-api/router"
	"github.com/savsgio/atreugo/v11"
)

type StaticController struct {
}

func NewStaticController() *StaticController {
	return &StaticController{}
}

func (r *StaticController) GetHandlers() []router.Handler {
	return []router.Handler{
		{
			Pattern: "/docs/api",
			Handler: r.ApiDocs,
			Method:  router.GET,
		},
		{
			Pattern: "/docs/api/openapi.yaml",
			Handler: r.DocsYaml,
			Method:  router.GET,
		},
		{
			Pattern: "/docs/release-notes.html",
			Handler: r.ReleaseNotes,
			Method:  router.GET,
		},
		{
			Pattern: "/docs",
			Handler: r.Docs,
			Method:  router.GET,
		},
		{
			Pattern: "/docs/index.html",
			Handler: r.Docs,
			Method:  router.GET,
		},
	}
}

func (c *StaticController) ApiDocs(ctx *atreugo.RequestCtx) error {
	ctx.Response.Header.SetContentType("text/html")
	ctx.Response.SendFile("docs/redoc-dynamic.html")
	return nil
}

func (c *StaticController) DocsYaml(ctx *atreugo.RequestCtx) error {
	ctx.Response.Header.SetContentType("text/yaml")
	ctx.Response.SendFile("docs/openapi.yaml")
	return nil
}

func (c *StaticController) ReleaseNotes(ctx *atreugo.RequestCtx) error {
	ctx.Response.Header.SetContentType("text/html")
	ctx.Response.SendFile("docs/release-notes.html")
	return nil
}

func (c *StaticController) Docs(ctx *atreugo.RequestCtx) error {
	ctx.Response.Header.SetContentType("text/html")
	ctx.Response.SendFile("docs/index.html")
	return nil
}
