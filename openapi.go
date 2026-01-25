package openswag

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
)

// Docs is the main documentation instance
type Docs struct {
	config    Config
	endpoints []Endpoint
	spec      *Spec
	mu        sync.RWMutex
}

// New creates a new documentation instance
func New(config Config) *Docs {
	// Apply defaults
	if config.UI.Theme == "" {
		config.UI.Theme = "purple"
	}
	if config.UI.Layout == "" {
		config.UI.Layout = "modern"
	}

	return &Docs{
		config:    config,
		endpoints: make([]Endpoint, 0),
	}
}

// Add registers an endpoint
func (d *Docs) Add(endpoint Endpoint) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.endpoints = append(d.endpoints, endpoint)
	d.spec = nil // Invalidate cache
}

// AddAll registers multiple endpoints
func (d *Docs) AddAll(endpoints ...Endpoint) {
	for _, ep := range endpoints {
		d.Add(ep)
	}
}

// BuildSpec generates the OpenAPI spec
func (d *Docs) BuildSpec() *Spec {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.spec != nil {
		return d.spec
	}

	spec := &Spec{
		OpenAPI: "3.1.0",
		Info:    d.config.Info,
		Servers: d.config.Servers,
		Tags:    d.config.Tags,
		Paths:   make(map[string]PathItem),
		Components: Components{
			Schemas:         make(map[string]*Schema),
			SecuritySchemes: make(map[string]map[string]interface{}),
		},
	}

	// Add security schemes
	for _, scheme := range d.config.Auth.Schemes {
		spec.Components.SecuritySchemes[scheme.Name] = scheme.ToOpenAPI()
	}

	// Build paths from endpoints
	for _, ep := range d.endpoints {
		d.addEndpointToSpec(spec, ep)
	}

	d.spec = spec
	return spec
}

func (d *Docs) addEndpointToSpec(spec *Spec, ep Endpoint) {
	pathItem, exists := spec.Paths[ep.Path]
	if !exists {
		pathItem = PathItem{}
	}

	operation := d.buildOperation(spec, ep)

	method := strings.ToUpper(ep.Method)
	switch method {
	case "GET":
		pathItem.Get = operation
	case "POST":
		pathItem.Post = operation
	case "PUT":
		pathItem.Put = operation
	case "PATCH":
		pathItem.Patch = operation
	case "DELETE":
		pathItem.Delete = operation
	case "OPTIONS":
		pathItem.Options = operation
	case "HEAD":
		pathItem.Head = operation
	}

	spec.Paths[ep.Path] = pathItem
}

func (d *Docs) buildOperation(spec *Spec, ep Endpoint) *Operation {
	op := &Operation{
		Summary:     ep.Summary,
		Description: ep.Description,
		Tags:        ep.Tags,
		Deprecated:  ep.Deprecated,
		Responses:   make(map[string]ResponseSpec),
	}

	// Build parameters
	for _, param := range ep.Parameters {
		paramSpec := ParameterSpec{
			Name:        param.Name,
			In:          param.In,
			Description: param.Description,
			Required:    param.Required,
			Example:     param.Example,
		}
		if param.Schema != nil {
			paramSpec.Schema = param.Schema
		} else {
			paramSpec.Schema = &Schema{Type: "string"}
		}
		op.Parameters = append(op.Parameters, paramSpec)
	}

	// Build request body
	if ep.RequestBody != nil {
		contentType := ep.RequestBody.ContentType
		if contentType == "" {
			contentType = "application/json"
		}

		var schema *Schema
		if ep.RequestBody.Schema != nil {
			schema = SchemaFromType(ep.RequestBody.Schema)
		}

		op.RequestBody = &RequestBodySpec{
			Description: ep.RequestBody.Description,
			Required:    ep.RequestBody.Required,
			Content: map[string]MediaType{
				contentType: {Schema: schema},
			},
		}
	}

	// Build responses
	for code, resp := range ep.Responses {
		respSpec := ResponseSpec{
			Description: resp.Description,
		}

		if resp.Schema != nil {
			schema := SchemaFromType(resp.Schema)
			respSpec.Content = map[string]MediaType{
				"application/json": {Schema: schema},
			}
		}

		op.Responses[codeToString(code)] = respSpec
	}

	// Build security
	for _, secName := range ep.Security {
		op.Security = append(op.Security, map[string][]string{
			secName: {},
		})
	}

	return op
}

func codeToString(code int) string {
	return fmt.Sprintf("%d", code)
}

// SpecJSON returns the OpenAPI spec as JSON
func (d *Docs) SpecJSON() ([]byte, error) {
	spec := d.BuildSpec()
	return json.MarshalIndent(spec, "", "  ")
}
