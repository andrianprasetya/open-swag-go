package gin

import (
	"strings"

	"github.com/gin-gonic/gin"

	openswag "github.com/andrianprasetya/open-swag-go"
)

// Mount mounts the documentation on a Gin router
func Mount(r *gin.Engine, docs *openswag.Docs, basePath string) {
	if !strings.HasSuffix(basePath, "/") {
		basePath += "/"
	}
	r.GET(basePath, gin.WrapF(docs.Handler()))
	r.GET(basePath+"openapi.json", gin.WrapF(docs.SpecHandler()))
}

// MountGroup mounts the documentation on a Gin router group
func MountGroup(rg *gin.RouterGroup, docs *openswag.Docs) {
	rg.GET("", gin.WrapF(docs.Handler()))
	rg.GET("/openapi.json", gin.WrapF(docs.SpecHandler()))
}
