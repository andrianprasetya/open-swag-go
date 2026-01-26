package openswag

import (
	"encoding/json"
	"strings"
	"sync"

	"github.com/andrianprasetya/open-swag-go/pkg/schema"
	"github.com/andrianprasetya/open-swag-go/pkg/spec"
)

// Docs is the main documentation instance
type Docs struct {
	config    Config
	endpoints []Endpoint
	openapi   *spec.OpenAPI
	mu        sync.RWMutex
}

// Endpoint represents an API endpoint definition
type Endpoint struct {
	Method      string
	Path        string
	Summary     string
	Description string
	Tags        []string
	Parameters  []Parameter
	RequestBody *RequestBody
	Responses   map[int]Response
	Security    []string
	Deprecated  bool
}

// Parameter represents an API parameter
type Parameter struct {
	Name        string
	In          string // "path", "query", "header", "cookie"
	Description string
	Required    bool
	Schema      *spec.Schema
	Example     interface{}
}

// RequestBody represents a request body
type RequestBody struct {
	Description string
	Required    bool
	Schema      interface{}
	ContentType string
}

// Response represents an API response
type Response struct {
	Description string
	Schema      interface{}
}

// New creates a new documentation instance
func New(config Config) *Docs {
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
	d.openapi = nil
}

// AddAll registers multiple endpoints
func (d *Docs) AddAll(endpoints ...Endpoint) {
	for _, ep := range endpoints {
		d.Add(ep)
	}
}

// BuildSpec generates the OpenAPI spec
func (d *Docs) BuildSpec() *spec.OpenAPI {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.openapi != nil {
		return d.openapi
	}

	info := spec.NewInfo(d.config.Info.Title, d.config.Info.Version).
		WithDescription(d.config.Info.Description)

	if d.config.Info.Contact != nil {
		info = info.WithContact(
			d.config.Info.Contact.Name,
			d.config.Info.Contact.URL,
			d.config.Info.Contact.Email,
		)
	}

	if d.config.Info.License != nil {
		info = info.WithLicense(d.config.Info.License.Name, d.config.Info.License.URL)
	}

	openapi := spec.NewOpenAPI(info)

	// Add servers
	for _, srv := range d.config.Servers {
		openapi.AddServer(spec.NewServer(srv.URL).WithDescription(srv.Description))
	}

	// Add tags
	for _, tag := range d.config.Tags {
		openapi.AddTag(spec.Tag{Name: tag.Name, Description: tag.Description})
	}

	// Build paths from endpoints
	for _, ep := range d.endpoints {
		d.addEndpointToSpec(openapi, ep)
	}

	d.openapi = openapi
	return openapi
}

func (d *Docs) addEndpointToSpec(openapi *spec.OpenAPI, ep Endpoint) {
	pathItem := openapi.Paths[ep.Path]
	if pathItem == nil {
		pathItem = spec.NewPathItem()
	}

	operation := d.buildOperation(ep)

	method := strings.ToUpper(ep.Method)
	switch method {
	case "GET":
		pathItem.SetGet(operation)
	case "POST":
		pathItem.SetPost(operation)
	case "PUT":
		pathItem.SetPut(operation)
	case "PATCH":
		pathItem.SetPatch(operation)
	case "DELETE":
		pathItem.SetDelete(operation)
	}

	openapi.AddPath(ep.Path, pathItem)
}

func (d *Docs) buildOperation(ep Endpoint) *spec.Operation {
	op := spec.NewOperation(ep.Summary).
		WithDescription(ep.Description).
		WithTags(ep.Tags...).
		SetDeprecated(ep.Deprecated)

	// Build parameters
	for _, param := range ep.Parameters {
		p := spec.NewParameter(param.Name, param.In).
			WithDescription(param.Description).
			SetRequired(param.Required)

		if param.Schema != nil {
			p.WithSchema(param.Schema)
		} else {
			p.WithSchema(spec.NewSchema("string"))
		}

		op.AddParameter(p)
	}

	// Build request body
	if ep.RequestBody != nil {
		contentType := ep.RequestBody.ContentType
		if contentType == "" {
			contentType = "application/json"
		}

		var s *spec.Schema
		if ep.RequestBody.Schema != nil {
			schemaResult := schema.FromType(ep.RequestBody.Schema)
			s = convertSchema(schemaResult)
		}

		rb := spec.NewRequestBody(ep.RequestBody.Description, ep.RequestBody.Required).
			WithJSONContent(s)
		op.WithRequestBody(rb)
	}

	// Build responses
	for code, resp := range ep.Responses {
		r := spec.NewResponse(resp.Description)

		if resp.Schema != nil {
			schemaResult := schema.FromType(resp.Schema)
			s := convertSchema(schemaResult)
			r.WithContent("application/json", s)
		}

		op.AddResponse(intToString(code), r)
	}

	// Build security
	for _, secName := range ep.Security {
		op.WithSecurity(spec.SecurityRequirement{secName: {}})
	}

	return op
}

func convertSchema(s *schema.Schema) *spec.Schema {
	if s == nil {
		return nil
	}

	result := &spec.Schema{
		Type:        s.Type,
		Format:      s.Format,
		Description: s.Description,
		Example:     s.Example,
		Default:     s.Default,
		Enum:        s.Enum,
		Required:    s.Required,
		Pattern:     s.Pattern,
		Minimum:     s.Minimum,
		Maximum:     s.Maximum,
		MinLength:   s.MinLength,
		MaxLength:   s.MaxLength,
	}

	if s.Items != nil {
		result.Items = convertSchema(s.Items)
	}

	if len(s.Properties) > 0 {
		result.Properties = make(map[string]*spec.Schema)
		for k, v := range s.Properties {
			result.Properties[k] = convertSchema(v)
		}
	}

	return result
}

func intToString(n int) string {
	if n == 0 {
		return "0"
	}
	result := ""
	for n > 0 {
		result = string(rune('0'+n%10)) + result
		n /= 10
	}
	return result
}

// SpecJSON returns the OpenAPI spec as JSON
func (d *Docs) SpecJSON() ([]byte, error) {
	openapi := d.BuildSpec()
	return json.MarshalIndent(openapi, "", "  ")
}
