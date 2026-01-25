package openswag

import (
	_ "embed"
	"encoding/json"
	"net/http"
	"strings"
)

//go:embed templates/scalar.html
var scalarTemplate string

// Handler returns the documentation UI handler
func (d *Docs) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Build Scalar configuration
		scalarConfig := map[string]interface{}{
			"theme":       d.config.UI.Theme,
			"layout":      d.config.UI.Layout,
			"darkMode":    d.config.UI.DarkMode,
			"showSidebar": d.config.UI.ShowSidebar,
		}

		configJSON, _ := json.Marshal(scalarConfig)

		// Replace placeholders in template
		html := scalarTemplate
		html = strings.Replace(html, "{{SPEC_URL}}", "./openapi.json", 1)
		html = strings.Replace(html, "{{CONFIG}}", string(configJSON), 1)
		html = strings.Replace(html, "{{TITLE}}", d.config.Info.Title, 1)

		if d.config.UI.CustomCSS != "" {
			html = strings.Replace(html, "{{CUSTOM_CSS}}", d.config.UI.CustomCSS, 1)
		} else {
			html = strings.Replace(html, "{{CUSTOM_CSS}}", "", 1)
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(html))
	}
}

// SpecHandler returns the OpenAPI spec JSON handler
func (d *Docs) SpecHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		specJSON, err := d.SpecJSON()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Write(specJSON)
	}
}

// Mount registers both handlers on a mux
func (d *Docs) Mount(mux *http.ServeMux, basePath string) {
	if !strings.HasSuffix(basePath, "/") {
		basePath += "/"
	}

	mux.HandleFunc(basePath, d.Handler())
	mux.HandleFunc(basePath+"openapi.json", d.SpecHandler())
}

// MountWithRouter for custom routers (chi, gorilla, etc.)
func (d *Docs) MountWithRouter(mountFn func(pattern string, handler http.HandlerFunc), basePath string) {
	mountFn(basePath, d.Handler())
	mountFn(basePath+"/openapi.json", d.SpecHandler())
}

// HandlerHTTP returns the documentation UI as http.Handler
func (d *Docs) HandlerHTTP() http.Handler {
	return http.HandlerFunc(d.Handler())
}

// SpecHandlerHTTP returns the OpenAPI spec as http.Handler
func (d *Docs) SpecHandlerHTTP() http.Handler {
	return http.HandlerFunc(d.SpecHandler())
}
