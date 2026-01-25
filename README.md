# Open Swag Go

[![Go Reference](https://pkg.go.dev/badge/github.com/andrianprasetya/open-swag-go.svg)](https://pkg.go.dev/github.com/andrianprasetya/open-swag-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/andrianprasetya/open-swag-go)](https://goreportcard.com/report/github.com/andrianprasetya/open-swag-go)

A modern OpenAPI documentation package for Go with decorator-like DX, similar to `adonis-open-swagger`.

## Features

- 🎨 Modern UI (Scalar) with dark mode
- 📝 Struct-based endpoint definitions (co-located with handlers)
- 🔐 Auth playground with credential persistence
- 🧪 Try-it console with code snippets
- 📊 Auto-generated examples from structs
- 🔄 Framework adapters (net/http, Chi, Gin, Echo, Fiber)

## Installation

```bash
go get github.com/andrianprasetya/open-swag-go
```

## Quick Start

```go
package main

import (
    "log"
    "net/http"

    openswag "github.com/andrianprasetya/open-swag-go"
)

// Request/Response types
type CreateUserRequest struct {
    Name  string `json:"name" swagger:"required" example:"John Doe"`
    Email string `json:"email" swagger:"required" example:"john@example.com"`
}

type UserResponse struct {
    ID    string `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

func main() {
    docs := openswag.New(openswag.Config{
        Info: openswag.Info{
            Title:   "My API",
            Version: "1.0.0",
        },
        UI: openswag.UIConfig{
            Theme:    "purple",
            DarkMode: true,
        },
    })

    // Add endpoint
    docs.Add(openswag.Endpoint{
        Method:      "POST",
        Path:        "/users",
        Summary:     "Create user",
        Tags:        []string{"Users"},
        RequestBody: openswag.Body(CreateUserRequest{}),
        Responses: openswag.Responses{
            201: openswag.Response("User created", UserResponse{}),
        },
    })

    mux := http.NewServeMux()
    docs.Mount(mux, "/docs")

    log.Println("Docs at http://localhost:8080/docs/")
    http.ListenAndServe(":8080", mux)
}
```

## Co-located Swagger Definitions

Define swagger docs next to your handlers (like decorators):

```go
// internal/handlers/user.go
package handlers

import (
    "net/http"
    openswag "github.com/andrianprasetya/open-swag-go"
    "myapp/internal/dto"
)

// Handler
func CreateUser(w http.ResponseWriter, r *http.Request) {
    // implementation
}

// Swagger doc (co-located with handler)
var CreateUserDoc = openswag.Endpoint{
    Method:      "POST",
    Path:        "/users",
    Handler:     CreateUser,
    Summary:     "Create a new user",
    Tags:        []string{"Users"},
    RequestBody: openswag.Body(dto.CreateUserRequest{}),
    Responses: openswag.Responses{
        201: openswag.Response("User created", dto.UserResponse{}),
        400: openswag.Response("Bad request", dto.ErrorResponse{}),
    },
    Security: []string{"bearerAuth"},
}

// Export all endpoints
var Endpoints = []openswag.Endpoint{CreateUserDoc}
```

Then in your docs setup:

```go
docs.AddAll(handlers.Endpoints...)
```

## Authentication

```go
docs := openswag.New(openswag.Config{
    // ...
    Auth: openswag.AuthConfig{
        PersistCredentials: true,
        Schemes: []openswag.AuthScheme{
            openswag.BearerAuth("bearerAuth"),
            openswag.APIKeyAuth("apiKey", "X-API-Key"),
        },
    },
})
```

## Parameters

```go
openswag.Endpoint{
    Method: "GET",
    Path:   "/users/{id}",
    Parameters: []openswag.Parameter{
        openswag.PathParam("id", "User ID"),
        openswag.QueryParam("include", "Related resources to include"),
        openswag.HeaderParam("X-Request-ID", "Request tracking ID"),
    },
}
```

## Struct Tags

```go
type CreateUserRequest struct {
    Name  string `json:"name" swagger:"required" example:"John Doe" description:"User's full name"`
    Email string `json:"email" swagger:"required,format=email" example:"john@example.com"`
    Age   int    `json:"age" validate:"required,min=18" example:"25"`
}
```

Supported tags:
- `swagger:"required"` - Mark field as required
- `example:"value"` - Set example value
- `description:"text"` - Set field description
- `format:"email"` - Set format (email, uuid, date-time, etc.)
- `validate:"required"` - Also marks as required (validator library)
- `binding:"required"` - Also marks as required (Gin)

## Framework Adapters

### Chi
```go
import (
    openswag "github.com/andrianprasetya/open-swag-go"
    chiadapter "github.com/andrianprasetya/open-swag-go/adapters/chi"
    "github.com/go-chi/chi/v5"
)

r := chi.NewRouter()
chiadapter.Mount(r, docs, "/docs")
```

### Gin
```go
import (
    openswag "github.com/andrianprasetya/open-swag-go"
    ginadapter "github.com/andrianprasetya/open-swag-go/adapters/gin"
    "github.com/gin-gonic/gin"
)

r := gin.Default()
ginadapter.Mount(r, docs, "/docs")
```

### Echo
```go
import (
    openswag "github.com/andrianprasetya/open-swag-go"
    echoadapter "github.com/andrianprasetya/open-swag-go/adapters/echo"
    "github.com/labstack/echo/v4"
)

e := echo.New()
echoadapter.Mount(e, docs, "/docs")
```

### Fiber
```go
import (
    openswag "github.com/andrianprasetya/open-swag-go"
    fiberadapter "github.com/andrianprasetya/open-swag-go/adapters/fiber"
    "github.com/gofiber/fiber/v2"
)

app := fiber.New()
fiberadapter.Mount(app, docs, "/docs")
```

## License

MIT License
