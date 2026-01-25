package echo

import (
	"github.com/labstack/echo/v4"

	openswag "github.com/andrianprasetya/open-swag-go"
)

// Mount mounts the documentation on an Echo router
func Mount(e *echo.Echo, docs *openswag.Docs, basePath string) {
	e.GET(basePath, echo.WrapHandler(docs.HandlerHTTP()))
	e.GET(basePath+"/openapi.json", echo.WrapHandler(docs.SpecHandlerHTTP()))
}

// MountGroup mounts the documentation on an Echo group
func MountGroup(g *echo.Group, docs *openswag.Docs) {
	g.GET("", echo.WrapHandler(docs.HandlerHTTP()))
	g.GET("/openapi.json", echo.WrapHandler(docs.SpecHandlerHTTP()))
}
