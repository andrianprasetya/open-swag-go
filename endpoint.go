package openswag

import "net/http"

// Endpoint represents an API endpoint definition
type Endpoint struct {
	Method      string
	Path        string
	Handler     http.HandlerFunc
	Summary     string
	Description string
	Tags        []string
	Parameters  []Parameter
	RequestBody *RequestBody
	Responses   Responses
	Security    []string
	Deprecated  bool
}

// Parameter represents a path/query/header parameter
type Parameter struct {
	Name        string      `json:"name"`
	In          string      `json:"in"`
	Description string      `json:"description,omitempty"`
	Required    bool        `json:"required,omitempty"`
	Schema      *Schema     `json:"schema,omitempty"`
	Example     interface{} `json:"example,omitempty"`
}

// Responses is a map of status code to response
type Responses map[int]ResponseDef

// ResponseDef represents a response definition
type ResponseDef struct {
	Description string
	Schema      interface{}
	Headers     map[string]Header
}

// Header represents a response header
type Header struct {
	Description string  `json:"description,omitempty"`
	Schema      *Schema `json:"schema,omitempty"`
}

// RequestBody represents request body
type RequestBody struct {
	Description string
	Required    bool
	Schema      interface{}
	ContentType string
}

// Body creates a JSON request body
func Body(schema interface{}) *RequestBody {
	return &RequestBody{
		Required:    true,
		Schema:      schema,
		ContentType: "application/json",
	}
}

// BodyWithDesc creates a JSON request body with description
func BodyWithDesc(description string, schema interface{}) *RequestBody {
	return &RequestBody{
		Description: description,
		Required:    true,
		Schema:      schema,
		ContentType: "application/json",
	}
}

// FormBody creates a multipart form request body
func FormBody(schema interface{}) *RequestBody {
	return &RequestBody{
		Required:    true,
		Schema:      schema,
		ContentType: "multipart/form-data",
	}
}

// Response creates a response definition
func Response(description string, schema interface{}) ResponseDef {
	return ResponseDef{
		Description: description,
		Schema:      schema,
	}
}

// PathParam creates a path parameter
func PathParam(name, description string) Parameter {
	return Parameter{
		Name:        name,
		In:          "path",
		Description: description,
		Required:    true,
	}
}

// QueryParam creates a query parameter
func QueryParam(name, description string) Parameter {
	return Parameter{
		Name:        name,
		In:          "query",
		Description: description,
		Required:    false,
	}
}

// RequiredQueryParam creates a required query parameter
func RequiredQueryParam(name, description string) Parameter {
	return Parameter{
		Name:        name,
		In:          "query",
		Description: description,
		Required:    true,
	}
}

// HeaderParam creates a header parameter
func HeaderParam(name, description string) Parameter {
	return Parameter{
		Name:        name,
		In:          "header",
		Description: description,
		Required:    false,
	}
}

// RequiredHeaderParam creates a required header parameter
func RequiredHeaderParam(name, description string) Parameter {
	return Parameter{
		Name:        name,
		In:          "header",
		Description: description,
		Required:    true,
	}
}
