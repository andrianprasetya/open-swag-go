package fiber

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"

	openswag "github.com/andrianprasetya/open-swag-go"
)

// Mount mounts the documentation on a Fiber app
func Mount(app *fiber.App, docs *openswag.Docs, basePath string) {
	if !strings.HasSuffix(basePath, "/") {
		basePath += "/"
	}
	app.Get(basePath, adaptor.HTTPHandlerFunc(docs.Handler()))
	app.Get(basePath+"openapi.json", adaptor.HTTPHandlerFunc(docs.SpecHandler()))
}

// MountGroup mounts the documentation on a Fiber router group
func MountGroup(g fiber.Router, docs *openswag.Docs) {
	g.Get("/", adaptor.HTTPHandlerFunc(docs.Handler()))
	g.Get("/openapi.json", adaptor.HTTPHandlerFunc(docs.SpecHandler()))
}
