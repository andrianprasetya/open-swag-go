# Go OpenAPI UI - Development Walkthrough

A comprehensive guide to building a modern OpenAPI documentation package for Go with **decorator-like behavior** similar to `adonis-open-swagger`.

## How adonis-open-swagger Works (Reference)

```typescript
// In AdonisJS - decorators on controller methods
@SwaggerInfo({
    summary: "Create user",
    description: "Create a new user account",
    tags: ["Users"],
})
@SwaggerRequestBody("User data", createUserSchema)
@SwaggerResponse(200, "Success", userResponseSchema)
@SwaggerResponse(400, "Bad request", errorSchema)
@SwaggerSecurity([{ "bearerAuth": [] }])
public async create({ request, response }: HttpContext) {
    // handler code
}
```

## Go Equivalent Pattern (Struct-Based)

Since Go doesn't have decorators, we use **struct-based definitions**. But we want the swagger definition **co-located with the handler** (like AdonisJS decorators).

### Option A: Define swagger inside handler file (Recommended)
```go
// internal/handlers/user.go
package handlers

import (
    "net/http"
    openapi "github.com/andrianprasetya/open-swag-go"
    "myapp/internal/dto"
)

// CreateUser handler
func CreateUser(w http.ResponseWriter, r *http.Request) {
    // Your implementation
}

// CreateUserDoc - swagger definition co-located with handler
var CreateUserDoc = openapi.Endpoint{
    Method:      "POST",
    Path:        "/users",
    Handler:     CreateUser,
    Summary:     "Create a new user",
    Description: "Create a new user account",
    Tags:        []string{"Users"},
    RequestBody: openapi.Body(dto.CreateUserRequest{}),
    Responses: openapi.Responses{
        201: openapi.Response("User created", dto.UserResponse{}),
        400: openapi.Response("Bad request", dto.ErrorResponse{}),
    },
    Security: []string{"bearerAuth"},
}

// GetUser handler
func GetUser(w http.ResponseWriter, r *http.Request) {
    // Your implementation
}

// GetUserDoc - swagger definition
var GetUserDoc = openapi.Endpoint{
    Method:      "GET",
    Path:        "/users/{id}",
    Handler:     GetUser,
    Summary:     "Get user by ID",
    Tags:        []string{"Users"},
    Parameters:  []openapi.Parameter{openapi.PathParam("id", "User ID")},
    Responses: openapi.Responses{
        200: openapi.Response("User found", dto.UserResponse{}),
        404: openapi.Response("Not found", dto.ErrorResponse{}),
    },
    Security: []string{"bearerAuth"},
}

// AllUserEndpoints - export all endpoints from this handler
var AllUserEndpoints = []openapi.Endpoint{
    CreateUserDoc,
    GetUserDoc,
}
```

### Option B: Use a wrapper function
```go
// internal/handlers/user.go
package handlers

import (
    "net/http"
    openapi "github.com/andrianprasetya/open-swag-go"
    "myapp/internal/dto"
)

// CreateUser returns both handler and its swagger doc
func CreateUser() openapi.Endpoint {
    return openapi.Endpoint{
        Method:      "POST",
        Path:        "/users",
        Handler:     createUserHandler,
        Summary:     "Create a new user",
        Tags:        []string{"Users"},
        RequestBody: openapi.Body(dto.CreateUserRequest{}),
        Responses: openapi.Responses{
            201: openapi.Response("User created", dto.UserResponse{}),
        },
        Security: []string{"bearerAuth"},
    }
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
    // Your implementation
}
```

### Why Co-locate?
- ✅ Swagger definition is next to the handler code
- ✅ Easy to update both together
- ✅ Similar DX to AdonisJS decorators
- ✅ IDE can navigate between handler and its doc

## Table of Contents

1. [Phase 0: GitHub Setup](#phase-0-github-setup)
2. [Project Overview](#project-overview)
3. [UI Choice: Scalar](#ui-choice-scalar)
4. [Project Structure](#project-structure)
5. [Phase 1: Core Foundation](#phase-1-core-foundation)
6. [Phase 2: Schema & Examples](#phase-2-schema--examples)
7. [Phase 3: Auth Playground](#phase-3-auth-playground)
8. [Phase 4: Try-It Console](#phase-4-try-it-console)
9. [Phase 5: Version Tracking](#phase-5-version-tracking)
10. [Phase 6: Framework Adapters](#phase-6-framework-adapters)
11. [Phase 7: Publishing & Versioning](#phase-7-publishing--versioning)
12. [Real Project Integration](#real-project-integration)

---

## Phase 0: GitHub Setup

### Step 0.1: Create GitHub Repository

```bash
# Create new repository on GitHub: andrianprasetya/open-swag-go
# Then clone it locally
git clone https://github.com/andrianprasetya/open-swag-go.git
cd open-swag-go
```

### Step 0.2: Initialize Go Module

```bash
# Initialize with your GitHub path
go mod init github.com/andrianprasetya/open-swag-go

# Create initial structure
mkdir -p pkg/{spec,schema,examples,auth,tryit,versioning,ui/templates}
mkdir -p adapters/{chi,gin,echo,fiber}
mkdir -p examples/{basic,full-featured}
```

### Step 0.3: Create Essential Files

**README.md**
```markdown
# Open Swag Go

A modern OpenAPI documentation package for Go with decorator-like DX.

## Features

- 🎨 Modern UI (Scalar) with dark mode
- 📝 Struct-based endpoint definitions (co-located with handlers)
- 🔐 Auth playground with credential persistence
- 🧪 Try-it console with code snippets
- 📊 Auto-generated examples from structs
- 🔄 Version diff & breaking change detection

## Installation

\`\`\`bash
go get github.com/andrianprasetya/open-swag-go
\`\`\`

## Quick Start

\`\`\`go
package main

import (
    "net/http"
    openapi "github.com/andrianprasetya/open-swag-go"
)

func main() {
    docs := openapi.New(openapi.Config{
        Info: openapi.Info{
            Title:   "My API",
            Version: "1.0.0",
        },
    })
    
    docs.Add(openapi.Endpoint{
        Method:  "GET",
        Path:    "/users",
        Summary: "List users",
    })
    
    mux := http.NewServeMux()
    docs.Mount(mux, "/docs")
    
    http.ListenAndServe(":8080", mux)
}
\`\`\`

## License

MIT
```

**LICENSE** (MIT)
```
MIT License

Copyright (c) 2024 Andrian Prasetya

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```

**.gitignore**
```
# Binaries
*.exe
*.exe~
*.dll
*.so
*.dylib

# Test binary
*.test

# Output
*.out

# Go workspace
go.work

# IDE
.idea/
.vscode/
*.swp
*.swo

# OS
.DS_Store
Thumbs.db

# Vendor (optional)
# vendor/
```

### Step 0.4: Initial Commit

```bash
git add .
git commit -m "Initial commit: project structure"
git push origin main
```

---

## Project Overview

### Goals
- **Decorator-like DX** similar to adonis-open-swagger
- Better DX than swaggo (no comment parsing)
- Modern UI with dark mode, responsive design
- Auth playground with credential persistence
- Auto-generated examples from structs
- Built-in API tester with code snippets
- Version diff with breaking change detection

### Tech Stack
- **Language**: Go 1.22+
- **UI**: Scalar (embedded via CDN or bundled)
- **Spec**: OpenAPI 3.1
- **Schema**: JSON Schema draft 2020-12

---

## UI Choice: Scalar

### Why Scalar?

```
✅ Modern, clean design
✅ Built-in dark mode
✅ Responsive (mobile-friendly)
✅ Sidebar with search
✅ Tag grouping
✅ Collapsible schemas
✅ Try-it console built-in
✅ Multiple themes
✅ Active development
✅ Easy to embed (CDN or npm)
✅ Customizable via CSS
```

### How to Embed Scalar

**Option 1: CDN (Simplest)**
```html
<!DOCTYPE html>
<html>
<head>
  <title>API Documentation</title>
  <meta charset="utf-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1" />
</head>
<body>
  <script
    id="api-reference"
    data-url="/openapi.json"
    data-configuration='{
      "theme": "purple",
      "layout": "modern",
      "darkMode": true,
      "showSidebar": true
    }'>
  </script>
  <script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
</body>
</html>
```

**Option 2: Embed in Go (Recommended)**
```go
// Embed the HTML template
//go:embed templates/scalar.html
var scalarHTML string

func (d *Docs) Handler() http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        html := strings.Replace(scalarHTML, "{{SPEC_URL}}", d.specURL, 1)
        html = strings.Replace(html, "{{CONFIG}}", d.configJSON, 1)
        w.Header().Set("Content-Type", "text/html")
        w.Write([]byte(html))
    })
}
```

---

## Project Structure

```
open-swag-go/
├── go.mod
├── go.sum
├── README.md
├── LICENSE
│
├── openapi.go                 # Main entry point
├── config.go                  # Configuration structs
├── handler.go                 # HTTP handlers
│
├── pkg/
│   ├── spec/
│   │   ├── openapi.go         # OpenAPI spec builder
│   │   ├── info.go            # Info object
│   │   ├── paths.go           # Paths & operations
│   │   ├── components.go      # Components (schemas, security)
│   │   └── server.go          # Server objects
│   │
│   ├── schema/
│   │   ├── converter.go       # Go struct → JSON Schema
│   │   ├── tags.go            # Struct tag parser
│   │   ├── types.go           # Type mappings
│   │   └── validator.go       # Schema validation
│   │
│   ├── examples/
│   │   ├── generator.go       # Example generator
│   │   ├── faker.go           # Fake data integration
│   │   └── templates.go       # Example templates
│   │
│   ├── auth/
│   │   ├── schemes.go         # Security schemes
│   │   ├── playground.go      # Auth playground config
│   │   └── persist.go         # Credential persistence
│   │
│   ├── tryit/
│   │   ├── console.go         # Try-it configuration
│   │   ├── snippets/
│   │   │   ├── curl.go
│   │   │   ├── javascript.go
│   │   │   ├── go.go
│   │   │   └── python.go
│   │   ├── history.go         # Request history
│   │   └── environments.go    # Environment variables
│   │
│   ├── versioning/
│   │   ├── diff.go            # Spec differ
│   │   ├── breaking.go        # Breaking change detection
│   │   ├── changelog.go       # Changelog generator
│   │   └── migration.go       # Migration guide
│   │
│   └── ui/
│       ├── scalar.go          # Scalar integration
│       ├── themes.go          # Theme configuration
│       ├── templates/
│       │   └── scalar.html    # Embedded HTML template
│       └── assets/            # Static assets (optional)
│
├── adapters/
│   ├── nethttp/
│   │   └── adapter.go         # net/http adapter
│   ├── chi/
│   │   └── adapter.go         # Chi router adapter
│   ├── gin/
│   │   └── adapter.go         # Gin adapter
│   ├── echo/
│   │   └── adapter.go         # Echo adapter
│   └── fiber/
│       └── adapter.go         # Fiber adapter
│
└── examples/
    ├── basic/
    │   └── main.go
    ├── with-auth/
    │   └── main.go
    ├── version-diff/
    │   └── main.go
    └── full-featured/
        └── main.go
```

---

## Phase 1: Core Foundation

### Step 1.1: Configuration Structs

```go
// config.go
package openapi

// Config is the main configuration for the documentation
type Config struct {
    // OpenAPI Info
    Info Info `json:"info"`
    
    // Servers
    Servers []Server `json:"servers,omitempty"`
    
    // Tags for grouping
    Tags []Tag `json:"tags,omitempty"`
    
    // UI Configuration
    UI UIConfig `json:"ui"`
    
    // Auth Configuration
    Auth AuthConfig `json:"auth,omitempty"`
    
    // Examples Configuration
    Examples ExampleConfig `json:"examples,omitempty"`
    
    // Try-It Configuration
    TryIt TryItConfig `json:"tryIt,omitempty"`
    
    // Versioning Configuration
    Versioning VersionConfig `json:"versioning,omitempty"`
}

type Info struct {
    Title          string  `json:"title"`
    Version        string  `json:"version"`
    Description    string  `json:"description,omitempty"`
    TermsOfService string  `json:"termsOfService,omitempty"`
    Contact        Contact `json:"contact,omitempty"`
    License        License `json:"license,omitempty"`
}

type Contact struct {
    Name  string `json:"name,omitempty"`
    URL   string `json:"url,omitempty"`
    Email string `json:"email,omitempty"`
}

type License struct {
    Name string `json:"name"`
    URL  string `json:"url,omitempty"`
}

type Server struct {
    URL         string `json:"url"`
    Description string `json:"description,omitempty"`
}

type Tag struct {
    Name        string `json:"name"`
    Description string `json:"description,omitempty"`
}

type UIConfig struct {
    Theme              string `json:"theme"`              // "dark", "light", "purple", etc.
    DarkMode           bool   `json:"darkMode"`
    Layout             string `json:"layout"`             // "modern", "classic"
    ShowSidebar        bool   `json:"showSidebar"`
    SidebarSearch      bool   `json:"sidebarSearch"`
    TagGrouping        bool   `json:"tagGrouping"`
    CollapsibleSchemas bool   `json:"collapsibleSchemas"`
    CustomCSS          string `json:"customCss,omitempty"`
}
```

### Step 1.2: Main Entry Point

```go
// openapi.go
package openapi

import (
    "encoding/json"
    "net/http"
    "sync"
)

// Docs is the main documentation instance
type Docs struct {
    config    Config
    endpoints []Endpoint
    spec      *OpenAPISpec
    specJSON  []byte
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
func (d *Docs) Add(endpoint Endpoint) *Docs {
    d.mu.Lock()
    defer d.mu.Unlock()
    d.endpoints = append(d.endpoints, endpoint)
    d.spec = nil // Invalidate cache
    return d
}

// Endpoint creates a new endpoint builder
func (d *Docs) Endpoint(path, method string) *EndpointBuilder {
    return &EndpointBuilder{
        docs: d,
        endpoint: Endpoint{
            Path:   path,
            Method: method,
        },
    }
}

// Build generates the OpenAPI spec
func (d *Docs) Build() (*OpenAPISpec, error) {
    d.mu.Lock()
    defer d.mu.Unlock()
    
    if d.spec != nil {
        return d.spec, nil
    }
    
    spec := &OpenAPISpec{
        OpenAPI: "3.1.0",
        Info:    d.config.Info,
        Servers: d.config.Servers,
        Tags:    d.config.Tags,
        Paths:   make(map[string]PathItem),
        Components: Components{
            Schemas:         make(map[string]Schema),
            SecuritySchemes: make(map[string]SecurityScheme),
        },
    }
    
    // Build paths from endpoints
    for _, ep := range d.endpoints {
        d.addEndpointToSpec(spec, ep)
    }
    
    d.spec = spec
    return spec, nil
}

// SpecJSON returns the OpenAPI spec as JSON
func (d *Docs) SpecJSON() ([]byte, error) {
    spec, err := d.Build()
    if err != nil {
        return nil, err
    }
    return json.MarshalIndent(spec, "", "  ")
}
```

### Step 1.3: Endpoint Builder (Fluent API)

```go
// endpoint.go
package openapi

// Endpoint represents an API endpoint
type Endpoint struct {
    Path        string
    Method      string
    Summary     string
    Description string
    Tags        []string
    Parameters  []Parameter
    RequestBody *RequestBody
    Responses   map[string]Response
    Security    []map[string][]string
    Deprecated  bool
}

// EndpointBuilder provides fluent API for building endpoints
type EndpointBuilder struct {
    docs     *Docs
    endpoint Endpoint
}

func (b *EndpointBuilder) Summary(summary string) *EndpointBuilder {
    b.endpoint.Summary = summary
    return b
}

func (b *EndpointBuilder) Description(desc string) *EndpointBuilder {
    b.endpoint.Description = desc
    return b
}

func (b *EndpointBuilder) Tags(tags ...string) *EndpointBuilder {
    b.endpoint.Tags = tags
    return b
}

func (b *EndpointBuilder) Param(param Parameter) *EndpointBuilder {
    b.endpoint.Parameters = append(b.endpoint.Parameters, param)
    return b
}

// PathParam creates a path parameter
func PathParam(name string) *ParameterBuilder {
    return &ParameterBuilder{
        param: Parameter{
            Name:     name,
            In:       "path",
            Required: true,
        },
    }
}

// QueryParam creates a query parameter
func QueryParam(name string) *ParameterBuilder {
    return &ParameterBuilder{
        param: Parameter{
            Name: name,
            In:   "query",
        },
    }
}

// HeaderParam creates a header parameter
func HeaderParam(name string) *ParameterBuilder {
    return &ParameterBuilder{
        param: Parameter{
            Name: name,
            In:   "header",
        },
    }
}

func (b *EndpointBuilder) Body(schema interface{}) *EndpointBuilder {
    b.endpoint.RequestBody = &RequestBody{
        Required: true,
        Content: map[string]MediaType{
            "application/json": {
                Schema: schemaFromType(schema),
            },
        },
    }
    return b
}

func (b *EndpointBuilder) Response(code int, description string, schema interface{}) *EndpointBuilder {
    if b.endpoint.Responses == nil {
        b.endpoint.Responses = make(map[string]Response)
    }
    
    resp := Response{Description: description}
    if schema != nil {
        resp.Content = map[string]MediaType{
            "application/json": {
                Schema: schemaFromType(schema),
            },
        }
    }
    
    b.endpoint.Responses[fmt.Sprintf("%d", code)] = resp
    return b
}

func (b *EndpointBuilder) Security(schemes ...string) *EndpointBuilder {
    for _, scheme := range schemes {
        b.endpoint.Security = append(b.endpoint.Security, map[string][]string{
            scheme: {},
        })
    }
    return b
}

func (b *EndpointBuilder) Deprecated() *EndpointBuilder {
    b.endpoint.Deprecated = true
    return b
}

// Register adds the endpoint to the docs
func (b *EndpointBuilder) Register() *Docs {
    return b.docs.Add(b.endpoint)
}
```

### Step 1.4: HTTP Handlers

```go
// handler.go
package openapi

import (
    _ "embed"
    "encoding/json"
    "net/http"
    "strings"
)

//go:embed pkg/ui/templates/scalar.html
var scalarTemplate string

// Handler returns the documentation UI handler
func (d *Docs) Handler() http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Build Scalar configuration
        scalarConfig := map[string]interface{}{
            "theme":      d.config.UI.Theme,
            "layout":     d.config.UI.Layout,
            "darkMode":   d.config.UI.DarkMode,
            "showSidebar": d.config.UI.ShowSidebar,
        }
        
        configJSON, _ := json.Marshal(scalarConfig)
        
        // Replace placeholders in template
        html := scalarTemplate
        html = strings.Replace(html, "{{SPEC_URL}}", "/openapi.json", 1)
        html = strings.Replace(html, "{{CONFIG}}", string(configJSON), 1)
        html = strings.Replace(html, "{{TITLE}}", d.config.Info.Title, 1)
        
        if d.config.UI.CustomCSS != "" {
            html = strings.Replace(html, "{{CUSTOM_CSS}}", d.config.UI.CustomCSS, 1)
        } else {
            html = strings.Replace(html, "{{CUSTOM_CSS}}", "", 1)
        }
        
        w.Header().Set("Content-Type", "text/html; charset=utf-8")
        w.Write([]byte(html))
    })
}

// SpecHandler returns the OpenAPI spec JSON handler
func (d *Docs) SpecHandler() http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        specJSON, err := d.SpecJSON()
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        
        w.Header().Set("Content-Type", "application/json")
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Write(specJSON)
    })
}

// Mount registers both handlers on a mux
func (d *Docs) Mount(mux *http.ServeMux, basePath string) {
    if !strings.HasSuffix(basePath, "/") {
        basePath += "/"
    }
    
    mux.Handle(basePath, d.Handler())
    mux.Handle(basePath+"openapi.json", d.SpecHandler())
}
```

### Step 1.5: Scalar HTML Template

```html
<!-- pkg/ui/templates/scalar.html -->
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{TITLE}} - API Documentation</title>
    <style>
        body { margin: 0; padding: 0; }
        {{CUSTOM_CSS}}
    </style>
</head>
<body>
    <script
        id="api-reference"
        data-url="{{SPEC_URL}}"
        data-configuration='{{CONFIG}}'>
    </script>
    <script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
</body>
</html>
```

---

## Phase 2: Schema & Examples

### Step 2.1: Struct to JSON Schema Converter

```go
// pkg/schema/converter.go
package schema

import (
    "reflect"
    "strings"
)

// Schema represents a JSON Schema
type Schema struct {
    Type        string             `json:"type,omitempty"`
    Format      string             `json:"format,omitempty"`
    Description string             `json:"description,omitempty"`
    Properties  map[string]*Schema `json:"properties,omitempty"`
    Required    []string           `json:"required,omitempty"`
    Items       *Schema            `json:"items,omitempty"`
    Enum        []interface{}      `json:"enum,omitempty"`
    Example     interface{}        `json:"example,omitempty"`
    Default     interface{}        `json:"default,omitempty"`
    Minimum     *float64           `json:"minimum,omitempty"`
    Maximum     *float64           `json:"maximum,omitempty"`
    MinLength   *int               `json:"minLength,omitempty"`
    MaxLength   *int               `json:"maxLength,omitempty"`
    Pattern     string             `json:"pattern,omitempty"`
    Ref         string             `json:"$ref,omitempty"`
}

// FromType converts a Go type to JSON Schema
func FromType(t interface{}) *Schema {
    return fromReflectType(reflect.TypeOf(t))
}

func fromReflectType(t reflect.Type) *Schema {
    // Handle pointer types
    if t.Kind() == reflect.Ptr {
        return fromReflectType(t.Elem())
    }
    
    switch t.Kind() {
    case reflect.String:
        return &Schema{Type: "string"}
    case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
        return &Schema{Type: "integer"}
    case reflect.Float32, reflect.Float64:
        return &Schema{Type: "number"}
    case reflect.Bool:
        return &Schema{Type: "boolean"}
    case reflect.Slice, reflect.Array:
        return &Schema{
            Type:  "array",
            Items: fromReflectType(t.Elem()),
        }
    case reflect.Struct:
        return fromStruct(t)
    case reflect.Map:
        return &Schema{
            Type: "object",
            // additionalProperties could be added here
        }
    default:
        return &Schema{Type: "object"}
    }
}

func fromStruct(t reflect.Type) *Schema {
    schema := &Schema{
        Type:       "object",
        Properties: make(map[string]*Schema),
        Required:   []string{},
    }
    
    for i := 0; i < t.NumField(); i++ {
        field := t.Field(i)
        
        // Skip unexported fields
        if !field.IsExported() {
            continue
        }
        
        // Get JSON tag
        jsonTag := field.Tag.Get("json")
        if jsonTag == "-" {
            continue
        }
        
        name := strings.Split(jsonTag, ",")[0]
        if name == "" {
            name = field.Name
        }
        
        // Parse swagger/openapi tags
        fieldSchema := fromReflectType(field.Type)
        parseFieldTags(field, fieldSchema)
        
        schema.Properties[name] = fieldSchema
        
        // Check if required
        if isRequired(field) {
            schema.Required = append(schema.Required, name)
        }
    }
    
    return schema
}

func parseFieldTags(field reflect.StructField, schema *Schema) {
    // Parse example tag
    if example := field.Tag.Get("example"); example != "" {
        schema.Example = example
    }
    
    // Parse description tag
    if desc := field.Tag.Get("description"); desc != "" {
        schema.Description = desc
    }
    
    // Parse format tag
    if format := field.Tag.Get("format"); format != "" {
        schema.Format = format
    }
    
    // Parse swagger tag for additional options
    if swagger := field.Tag.Get("swagger"); swagger != "" {
        parseSwaggerTag(swagger, schema)
    }
}

func parseSwaggerTag(tag string, schema *Schema) {
    parts := strings.Split(tag, ",")
    for _, part := range parts {
        kv := strings.SplitN(part, "=", 2)
        key := strings.TrimSpace(kv[0])
        
        switch key {
        case "required":
            // Handled separately
        case "format":
            if len(kv) > 1 {
                schema.Format = kv[1]
            }
        case "description":
            if len(kv) > 1 {
                schema.Description = kv[1]
            }
        case "example":
            if len(kv) > 1 {
                schema.Example = kv[1]
            }
        }
    }
}

func isRequired(field reflect.StructField) bool {
    // Check swagger tag
    if swagger := field.Tag.Get("swagger"); strings.Contains(swagger, "required") {
        return true
    }
    
    // Check validate tag (for validator libraries)
    if validate := field.Tag.Get("validate"); strings.Contains(validate, "required") {
        return true
    }
    
    // Check binding tag (for Gin)
    if binding := field.Tag.Get("binding"); strings.Contains(binding, "required") {
        return true
    }
    
    return false
}
```

### Step 2.2: Example Generator

```go
// pkg/examples/generator.go
package examples

import (
    "reflect"
    "time"
    
    "github.com/go-faker/faker/v4"
)

// Config for example generation
type Config struct {
    UseFaker     bool
    TypeExamples map[string]interface{}
}

// Generator generates example values
type Generator struct {
    config Config
}

// New creates a new example generator
func New(config Config) *Generator {
    if config.TypeExamples == nil {
        config.TypeExamples = defaultTypeExamples()
    }
    return &Generator{config: config}
}

func defaultTypeExamples() map[string]interface{} {
    return map[string]interface{}{
        "uuid":      "550e8400-e29b-41d4-a716-446655440000",
        "email":     "user@example.com",
        "uri":       "https://example.com",
        "hostname":  "api.example.com",
        "ipv4":      "192.168.1.1",
        "date":      "2024-01-15",
        "date-time": "2024-01-15T10:30:00Z",
    }
}

// Generate creates an example from a Go type
func (g *Generator) Generate(t interface{}) interface{} {
    return g.generateFromType(reflect.TypeOf(t))
}

func (g *Generator) generateFromType(t reflect.Type) interface{} {
    if t.Kind() == reflect.Ptr {
        return g.generateFromType(t.Elem())
    }
    
    switch t.Kind() {
    case reflect.String:
        if g.config.UseFaker {
            return faker.Word()
        }
        return "string"
    case reflect.Int, reflect.Int32, reflect.Int64:
        return 42
    case reflect.Float32, reflect.Float64:
        return 3.14
    case reflect.Bool:
        return true
    case reflect.Slice:
        elem := g.generateFromType(t.Elem())
        return []interface{}{elem}
    case reflect.Struct:
        return g.generateFromStruct(t)
    default:
        return nil
    }
}

func (g *Generator) generateFromStruct(t reflect.Type) map[string]interface{} {
    result := make(map[string]interface{})
    
    for i := 0; i < t.NumField(); i++ {
        field := t.Field(i)
        if !field.IsExported() {
            continue
        }
        
        // Get field name from json tag
        name := getJSONFieldName(field)
        if name == "-" {
            continue
        }
        
        // Check for explicit example tag
        if example := field.Tag.Get("example"); example != "" {
            result[name] = example
            continue
        }
        
        // Check for faker tag
        if g.config.UseFaker {
            if fakerTag := field.Tag.Get("faker"); fakerTag != "" {
                result[name] = g.generateFromFakerTag(fakerTag)
                continue
            }
        }
        
        // Check format for type examples
        if format := field.Tag.Get("format"); format != "" {
            if example, ok := g.config.TypeExamples[format]; ok {
                result[name] = example
                continue
            }
        }
        
        // Generate based on type
        result[name] = g.generateFromType(field.Type)
    }
    
    return result
}

func (g *Generator) generateFromFakerTag(tag string) interface{} {
    switch tag {
    case "email":
        return faker.Email()
    case "name":
        return faker.Name()
    case "phone":
        return faker.Phonenumber()
    case "uuid":
        return faker.UUIDHyphenated()
    case "url":
        return faker.URL()
    case "sentence":
        return faker.Sentence()
    case "paragraph":
        return faker.Paragraph()
    default:
        return faker.Word()
    }
}

func getJSONFieldName(field reflect.StructField) string {
    tag := field.Tag.Get("json")
    if tag == "" {
        return field.Name
    }
    parts := strings.Split(tag, ",")
    return parts[0]
}
```

---

## Phase 3: Auth Playground

### Step 3.1: Auth Configuration

```go
// pkg/auth/schemes.go
package auth

// Config for authentication
type Config struct {
    // Persist credentials in browser localStorage
    PersistCredentials bool `json:"persistCredentials"`
    
    // Security schemes
    Schemes []Scheme `json:"schemes"`
    
    // Default scheme
    DefaultScheme string `json:"defaultScheme,omitempty"`
}

// Scheme represents a security scheme
type Scheme struct {
    Name        string `json:"name"`
    Type        string `json:"type"`        // "bearer", "apiKey", "oauth2", "basic"
    In          string `json:"in"`          // "header", "cookie", "query"
    HeaderName  string `json:"headerName,omitempty"`
    Description string `json:"description,omitempty"`
    
    // OAuth2 specific
    OAuth2 *OAuth2Config `json:"oauth2,omitempty"`
}

type OAuth2Config struct {
    AuthorizationURL string            `json:"authorizationUrl"`
    TokenURL         string            `json:"tokenUrl"`
    Scopes           map[string]string `json:"scopes"`
}

// ToOpenAPI converts to OpenAPI security scheme
func (s Scheme) ToOpenAPI() map[string]interface{} {
    scheme := map[string]interface{}{
        "description": s.Description,
    }
    
    switch s.Type {
    case "bearer":
        scheme["type"] = "http"
        scheme["scheme"] = "bearer"
        scheme["bearerFormat"] = "JWT"
    case "apiKey":
        scheme["type"] = "apiKey"
        scheme["in"] = s.In
        scheme["name"] = s.HeaderName
    case "basic":
        scheme["type"] = "http"
        scheme["scheme"] = "basic"
    case "oauth2":
        scheme["type"] = "oauth2"
        if s.OAuth2 != nil {
            scheme["flows"] = map[string]interface{}{
                "authorizationCode": map[string]interface{}{
                    "authorizationUrl": s.OAuth2.AuthorizationURL,
                    "tokenUrl":         s.OAuth2.TokenURL,
                    "scopes":           s.OAuth2.Scopes,
                },
            }
        }
    }
    
    return scheme
}

// Helpers for creating schemes
func BearerAuth(name string) Scheme {
    return Scheme{
        Name:        name,
        Type:        "bearer",
        Description: "JWT Bearer token authentication",
    }
}

func APIKeyAuth(name, headerName string) Scheme {
    return Scheme{
        Name:        name,
        Type:        "apiKey",
        In:          "header",
        HeaderName:  headerName,
        Description: "API Key authentication",
    }
}

func CookieAuth(name, cookieName string) Scheme {
    return Scheme{
        Name:        name,
        Type:        "apiKey",
        In:          "cookie",
        HeaderName:  cookieName,
        Description: "Cookie-based authentication",
    }
}

func BasicAuth(name string) Scheme {
    return Scheme{
        Name:        name,
        Type:        "basic",
        Description: "HTTP Basic authentication",
    }
}
```

---

## Phase 4: Try-It Console

### Step 4.1: Code Snippet Generator

```go
// pkg/tryit/snippets/generator.go
package snippets

import (
    "encoding/json"
    "fmt"
    "strings"
)

// Request represents an API request for snippet generation
type Request struct {
    Method      string
    URL         string
    Headers     map[string]string
    Body        interface{}
    ContentType string
}

// Generator generates code snippets
type Generator struct{}

func New() *Generator {
    return &Generator{}
}

// Generate creates snippets for all supported languages
func (g *Generator) Generate(req Request) map[string]string {
    return map[string]string{
        "curl":       g.Curl(req),
        "javascript": g.JavaScript(req),
        "go":         g.Go(req),
        "python":     g.Python(req),
        "php":        g.PHP(req),
    }
}

// Curl generates a curl command
func (g *Generator) Curl(req Request) string {
    var parts []string
    parts = append(parts, fmt.Sprintf("curl -X %s \"%s\"", req.Method, req.URL))
    
    for key, value := range req.Headers {
        parts = append(parts, fmt.Sprintf("  -H \"%s: %s\"", key, value))
    }
    
    if req.Body != nil {
        bodyJSON, _ := json.Marshal(req.Body)
        parts = append(parts, fmt.Sprintf("  -d '%s'", string(bodyJSON)))
    }
    
    return strings.Join(parts, " \\\n")
}

// JavaScript generates fetch code
func (g *Generator) JavaScript(req Request) string {
    headers := "{\n"
    for key, value := range req.Headers {
        headers += fmt.Sprintf("    \"%s\": \"%s\",\n", key, value)
    }
    headers += "  }"
    
    bodyStr := ""
    if req.Body != nil {
        bodyJSON, _ := json.MarshalIndent(req.Body, "  ", "  ")
        bodyStr = fmt.Sprintf(",\n  body: JSON.stringify(%s)", string(bodyJSON))
    }
    
    return fmt.Sprintf(`const response = await fetch("%s", {
  method: "%s",
  headers: %s%s
});

const data = await response.json();
console.log(data);`, req.URL, req.Method, headers, bodyStr)
}

// Go generates Go http code
func (g *Generator) Go(req Request) string {
    bodySetup := ""
    bodyArg := "nil"
    
    if req.Body != nil {
        bodyJSON, _ := json.MarshalIndent(req.Body, "", "  ")
        bodySetup = fmt.Sprintf(`body := strings.NewReader(`+"`%s`"+`)
`, string(bodyJSON))
        bodyArg = "body"
    }
    
    headerSetup := ""
    for key, value := range req.Headers {
        headerSetup += fmt.Sprintf(`req.Header.Set("%s", "%s")
`, key, value)
    }
    
    return fmt.Sprintf(`%sreq, err := http.NewRequest("%s", "%s", %s)
if err != nil {
    log.Fatal(err)
}

%s
client := &http.Client{}
resp, err := client.Do(req)
if err != nil {
    log.Fatal(err)
}
defer resp.Body.Close()

var result map[string]interface{}
json.NewDecoder(resp.Body).Decode(&result)
fmt.Println(result)`, bodySetup, req.Method, req.URL, bodyArg, headerSetup)
}

// Python generates Python requests code
func (g *Generator) Python(req Request) string {
    headers := "{\n"
    for key, value := range req.Headers {
        headers += fmt.Sprintf("    \"%s\": \"%s\",\n", key, value)
    }
    headers += "}"
    
    bodyStr := ""
    if req.Body != nil {
        bodyJSON, _ := json.MarshalIndent(req.Body, "", "    ")
        bodyStr = fmt.Sprintf(",\n    json=%s", string(bodyJSON))
    }
    
    return fmt.Sprintf(`import requests

response = requests.%s(
    "%s",
    headers=%s%s
)

print(response.json())`, strings.ToLower(req.Method), req.URL, headers, bodyStr)
}

// PHP generates PHP code
func (g *Generator) PHP(req Request) string {
    headers := ""
    for key, value := range req.Headers {
        headers += fmt.Sprintf("    '%s: %s',\n", key, value)
    }
    
    bodyStr := ""
    if req.Body != nil {
        bodyJSON, _ := json.Marshal(req.Body)
        bodyStr = fmt.Sprintf("\nCURLOPT_POSTFIELDS => '%s',", string(bodyJSON))
    }
    
    return fmt.Sprintf(`<?php
$curl = curl_init();

curl_setopt_array($curl, [
    CURLOPT_URL => "%s",
    CURLOPT_RETURNTRANSFER => true,
    CURLOPT_CUSTOMREQUEST => "%s",
    CURLOPT_HTTPHEADER => [
%s    ],%s
]);

$response = curl_exec($curl);
curl_close($curl);

echo $response;`, req.URL, req.Method, headers, bodyStr)
}
```

### Step 4.2: Environment Configuration

```go
// pkg/tryit/environments.go
package tryit

// Config for Try-It console
type Config struct {
    Enabled        bool          `json:"enabled"`
    Snippets       []string      `json:"snippets"`
    SaveHistory    bool          `json:"saveHistory"`
    MaxHistorySize int           `json:"maxHistorySize"`
    Environments   []Environment `json:"environments"`
}

// Environment represents a set of variables
type Environment struct {
    Name      string            `json:"name"`
    Variables map[string]string `json:"variables"`
    IsDefault bool              `json:"isDefault"`
}

// DefaultConfig returns sensible defaults
func DefaultConfig() Config {
    return Config{
        Enabled:        true,
        Snippets:       []string{"curl", "javascript", "go", "python"},
        SaveHistory:    true,
        MaxHistorySize: 50,
        Environments: []Environment{
            {
                Name:      "Development",
                IsDefault: true,
                Variables: map[string]string{
                    "baseUrl": "http://localhost:8080",
                },
            },
        },
    }
}

// ToScalarConfig converts to Scalar-compatible configuration
func (c Config) ToScalarConfig() map[string]interface{} {
    return map[string]interface{}{
        "hiddenClients":   getHiddenClients(c.Snippets),
        "defaultHttpClient": map[string]interface{}{
            "targetKey": "shell",
            "clientKey": "curl",
        },
    }
}

func getHiddenClients(enabled []string) []string {
    all := []string{"curl", "javascript", "go", "python", "php", "ruby", "java", "csharp"}
    hidden := []string{}
    
    enabledMap := make(map[string]bool)
    for _, s := range enabled {
        enabledMap[s] = true
    }
    
    for _, client := range all {
        if !enabledMap[client] {
            hidden = append(hidden, client)
        }
    }
    
    return hidden
}
```

---

## Phase 5: Version Tracking

### Step 5.1: Spec Differ

```go
// pkg/versioning/diff.go
package versioning

import (
    "encoding/json"
    "fmt"
    "os"
)

// Diff represents differences between two specs
type Diff struct {
    OldVersion string           `json:"oldVersion"`
    NewVersion string           `json:"newVersion"`
    Changes    []Change         `json:"changes"`
    Breaking   []BreakingChange `json:"breaking"`
    Summary    Summary          `json:"summary"`
}

// Change represents a single change
type Change struct {
    Type        ChangeType `json:"type"`
    Path        string     `json:"path"`
    Method      string     `json:"method,omitempty"`
    Description string     `json:"description"`
    IsBreaking  bool       `json:"isBreaking"`
}

type ChangeType string

const (
    ChangeAdded    ChangeType = "added"
    ChangeRemoved  ChangeType = "removed"
    ChangeModified ChangeType = "modified"
)

// BreakingChange represents a breaking change
type BreakingChange struct {
    Path      string `json:"path"`
    Method    string `json:"method"`
    Reason    string `json:"reason"`
    Migration string `json:"migration"`
}

// Summary of changes
type Summary struct {
    AddedEndpoints    int `json:"addedEndpoints"`
    RemovedEndpoints  int `json:"removedEndpoints"`
    ModifiedEndpoints int `json:"modifiedEndpoints"`
    BreakingChanges   int `json:"breakingChanges"`
}

// Differ compares OpenAPI specs
type Differ struct{}

func NewDiffer() *Differ {
    return &Differ{}
}

// Compare compares two spec files
func (d *Differ) Compare(oldPath, newPath string) (*Diff, error) {
    oldSpec, err := loadSpec(oldPath)
    if err != nil {
        return nil, fmt.Errorf("failed to load old spec: %w", err)
    }
    
    newSpec, err := loadSpec(newPath)
    if err != nil {
        return nil, fmt.Errorf("failed to load new spec: %w", err)
    }
    
    return d.CompareSpecs(oldSpec, newSpec)
}

// CompareSpecs compares two parsed specs
func (d *Differ) CompareSpecs(oldSpec, newSpec map[string]interface{}) (*Diff, error) {
    diff := &Diff{
        OldVersion: getVersion(oldSpec),
        NewVersion: getVersion(newSpec),
        Changes:    []Change{},
        Breaking:   []BreakingChange{},
    }
    
    oldPaths := getPaths(oldSpec)
    newPaths := getPaths(newSpec)
    
    // Find added endpoints
    for path, methods := range newPaths {
        if _, exists := oldPaths[path]; !exists {
            for method := range methods {
                diff.Changes = append(diff.Changes, Change{
                    Type:        ChangeAdded,
                    Path:        path,
                    Method:      method,
                    Description: fmt.Sprintf("New endpoint: %s %s", method, path),
                })
                diff.Summary.AddedEndpoints++
            }
        }
    }
    
    // Find removed endpoints (breaking!)
    for path, methods := range oldPaths {
        if _, exists := newPaths[path]; !exists {
            for method := range methods {
                diff.Changes = append(diff.Changes, Change{
                    Type:        ChangeRemoved,
                    Path:        path,
                    Method:      method,
                    Description: fmt.Sprintf("Removed endpoint: %s %s", method, path),
                    IsBreaking:  true,
                })
                diff.Breaking = append(diff.Breaking, BreakingChange{
                    Path:      path,
                    Method:    method,
                    Reason:    "Endpoint removed",
                    Migration: "Update client code to use alternative endpoint or remove usage",
                })
                diff.Summary.RemovedEndpoints++
                diff.Summary.BreakingChanges++
            }
        }
    }
    
    // Find modified endpoints
    for path, oldMethods := range oldPaths {
        if newMethods, exists := newPaths[path]; exists {
            for method, oldOp := range oldMethods {
                if newOp, methodExists := newMethods[method]; methodExists {
                    changes := d.compareOperations(path, method, oldOp, newOp)
                    diff.Changes = append(diff.Changes, changes...)
                    
                    for _, change := range changes {
                        if change.IsBreaking {
                            diff.Summary.BreakingChanges++
                        }
                    }
                    
                    if len(changes) > 0 {
                        diff.Summary.ModifiedEndpoints++
                    }
                }
            }
        }
    }
    
    return diff, nil
}

func (d *Differ) compareOperations(path, method string, oldOp, newOp map[string]interface{}) []Change {
    changes := []Change{}
    
    // Compare request body
    oldBody := getRequestBody(oldOp)
    newBody := getRequestBody(newOp)
    
    if oldBody != nil && newBody == nil {
        changes = append(changes, Change{
            Type:        ChangeModified,
            Path:        path,
            Method:      method,
            Description: "Request body removed",
            IsBreaking:  true,
        })
    }
    
    // Compare required fields
    oldRequired := getRequiredFields(oldOp)
    newRequired := getRequiredFields(newOp)
    
    for _, field := range newRequired {
        if !contains(oldRequired, field) {
            changes = append(changes, Change{
                Type:        ChangeModified,
                Path:        path,
                Method:      method,
                Description: fmt.Sprintf("New required field: %s", field),
                IsBreaking:  true,
            })
        }
    }
    
    // Compare response codes
    oldResponses := getResponseCodes(oldOp)
    newResponses := getResponseCodes(newOp)
    
    for _, code := range oldResponses {
        if !contains(newResponses, code) {
            changes = append(changes, Change{
                Type:        ChangeModified,
                Path:        path,
                Method:      method,
                Description: fmt.Sprintf("Response code %s removed", code),
                IsBreaking:  true,
            })
        }
    }
    
    return changes
}

// Helper functions
func loadSpec(path string) (map[string]interface{}, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, err
    }
    
    var spec map[string]interface{}
    if err := json.Unmarshal(data, &spec); err != nil {
        return nil, err
    }
    
    return spec, nil
}

func getVersion(spec map[string]interface{}) string {
    if info, ok := spec["info"].(map[string]interface{}); ok {
        if version, ok := info["version"].(string); ok {
            return version
        }
    }
    return "unknown"
}

func getPaths(spec map[string]interface{}) map[string]map[string]map[string]interface{} {
    result := make(map[string]map[string]map[string]interface{})
    
    if paths, ok := spec["paths"].(map[string]interface{}); ok {
        for path, methods := range paths {
            result[path] = make(map[string]map[string]interface{})
            if methodMap, ok := methods.(map[string]interface{}); ok {
                for method, op := range methodMap {
                    if opMap, ok := op.(map[string]interface{}); ok {
                        result[path][method] = opMap
                    }
                }
            }
        }
    }
    
    return result
}

func contains(slice []string, item string) bool {
    for _, s := range slice {
        if s == item {
            return true
        }
    }
    return false
}
```

---

## Phase 6: Framework Adapters

### Step 6.1: Chi Adapter

```go
// adapters/chi/adapter.go
package chi

import (
    "github.com/go-chi/chi/v5"
    openapi "github.com/andrianprasetya/open-swag-go"
)

// Mount mounts the documentation on a Chi router
func Mount(r chi.Router, docs *openapi.Docs, basePath string) {
    r.Get(basePath, docs.Handler().ServeHTTP)
    r.Get(basePath+"/openapi.json", docs.SpecHandler().ServeHTTP)
}
```

### Step 6.2: Gin Adapter

```go
// adapters/gin/adapter.go
package gin

import (
    "github.com/gin-gonic/gin"
    openapi "github.com/andrianprasetya/open-swag-go"
)

// Mount mounts the documentation on a Gin router
func Mount(r *gin.Engine, docs *openapi.Docs, basePath string) {
    r.GET(basePath, gin.WrapH(docs.Handler()))
    r.GET(basePath+"/openapi.json", gin.WrapH(docs.SpecHandler()))
}
```

### Step 6.3: Echo Adapter

```go
// adapters/echo/adapter.go
package echo

import (
    "github.com/labstack/echo/v4"
    openapi "github.com/andrianprasetya/open-swag-go"
)

// Mount mounts the documentation on an Echo router
func Mount(e *echo.Echo, docs *openapi.Docs, basePath string) {
    e.GET(basePath, echo.WrapHandler(docs.Handler()))
    e.GET(basePath+"/openapi.json", echo.WrapHandler(docs.SpecHandler()))
}
```

### Step 6.4: Fiber Adapter

```go
// adapters/fiber/adapter.go
package fiber

import (
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/adaptor"
    openapi "github.com/andrianprasetya/open-swag-go"
)

// Mount mounts the documentation on a Fiber app
func Mount(app *fiber.App, docs *openapi.Docs, basePath string) {
    app.Get(basePath, adaptor.HTTPHandler(docs.Handler()))
    app.Get(basePath+"/openapi.json", adaptor.HTTPHandler(docs.SpecHandler()))
}
```

---

## Complete Usage Example

```go
// examples/full-featured/main.go
package main

import (
    "log"
    "net/http"
    
    openapi "github.com/andrianprasetya/open-swag-go"
    "github.com/andrianprasetya/open-swag-go/pkg/auth"
    "github.com/andrianprasetya/open-swag-go/pkg/tryit"
)

// Request/Response types
type CreateUserRequest struct {
    Name  string `json:"name" swagger:"required" example:"John Doe" description:"User's full name"`
    Email string `json:"email" swagger:"required,format=email" example:"john@example.com"`
    Age   int    `json:"age" example:"25" swagger:"min=18,max=120"`
}

type UserResponse struct {
    ID        string `json:"id" format:"uuid"`
    Name      string `json:"name"`
    Email     string `json:"email"`
    CreatedAt string `json:"created_at" format:"date-time"`
}

type ErrorResponse struct {
    Error   string `json:"error"`
    Message string `json:"message"`
}

func main() {
    // Create documentation
    docs := openapi.New(openapi.Config{
        Info: openapi.Info{
            Title:       "My Awesome API",
            Version:     "2.0.0",
            Description: "A production-ready API with comprehensive documentation",
            Contact: openapi.Contact{
                Name:  "API Support",
                Email: "support@example.com",
            },
        },
        Servers: []openapi.Server{
            {URL: "http://localhost:8080", Description: "Development"},
            {URL: "https://api.example.com", Description: "Production"},
        },
        Tags: []openapi.Tag{
            {Name: "Users", Description: "User management endpoints"},
            {Name: "Auth", Description: "Authentication endpoints"},
        },
        
        // UI Configuration
        UI: openapi.UIConfig{
            Theme:              "purple",
            DarkMode:           true,
            Layout:             "modern",
            ShowSidebar:        true,
            SidebarSearch:      true,
            TagGrouping:        true,
            CollapsibleSchemas: true,
        },
        
        // Auth Playground
        Auth: auth.Config{
            PersistCredentials: true,
            Schemes: []auth.Scheme{
                auth.BearerAuth("bearerAuth"),
                auth.APIKeyAuth("apiKey", "X-API-Key"),
            },
        },
        
        // Examples
        Examples: openapi.ExampleConfig{
            AutoGenerate: true,
            UseFaker:     true,
        },
        
        // Try-It Console
        TryIt: tryit.Config{
            Enabled:        true,
            Snippets:       []string{"curl", "javascript", "go", "python"},
            SaveHistory:    true,
            MaxHistorySize: 50,
            Environments: []tryit.Environment{
                {
                    Name:      "Development",
                    IsDefault: true,
                    Variables: map[string]string{
                        "baseUrl": "http://localhost:8080",
                    },
                },
                {
                    Name: "Production",
                    Variables: map[string]string{
                        "baseUrl": "https://api.example.com",
                    },
                },
            },
        },
        
        // Version Tracking
        Versioning: openapi.VersionConfig{
            Enabled:            true,
            ShowDiff:           true,
            HighlightBreaking:  true,
            ShowMigrationGuide: true,
        },
    })
    
    // Register endpoints using fluent API
    docs.Endpoint("/users", "POST").
        Summary("Create a new user").
        Description("Create a new user account with the provided information").
        Tags("Users").
        Body(CreateUserRequest{}).
        Response(201, "User created successfully", UserResponse{}).
        Response(400, "Invalid request", ErrorResponse{}).
        Response(409, "User already exists", ErrorResponse{}).
        Security("bearerAuth").
        Register()
    
    docs.Endpoint("/users/{id}", "GET").
        Summary("Get user by ID").
        Description("Retrieve user details by their unique identifier").
        Tags("Users").
        Param(openapi.PathParam("id").
            Description("User ID").
            Format("uuid").
            Example("550e8400-e29b-41d4-a716-446655440000")).
        Response(200, "User found", UserResponse{}).
        Response(404, "User not found", ErrorResponse{}).
        Security("bearerAuth", "apiKey").
        Register()
    
    docs.Endpoint("/users/{id}", "PUT").
        Summary("Update user").
        Description("Update an existing user's information").
        Tags("Users").
        Param(openapi.PathParam("id").Format("uuid")).
        Body(CreateUserRequest{}).
        Response(200, "User updated", UserResponse{}).
        Response(404, "User not found", ErrorResponse{}).
        Security("bearerAuth").
        Register()
    
    docs.Endpoint("/users/{id}", "DELETE").
        Summary("Delete user").
        Description("Permanently delete a user account").
        Tags("Users").
        Param(openapi.PathParam("id").Format("uuid")).
        Response(204, "User deleted", nil).
        Response(404, "User not found", ErrorResponse{}).
        Security("bearerAuth").
        Register()
    
    // Mount documentation
    mux := http.NewServeMux()
    docs.Mount(mux, "/docs/")
    
    // Your API routes
    mux.HandleFunc("/users", handleUsers)
    mux.HandleFunc("/users/", handleUser)
    
    log.Println("Server starting on :8080")
    log.Println("Documentation available at http://localhost:8080/docs/")
    log.Fatal(http.ListenAndServe(":8080", mux))
}

func handleUsers(w http.ResponseWriter, r *http.Request) {
    // Your handler implementation
}

func handleUser(w http.ResponseWriter, r *http.Request) {
    // Your handler implementation
}
```

---

## Development Roadmap

### MVP (v0.1.0) - 2 weeks
- [ ] Core config structs
- [ ] Endpoint builder (fluent API)
- [ ] Struct → JSON Schema converter
- [ ] Basic Scalar UI integration
- [ ] net/http handler
- [ ] Basic example

### v0.2.0 - 1 week
- [ ] Example generator with faker
- [ ] Struct tag parsing (swagger, example, format)
- [ ] Auth schemes (bearer, apiKey, basic)

### v0.3.0 - 1 week
- [ ] Try-It console configuration
- [ ] Code snippet generator (curl, js, go, python)
- [ ] Environment variables

### v0.4.0 - 1 week
- [ ] Version diff
- [ ] Breaking change detection
- [ ] Migration guide generator

### v0.5.0 - 1 week
- [ ] Framework adapters (Chi, Gin, Echo, Fiber)
- [ ] CLI tool for spec generation
- [ ] Watch mode for development

### v1.0.0
- [ ] Full documentation
- [ ] Comprehensive tests
- [ ] Performance optimization
- [ ] Production ready

---

## Resources

- [OpenAPI 3.1 Specification](https://spec.openapis.org/oas/v3.1.0)
- [JSON Schema](https://json-schema.org/)
- [Scalar Documentation](https://github.com/scalar/scalar)
- [Go Reflection](https://go.dev/blog/laws-of-reflection)

---

## License

MIT License - Feel free to use this guide and code for your project!


---

## Phase 7: Publishing & Versioning

### Step 7.1: Semantic Versioning

Follow [Semantic Versioning](https://semver.org/):
- **MAJOR** (v1.0.0 → v2.0.0): Breaking changes
- **MINOR** (v1.0.0 → v1.1.0): New features, backward compatible
- **PATCH** (v1.0.0 → v1.0.1): Bug fixes, backward compatible

### Step 7.2: Create Git Tags for Releases

```bash
# After completing MVP features
git add .
git commit -m "feat: initial release with core features"

# Create version tag
git tag v0.1.0
git push origin main --tags

# For subsequent releases
git tag v0.2.0
git push origin main --tags
```

### Step 7.3: GitHub Actions for CI/CD

Create `.github/workflows/ci.yml`:
```yaml
name: CI

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.22'
      
      - name: Install dependencies
        run: go mod download
      
      - name: Run tests
        run: go test -v ./...
      
      - name: Run linter
        uses: golangci/golangci-lint-action@v4
        with:
          version: latest

  build:
    runs-on: ubuntu-latest
    needs: test
    steps:
      - uses: actions/checkout@v4
      
      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.22'
      
      - name: Build
        run: go build -v ./...
```

Create `.github/workflows/release.yml`:
```yaml
name: Release

on:
  push:
    tags:
      - 'v*'

jobs:
  release:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0
      
      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.22'
      
      - name: Run tests
        run: go test -v ./...
      
      - name: Create GitHub Release
        uses: softprops/action-gh-release@v1
        with:
          generate_release_notes: true
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

### Step 7.4: Go Package Publishing

Go packages are automatically available via `go get` once pushed to GitHub with tags:

```bash
# Users can install your package with:
go get github.com/andrianprasetya/open-swag-go@v0.1.0

# Or latest:
go get github.com/andrianprasetya/open-swag-go@latest
```

### Step 7.5: Version in Code

Add version constant to your package:

```go
// version.go
package openapi

const (
    Version = "0.1.0"
)

// GetVersion returns the package version
func GetVersion() string {
    return Version
}
```

### Step 7.6: CHANGELOG.md

Create `CHANGELOG.md`:
```markdown
# Changelog

All notable changes to this project will be documented in this file.

## [Unreleased]

## [0.1.0] - 2024-XX-XX

### Added
- Initial release
- Core endpoint definition structs
- Scalar UI integration
- Basic schema generation from Go structs
- net/http handler support

### Features
- `openapi.New()` - Create new docs instance
- `docs.Add()` - Add endpoint definition
- `docs.Mount()` - Mount handlers on mux
- `openapi.Body()` - Request body helper
- `openapi.Response()` - Response helper
- `openapi.PathParam()` / `QueryParam()` - Parameter helpers

## [0.2.0] - TBD

### Added
- Example generator with faker support
- Auth playground configuration
- Framework adapters (Chi, Gin, Echo, Fiber)
```

### Step 7.7: Release Checklist

Before each release:
- [ ] Update `Version` constant in `version.go`
- [ ] Update `CHANGELOG.md`
- [ ] Run all tests: `go test ./...`
- [ ] Run linter: `golangci-lint run`
- [ ] Create git tag: `git tag vX.Y.Z`
- [ ] Push with tags: `git push origin main --tags`

### Step 7.8: pkg.go.dev

Once you push a tagged version, your package will automatically appear on [pkg.go.dev](https://pkg.go.dev):

```
https://pkg.go.dev/github.com/andrianprasetya/open-swag-go
```

Add a badge to your README:
```markdown
[![Go Reference](https://pkg.go.dev/badge/github.com/andrianprasetya/open-swag-go.svg)](https://pkg.go.dev/github.com/andrianprasetya/open-swag-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/andrianprasetya/open-swag-go)](https://goreportcard.com/report/github.com/andrianprasetya/open-swag-go)
```

---

## Real Project Integration

### Step 1: Create Your Package (github.com/andrianprasetya/open-swag-go)

```bash
mkdir go-openapi-ui && cd go-openapi-ui
go mod init github.com/andrianprasetya/open-swag-go
```

### Step 2: Core Package Files

**openapi.go** - Main entry
```go
package openapi

import (
    "encoding/json"
    "net/http"
    "sync"
)

type Docs struct {
    config    Config
    endpoints []Endpoint
    spec      *Spec
    mu        sync.RWMutex
}

func New(config Config) *Docs {
    return &Docs{
        config:    config,
        endpoints: []Endpoint{},
    }
}

// Add registers an endpoint
func (d *Docs) Add(endpoint Endpoint) {
    d.mu.Lock()
    defer d.mu.Unlock()
    d.endpoints = append(d.endpoints, endpoint)
    d.spec = nil // invalidate cache
}

// AddAll registers multiple endpoints
func (d *Docs) AddAll(endpoints ...Endpoint) {
    for _, ep := range endpoints {
        d.Add(ep)
    }
}
```

**endpoint.go** - Endpoint struct
```go
package openapi

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
    Name        string
    In          string // "path", "query", "header"
    Description string
    Required    bool
    Schema      Schema
    Example     interface{}
}

// Responses is a map of status code to response
type Responses map[int]ResponseDef

// ResponseDef represents a response definition
type ResponseDef struct {
    Description string
    Schema      interface{}
    Headers     map[string]Header
}

// RequestBody represents request body
type RequestBody struct {
    Description string
    Required    bool
    Schema      interface{}
    ContentType string
}

// Helper functions
func Body(schema interface{}) *RequestBody {
    return &RequestBody{
        Required:    true,
        Schema:      schema,
        ContentType: "application/json",
    }
}

func BodyWithDesc(description string, schema interface{}) *RequestBody {
    return &RequestBody{
        Description: description,
        Required:    true,
        Schema:      schema,
        ContentType: "application/json",
    }
}

func FormBody(schema interface{}) *RequestBody {
    return &RequestBody{
        Required:    true,
        Schema:      schema,
        ContentType: "multipart/form-data",
    }
}

func Response(description string, schema interface{}) ResponseDef {
    return ResponseDef{
        Description: description,
        Schema:      schema,
    }
}

// Parameter helpers
func PathParam(name, description string) Parameter {
    return Parameter{
        Name:        name,
        In:          "path",
        Description: description,
        Required:    true,
    }
}

func QueryParam(name, description string) Parameter {
    return Parameter{
        Name:        name,
        In:          "query",
        Description: description,
        Required:    false,
    }
}

func HeaderParam(name, description string) Parameter {
    return Parameter{
        Name:        name,
        In:          "header",
        Description: description,
        Required:    false,
    }
}
```

**config.go** - Configuration
```go
package openapi

type Config struct {
    Info     Info
    Servers  []Server
    Tags     []Tag
    UI       UIConfig
    Auth     AuthConfig
    TryIt    TryItConfig
}

type Info struct {
    Title       string
    Version     string
    Description string
    Contact     *Contact
    License     *License
}

type Contact struct {
    Name  string
    Email string
    URL   string
}

type License struct {
    Name string
    URL  string
}

type Server struct {
    URL         string
    Description string
}

type Tag struct {
    Name        string
    Description string
}

type UIConfig struct {
    Theme              string // "purple", "dark", "light"
    DarkMode           bool
    ShowSidebar        bool
    SidebarSearch      bool
    CollapsibleSchemas bool
    CustomCSS          string
}

type AuthConfig struct {
    PersistCredentials bool
    Schemes            []AuthScheme
}

type AuthScheme struct {
    Name        string
    Type        string // "bearer", "apiKey", "basic"
    In          string // "header", "cookie", "query"
    HeaderName  string
    Description string
}

type TryItConfig struct {
    Enabled      bool
    Snippets     []string // "curl", "javascript", "go", "python"
    Environments []Environment
}

type Environment struct {
    Name      string
    Variables map[string]string
    IsDefault bool
}

// Helper functions for auth schemes
func BearerAuth(name string) AuthScheme {
    return AuthScheme{
        Name:        name,
        Type:        "bearer",
        Description: "JWT Bearer token",
    }
}

func APIKeyAuth(name, headerName string) AuthScheme {
    return AuthScheme{
        Name:        name,
        Type:        "apiKey",
        In:          "header",
        HeaderName:  headerName,
        Description: "API Key authentication",
    }
}

func CookieAuth(name, cookieName string) AuthScheme {
    return AuthScheme{
        Name:        name,
        Type:        "apiKey",
        In:          "cookie",
        HeaderName:  cookieName,
        Description: "Cookie authentication",
    }
}
```

**handler.go** - HTTP handlers
```go
package openapi

import (
    _ "embed"
    "encoding/json"
    "net/http"
    "strings"
)

//go:embed templates/scalar.html
var scalarHTML string

// Handler returns the documentation UI
func (d *Docs) Handler() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        config := map[string]interface{}{
            "theme":      d.config.UI.Theme,
            "darkMode":   d.config.UI.DarkMode,
            "showSidebar": d.config.UI.ShowSidebar,
        }
        configJSON, _ := json.Marshal(config)
        
        html := scalarHTML
        html = strings.Replace(html, "{{SPEC_URL}}", "./openapi.json", 1)
        html = strings.Replace(html, "{{CONFIG}}", string(configJSON), 1)
        html = strings.Replace(html, "{{TITLE}}", d.config.Info.Title, 1)
        
        w.Header().Set("Content-Type", "text/html")
        w.Write([]byte(html))
    }
}

// SpecHandler returns the OpenAPI JSON spec
func (d *Docs) SpecHandler() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        spec := d.BuildSpec()
        w.Header().Set("Content-Type", "application/json")
        w.Header().Set("Access-Control-Allow-Origin", "*")
        json.NewEncoder(w).Encode(spec)
    }
}

// Mount registers handlers on a mux
func (d *Docs) Mount(mux *http.ServeMux, basePath string) {
    if !strings.HasSuffix(basePath, "/") {
        basePath += "/"
    }
    mux.HandleFunc(basePath, d.Handler())
    mux.HandleFunc(basePath+"openapi.json", d.SpecHandler())
}

// MountWithRouter for chi/gorilla mux
func (d *Docs) MountWithRouter(mountFn func(pattern string, handler http.HandlerFunc), basePath string) {
    mountFn(basePath, d.Handler())
    mountFn(basePath+"/openapi.json", d.SpecHandler())
}
```

**templates/scalar.html**
```html
<!DOCTYPE html>
<html>
<head>
    <title>{{TITLE}} - API Docs</title>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
</head>
<body>
    <script
        id="api-reference"
        data-url="{{SPEC_URL}}"
        data-configuration='{{CONFIG}}'>
    </script>
    <script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
</body>
</html>
```

---

## Using the Package in Your Real Project

### Project Structure
```
my-go-api/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── docs/
│   │   └── swagger.go      # Documentation setup
│   ├── dto/
│   │   ├── user.go
│   │   └── error.go
│   └── handlers/
│       └── user.go
├── go.mod
└── go.sum
```

### Step 1: Install the package
```bash
go get github.com/andrianprasetya/open-swag-go
```

### Step 2: Define DTOs
```go
// internal/dto/user.go
package dto

type CreateUserRequest struct {
    Name  string `json:"name" validate:"required" example:"John Doe"`
    Email string `json:"email" validate:"required,email" example:"john@example.com"`
    Age   int    `json:"age" validate:"min=18" example:"25"`
}

type UpdateUserRequest struct {
    Name string `json:"name" example:"John Updated"`
    Age  int    `json:"age" example:"26"`
}

type UserResponse struct {
    ID        string `json:"id" format:"uuid"`
    Name      string `json:"name"`
    Email     string `json:"email"`
    Age       int    `json:"age"`
    CreatedAt string `json:"created_at" format:"date-time"`
}

type UserListResponse struct {
    Data  []UserResponse `json:"data"`
    Total int            `json:"total"`
    Page  int            `json:"page"`
}
```

```go
// internal/dto/error.go
package dto

type ErrorResponse struct {
    Code    int    `json:"code" example:"400"`
    Message string `json:"message" example:"Bad request"`
}

type ValidationError struct {
    Code   int               `json:"code" example:"422"`
    Errors map[string]string `json:"errors"`
}
```

### Step 3: Create Handlers with Co-located Swagger Docs
```go
// internal/handlers/user.go
package handlers

import (
    "encoding/json"
    "net/http"
    
    openapi "github.com/andrianprasetya/open-swag-go"
    "my-go-api/internal/dto"
)

// ============================================
// CREATE USER
// ============================================

func CreateUser(w http.ResponseWriter, r *http.Request) {
    // Your implementation
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]string{"id": "123"})
}

// CreateUserDoc - Swagger definition (like @SwaggerInfo decorator)
var CreateUserDoc = openapi.Endpoint{
    Method:      "POST",
    Path:        "/users",
    Handler:     CreateUser,
    Summary:     "Create a new user",
    Description: "Create a new user account with the provided information",
    Tags:        []string{"Users"},
    RequestBody: openapi.Body(dto.CreateUserRequest{}),
    Responses: openapi.Responses{
        201: openapi.Response("User created successfully", dto.UserResponse{}),
        400: openapi.Response("Invalid request body", dto.ErrorResponse{}),
        422: openapi.Response("Validation failed", dto.ValidationError{}),
    },
    Security: []string{"bearerAuth"},
}

// ============================================
// GET USER
// ============================================

func GetUser(w http.ResponseWriter, r *http.Request) {
    // Your implementation
}

var GetUserDoc = openapi.Endpoint{
    Method:      "GET",
    Path:        "/users/{id}",
    Handler:     GetUser,
    Summary:     "Get user by ID",
    Description: "Retrieve a single user by their unique identifier",
    Tags:        []string{"Users"},
    Parameters: []openapi.Parameter{
        openapi.PathParam("id", "User ID (UUID format)"),
    },
    Responses: openapi.Responses{
        200: openapi.Response("User found", dto.UserResponse{}),
        404: openapi.Response("User not found", dto.ErrorResponse{}),
    },
    Security: []string{"bearerAuth"},
}

// ============================================
// UPDATE USER
// ============================================

func UpdateUser(w http.ResponseWriter, r *http.Request) {
    // Your implementation
}

var UpdateUserDoc = openapi.Endpoint{
    Method:      "PUT",
    Path:        "/users/{id}",
    Handler:     UpdateUser,
    Summary:     "Update user",
    Description: "Update an existing user's information",
    Tags:        []string{"Users"},
    Parameters: []openapi.Parameter{
        openapi.PathParam("id", "User ID"),
    },
    RequestBody: openapi.Body(dto.UpdateUserRequest{}),
    Responses: openapi.Responses{
        200: openapi.Response("User updated", dto.UserResponse{}),
        404: openapi.Response("User not found", dto.ErrorResponse{}),
        422: openapi.Response("Validation failed", dto.ValidationError{}),
    },
    Security: []string{"bearerAuth"},
}

// ============================================
// DELETE USER
// ============================================

func DeleteUser(w http.ResponseWriter, r *http.Request) {
    // Your implementation
}

var DeleteUserDoc = openapi.Endpoint{
    Method:      "DELETE",
    Path:        "/users/{id}",
    Handler:     DeleteUser,
    Summary:     "Delete user",
    Description: "Permanently delete a user account",
    Tags:        []string{"Users"},
    Parameters: []openapi.Parameter{
        openapi.PathParam("id", "User ID"),
    },
    Responses: openapi.Responses{
        204: openapi.Response("User deleted", nil),
        404: openapi.Response("User not found", dto.ErrorResponse{}),
    },
    Security: []string{"bearerAuth"},
}

// ============================================
// LIST USERS
// ============================================

func ListUsers(w http.ResponseWriter, r *http.Request) {
    // Your implementation
}

var ListUsersDoc = openapi.Endpoint{
    Method:      "GET",
    Path:        "/users",
    Handler:     ListUsers,
    Summary:     "List all users",
    Description: "Get a paginated list of users",
    Tags:        []string{"Users"},
    Parameters: []openapi.Parameter{
        openapi.QueryParam("page", "Page number"),
        openapi.QueryParam("limit", "Items per page"),
        openapi.QueryParam("search", "Search by name or email"),
    },
    Responses: openapi.Responses{
        200: openapi.Response("Users retrieved", dto.UserListResponse{}),
    },
    Security: []string{"bearerAuth"},
}

// ============================================
// EXPORT ALL ENDPOINTS
// ============================================

// Endpoints exports all user endpoints for swagger registration
var Endpoints = []openapi.Endpoint{
    CreateUserDoc,
    GetUserDoc,
    UpdateUserDoc,
    DeleteUserDoc,
    ListUsersDoc,
}
```

### Step 4: Setup Documentation (Simple - just import from handlers)
```go
// internal/docs/swagger.go
package docs

import (
    openapi "github.com/andrianprasetya/open-swag-go"
    
    "my-go-api/internal/handlers"
)

func Setup() *openapi.Docs {
    docs := openapi.New(openapi.Config{
        Info: openapi.Info{
            Title:       "My Go API",
            Version:     "1.0.0",
            Description: "A production-ready Go API",
            Contact: &openapi.Contact{
                Name:  "API Support",
                Email: "support@example.com",
            },
        },
        Servers: []openapi.Server{
            {URL: "http://localhost:8080", Description: "Development"},
            {URL: "https://api.example.com", Description: "Production"},
        },
        Tags: []openapi.Tag{
            {Name: "Users", Description: "User management endpoints"},
            {Name: "Auth", Description: "Authentication endpoints"},
        },
        UI: openapi.UIConfig{
            Theme:              "purple",
            DarkMode:           true,
            ShowSidebar:        true,
            SidebarSearch:      true,
            CollapsibleSchemas: true,
        },
        Auth: openapi.AuthConfig{
            PersistCredentials: true,
            Schemes: []openapi.AuthScheme{
                openapi.BearerAuth("bearerAuth"),
                openapi.APIKeyAuth("apiKey", "X-API-Key"),
            },
        },
        TryIt: openapi.TryItConfig{
            Enabled:  true,
            Snippets: []string{"curl", "javascript", "go", "python"},
        },
    })
    
    // Simply import endpoints from handlers - they're already defined there!
    docs.AddAll(handlers.Endpoints...)
    
    // Add more handler endpoints as needed
    // docs.AddAll(authHandlers.Endpoints...)
    // docs.AddAll(productHandlers.Endpoints...)
    
    return docs
}
```

### Step 5: Main Entry Point
```go
// cmd/server/main.go
package main

import (
    "log"
    "net/http"
    
    "my-go-api/internal/docs"
    "my-go-api/internal/handlers"
)

func main() {
    mux := http.NewServeMux()
    
    // Setup and mount documentation
    swagger := docs.Setup()
    swagger.Mount(mux, "/docs")
    
    // API routes
    mux.HandleFunc("POST /users", handlers.CreateUser)
    mux.HandleFunc("GET /users", handlers.ListUsers)
    mux.HandleFunc("GET /users/{id}", handlers.GetUser)
    mux.HandleFunc("PUT /users/{id}", handlers.UpdateUser)
    mux.HandleFunc("DELETE /users/{id}", handlers.DeleteUser)
    
    log.Println("🚀 Server running on http://localhost:8080")
    log.Println("📚 Docs available at http://localhost:8080/docs")
    log.Fatal(http.ListenAndServe(":8080", mux))
}
```

---

## Side-by-Side Comparison

| AdonisJS (adonis-open-swagger) | Go (go-openapi-ui) |
|-------------------------------|-------------------|
| Decorator on controller method | `var XxxDoc` next to handler function |
| `@SwaggerInfo({ summary, tags })` | `Summary: "", Tags: []string{}` |
| `@SwaggerParam(name, schema)` | `Parameters: []openapi.Parameter{}` |
| `@SwaggerRequestBody(desc, schema)` | `RequestBody: openapi.Body(struct{})` |
| `@SwaggerResponse(code, desc, schema)` | `Responses: openapi.Responses{code: ...}` |
| `@SwaggerSecurity([{ auth: [] }])` | `Security: []string{"auth"}` |
| `config/swagger.ts` | `internal/docs/swagger.go` |
| Auto-scan routes | `handlers.Endpoints` export |

### Key Pattern: Co-located Definition

**AdonisJS:**
```typescript
// Decorator directly above method
@SwaggerInfo({ summary: "Create user", tags: ["Users"] })
@SwaggerResponse(201, "Created", UserResponse)
public async create({ request }: HttpContext) {
    // handler code
}
```

**Go (same file, next to handler):**
```go
// Handler
func CreateUser(w http.ResponseWriter, r *http.Request) {
    // handler code
}

// Swagger doc right below handler
var CreateUserDoc = openapi.Endpoint{
    Method:   "POST",
    Path:     "/users",
    Handler:  CreateUser,
    Summary:  "Create user",
    Tags:     []string{"Users"},
    Responses: openapi.Responses{
        201: openapi.Response("Created", dto.UserResponse{}),
    },
}
```
