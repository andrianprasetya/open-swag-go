package openswag

// Spec represents the OpenAPI specification
type Spec struct {
	OpenAPI    string              `json:"openapi"`
	Info       Info                `json:"info"`
	Servers    []Server            `json:"servers,omitempty"`
	Tags       []Tag               `json:"tags,omitempty"`
	Paths      map[string]PathItem `json:"paths"`
	Components Components          `json:"components,omitempty"`
}

// PathItem represents operations on a path
type PathItem struct {
	Get     *Operation `json:"get,omitempty"`
	Post    *Operation `json:"post,omitempty"`
	Put     *Operation `json:"put,omitempty"`
	Patch   *Operation `json:"patch,omitempty"`
	Delete  *Operation `json:"delete,omitempty"`
	Options *Operation `json:"options,omitempty"`
	Head    *Operation `json:"head,omitempty"`
}

// Operation represents an API operation
type Operation struct {
	Summary     string                  `json:"summary,omitempty"`
	Description string                  `json:"description,omitempty"`
	Tags        []string                `json:"tags,omitempty"`
	OperationID string                  `json:"operationId,omitempty"`
	Parameters  []ParameterSpec         `json:"parameters,omitempty"`
	RequestBody *RequestBodySpec        `json:"requestBody,omitempty"`
	Responses   map[string]ResponseSpec `json:"responses"`
	Security    []map[string][]string   `json:"security,omitempty"`
	Deprecated  bool                    `json:"deprecated,omitempty"`
}

// ParameterSpec represents a parameter in the spec
type ParameterSpec struct {
	Name        string      `json:"name"`
	In          string      `json:"in"`
	Description string      `json:"description,omitempty"`
	Required    bool        `json:"required,omitempty"`
	Schema      *Schema     `json:"schema,omitempty"`
	Example     interface{} `json:"example,omitempty"`
}

// RequestBodySpec represents a request body in the spec
type RequestBodySpec struct {
	Description string               `json:"description,omitempty"`
	Required    bool                 `json:"required,omitempty"`
	Content     map[string]MediaType `json:"content"`
}

// MediaType represents a media type
type MediaType struct {
	Schema  *Schema     `json:"schema,omitempty"`
	Example interface{} `json:"example,omitempty"`
}

// ResponseSpec represents a response in the spec
type ResponseSpec struct {
	Description string                `json:"description"`
	Content     map[string]MediaType  `json:"content,omitempty"`
	Headers     map[string]HeaderSpec `json:"headers,omitempty"`
}

// HeaderSpec represents a header in the spec
type HeaderSpec struct {
	Description string  `json:"description,omitempty"`
	Schema      *Schema `json:"schema,omitempty"`
}

// Components represents the components section
type Components struct {
	Schemas         map[string]*Schema                `json:"schemas,omitempty"`
	SecuritySchemes map[string]map[string]interface{} `json:"securitySchemes,omitempty"`
}
