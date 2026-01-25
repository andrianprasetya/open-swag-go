package chi

import (
	"github.com/go-chi/chi/v5"

	openswag "github.com/andrianprasetya/open-swag-go"
)

// Mount mounts the documentation on a Chi router
func Mount(r chi.Router, docs *openswag.Docs, basePath string) {
	r.Get(basePath, docs.Handler())
	r.Get(basePath+"/openapi.json", docs.SpecHandler())
}
