# OpenSwag Integration Guide for AI Assistants

This document helps AI assistants understand how to integrate `open-swag-go` into Go backend projects using Gin, Fiber, or Echo frameworks.

## Package Overview

`open-swag-go` is an OpenAPI documentation package for Go that provides:
- Automatic OpenAPI 3.0 spec generation
- Scalar UI for interactive API documentation
- Framework adapters for Gin, Fiber, Echo, Chi, and net/http
- Authentication protection for docs (Basic Auth & API Key)

## Installation

```bash
go get github.com/andrianprasetya/open-swag-go@latest
```

## Core Concepts

### 1. Import the package
```go
import openswag "github.com/andrianprasetya/open-swag-go"
```

### 2. Import framework adapter
```go
// For Fiber
import fiberadapter "github.com/andrianprasetya/open-swag-go/adapters/fiber"

// For Gin
import ginadapter "github.com/andrianprasetya/open-swag-go/adapters/gin"

// For Echo
import echoadapter "github.com/andrianprasetya/open-swag-go/adapters/echo"
```

### 3. Create docs instance
```go
docs := openswag.New(openswag.Config{
    Info: openswag.Info{
        Title:       "My API",
        Version:     "1.0.0",
        Description: "API description here",
    },
    Servers: []openswag.Server{
        {URL: "http://localhost:8080", Description: "Local"},
        {URL: "https://api.example.com", Description: "Production"},
    },
})
```

### 4. Register endpoints
```go
docs.Add(openswag.Endpoint{
    Method:      "GET",
    Path:        "/api/users",
    Summary:     "Get all users",
    Description: "Returns a list of all users with pagination",
    Tags:        []string{"Users"},
    Response:    []User{},  // Response type for schema generation
})

docs.Add(openswag.Endpoint{
    Method:   "POST",
    Path:     "/api/users",
    Summary:  "Create user",
    Tags:     []string{"Users"},
    Body:     CreateUserRequest{},  // Request body type
    Response: User{},
})
```

### 5. Mount to framework
```go
// Fiber
fiberadapter.Mount(app, docs, "/docs")

// Gin
ginadapter.Mount(router, docs, "/docs")

// Echo
echoadapter.Mount(e, docs, "/docs")
```

---

## Complete Integration Examples

### Fiber Integration

```go
package main

import (
    "log"

    openswag "github.com/andrianprasetya/open-swag-go"
    fiberadapter "github.com/andrianprasetya/open-swag-go/adapters/fiber"
    "github.com/gofiber/fiber/v2"
)

type User struct {
    ID    int    `json:"id" example:"1"`
    Name  string `json:"name" example:"John Doe"`
    Email string `json:"email" example:"john@example.com"`
}

type CreateUserRequest struct {
    Name  string `json:"name" validate:"required" example:"John Doe"`
    Email string `json:"email" validate:"required,email" example:"john@example.com"`
}

type ErrorResponse struct {
    Error   string `json:"error" example:"Something went wrong"`
    Code    int    `json:"code" example:"400"`
}

func main() {
    app := fiber.New()

    // Initialize OpenSwag
    docs := openswag.New(openswag.Config{
        Info: openswag.Info{
            Title:       "My Fiber API",
            Version:     "1.0.0",
            Description: "A sample API built with Fiber",
        },
        Servers: []openswag.Server{
            {URL: "http://localhost:3000", Description: "Local"},
        },
    })

    // Register all endpoints
    docs.AddAll([]openswag.Endpoint{
        {
            Method:      "GET",
            Path:        "/api/users",
            Summary:     "List all users",
            Description: "Returns paginated list of users",
            Tags:        []string{"Users"},
            Response:    []User{},
        },
        {
            Method:   "GET",
            Path:     "/api/users/:id",
            Summary:  "Get user by ID",
            Tags:     []string{"Users"},
            Response: User{},
        },
        {
            Method:   "POST",
            Path:     "/api/users",
            Summary:  "Create new user",
            Tags:     []string{"Users"},
            Body:     CreateUserRequest{},
            Response: User{},
        },
        {
            Method:   "PUT",
            Path:     "/api/users/:id",
            Summary:  "Update user",
            Tags:     []string{"Users"},
            Body:     CreateUserRequest{},
            Response: User{},
        },
        {
            Method:  "DELETE",
            Path:    "/api/users/:id",
            Summary: "Delete user",
            Tags:    []string{"Users"},
        },
    })

    // Mount docs at /docs
    fiberadapter.Mount(app, docs, "/docs")

    // Actual routes
    api := app.Group("/api")
    api.Get("/users", getUsers)
    api.Get("/users/:id", getUserByID)
    api.Post("/users", createUser)
    api.Put("/users/:id", updateUser)
    api.Delete("/users/:id", deleteUser)

    log.Println("Server: http://localhost:3000")
    log.Println("Docs: http://localhost:3000/docs/")
    app.Listen(":3000")
}

func getUsers(c *fiber.Ctx) error {
    return c.JSON([]User{{ID: 1, Name: "John", Email: "john@example.com"}})
}

func getUserByID(c *fiber.Ctx) error {
    return c.JSON(User{ID: 1, Name: "John", Email: "john@example.com"})
}

func createUser(c *fiber.Ctx) error {
    return c.Status(201).JSON(User{ID: 1, Name: "John", Email: "john@example.com"})
}

func updateUser(c *fiber.Ctx) error {
    return c.JSON(User{ID: 1, Name: "John Updated", Email: "john@example.com"})
}

func deleteUser(c *fiber.Ctx) error {
    return c.SendStatus(204)
}
```

### Gin Integration

```go
package main

import (
    "log"

    openswag "github.com/andrianprasetya/open-swag-go"
    ginadapter "github.com/andrianprasetya/open-swag-go/adapters/gin"
    "github.com/gin-gonic/gin"
)

type User struct {
    ID    int    `json:"id" example:"1"`
    Name  string `json:"name" example:"John Doe"`
    Email string `json:"email" example:"john@example.com"`
}

type CreateUserRequest struct {
    Name  string `json:"name" binding:"required" example:"John Doe"`
    Email string `json:"email" binding:"required,email" example:"john@example.com"`
}

func main() {
    r := gin.Default()

    // Initialize OpenSwag
    docs := openswag.New(openswag.Config{
        Info: openswag.Info{
            Title:       "My Gin API",
            Version:     "1.0.0",
            Description: "A sample API built with Gin",
        },
        Servers: []openswag.Server{
            {URL: "http://localhost:8080", Description: "Local"},
        },
    })

    // Register endpoints
    docs.AddAll([]openswag.Endpoint{
        {
            Method:   "GET",
            Path:     "/api/users",
            Summary:  "List all users",
            Tags:     []string{"Users"},
            Response: []User{},
        },
        {
            Method:   "GET",
            Path:     "/api/users/:id",
            Summary:  "Get user by ID",
            Tags:     []string{"Users"},
            Response: User{},
        },
        {
            Method:   "POST",
            Path:     "/api/users",
            Summary:  "Create new user",
            Tags:     []string{"Users"},
            Body:     CreateUserRequest{},
            Response: User{},
        },
    })

    // Mount docs
    ginadapter.Mount(r, docs, "/docs")

    // Routes
    api := r.Group("/api")
    api.GET("/users", getUsers)
    api.GET("/users/:id", getUserByID)
    api.POST("/users", createUser)

    log.Println("Server: http://localhost:8080")
    log.Println("Docs: http://localhost:8080/docs/")
    r.Run(":8080")
}

func getUsers(c *gin.Context) {
    c.JSON(200, []User{{ID: 1, Name: "John", Email: "john@example.com"}})
}

func getUserByID(c *gin.Context) {
    c.JSON(200, User{ID: 1, Name: "John", Email: "john@example.com"})
}

func createUser(c *gin.Context) {
    c.JSON(201, User{ID: 1, Name: "John", Email: "john@example.com"})
}
```

### Echo Integration

```go
package main

import (
    "net/http"

    openswag "github.com/andrianprasetya/open-swag-go"
    echoadapter "github.com/andrianprasetya/open-swag-go/adapters/echo"
    "github.com/labstack/echo/v4"
)

type User struct {
    ID    int    `json:"id" example:"1"`
    Name  string `json:"name" example:"John Doe"`
    Email string `json:"email" example:"john@example.com"`
}

type CreateUserRequest struct {
    Name  string `json:"name" validate:"required" example:"John Doe"`
    Email string `json:"email" validate:"required,email" example:"john@example.com"`
}

func main() {
    e := echo.New()

    // Initialize OpenSwag
    docs := openswag.New(openswag.Config{
        Info: openswag.Info{
            Title:       "My Echo API",
            Version:     "1.0.0",
            Description: "A sample API built with Echo",
        },
        Servers: []openswag.Server{
            {URL: "http://localhost:8080", Description: "Local"},
        },
    })

    // Register endpoints
    docs.AddAll([]openswag.Endpoint{
        {
            Method:   "GET",
            Path:     "/api/users",
            Summary:  "List all users",
            Tags:     []string{"Users"},
            Response: []User{},
        },
        {
            Method:   "POST",
            Path:     "/api/users",
            Summary:  "Create new user",
            Tags:     []string{"Users"},
            Body:     CreateUserRequest{},
            Response: User{},
        },
    })

    // Mount docs
    echoadapter.Mount(e, docs, "/docs")

    // Routes
    api := e.Group("/api")
    api.GET("/users", getUsers)
    api.POST("/users", createUser)

    e.Logger.Fatal(e.Start(":8080"))
}

func getUsers(c echo.Context) error {
    return c.JSON(http.StatusOK, []User{{ID: 1, Name: "John", Email: "john@example.com"}})
}

func createUser(c echo.Context) error {
    return c.JSON(http.StatusCreated, User{ID: 1, Name: "John", Email: "john@example.com"})
}
```

---

## Clean Architecture Integration

For projects using clean architecture (handler/usecase/repository pattern), organize docs registration in handlers:

### Handler with Docs Registration

```go
// internal/handler/user_handler.go
package handler

import (
    "my-api/internal/domain"
    "my-api/internal/usecase"

    openswag "github.com/andrianprasetya/open-swag-go"
    "github.com/gofiber/fiber/v2"
)

type UserHandler struct {
    usecase usecase.UserUsecase
}

func NewUserHandler(uc usecase.UserUsecase) *UserHandler {
    return &UserHandler{usecase: uc}
}

// RegisterDocs registers all user-related endpoints
func (h *UserHandler) RegisterDocs(docs *openswag.Docs) {
    docs.AddAll([]openswag.Endpoint{
        {
            Method:      "GET",
            Path:        "/api/users",
            Summary:     "List users",
            Description: "Get paginated list of users",
            Tags:        []string{"Users"},
            Response:    []domain.User{},
        },
        {
            Method:   "POST",
            Path:     "/api/users",
            Summary:  "Create user",
            Tags:     []string{"Users"},
            Body:     domain.CreateUserRequest{},
            Response: domain.User{},
        },
    })
}

// RegisterRoutes registers all user routes
func (h *UserHandler) RegisterRoutes(router fiber.Router) {
    router.Get("/users", h.GetUsers)
    router.Post("/users", h.CreateUser)
}

func (h *UserHandler) GetUsers(c *fiber.Ctx) error {
    users, err := h.usecase.GetAll()
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }
    return c.JSON(users)
}

func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
    var req domain.CreateUserRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
    }
    user, err := h.usecase.Create(req)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }
    return c.Status(201).JSON(user)
}
```

### Main with Handler Registration

```go
// cmd/main.go
package main

import (
    openswag "github.com/andrianprasetya/open-swag-go"
    fiberadapter "github.com/andrianprasetya/open-swag-go/adapters/fiber"
    "github.com/gofiber/fiber/v2"

    "my-api/internal/handler"
    "my-api/internal/repository"
    "my-api/internal/usecase"
)

func main() {
    app := fiber.New()

    // Initialize dependencies
    userRepo := repository.NewUserRepository()
    userUC := usecase.NewUserUsecase(userRepo)
    userHandler := handler.NewUserHandler(userUC)

    // Initialize docs
    docs := openswag.New(openswag.Config{
        Info: openswag.Info{
            Title:   "My API",
            Version: "1.0.0",
        },
    })

    // Register docs from all handlers
    userHandler.RegisterDocs(docs)
    // productHandler.RegisterDocs(docs)
    // orderHandler.RegisterDocs(docs)

    // Mount docs UI
    fiberadapter.Mount(app, docs, "/docs")

    // Register routes
    api := app.Group("/api")
    userHandler.RegisterRoutes(api)

    app.Listen(":3000")
}
```

---

## Endpoint Configuration Options

### Full Endpoint Structure

```go
openswag.Endpoint{
    // Required
    Method: "GET",  // GET, POST, PUT, PATCH, DELETE
    Path:   "/api/resource/:id",

    // Metadata
    Summary:     "Short description",
    Description: "Detailed description of what this endpoint does",
    Tags:        []string{"TagName"},
    OperationID: "uniqueOperationId",
    Deprecated:  false,

    // Request
    Body: RequestStruct{},  // For POST, PUT, PATCH

    // Response
    Response: ResponseStruct{},  // Main success response type
}
```

### Struct Tags for Schema Generation

```go
type User struct {
    ID        int       `json:"id" example:"1"`
    Name      string    `json:"name" example:"John Doe"`
    Email     string    `json:"email" example:"john@example.com"`
    Age       int       `json:"age,omitempty" example:"25"`
    CreatedAt time.Time `json:"created_at" example:"2024-01-01T00:00:00Z"`
    IsActive  bool      `json:"is_active" example:"true"`
}

type CreateUserRequest struct {
    Name     string `json:"name" validate:"required" example:"John Doe"`
    Email    string `json:"email" validate:"required,email" example:"john@example.com"`
    Password string `json:"password" validate:"required,min=8" example:"secretpass123"`
}
```

---

## Configuration Options

### Basic Configuration

```go
docs := openswag.New(openswag.Config{
    Info: openswag.Info{
        Title:          "API Title",
        Version:        "1.0.0",
        Description:    "API Description",
        TermsOfService: "https://example.com/terms",
        Contact: &openswag.Contact{
            Name:  "API Support",
            Email: "support@example.com",
            URL:   "https://example.com/support",
        },
        License: &openswag.License{
            Name: "MIT",
            URL:  "https://opensource.org/licenses/MIT",
        },
    },
    Servers: []openswag.Server{
        {URL: "http://localhost:8080", Description: "Development"},
        {URL: "https://staging.example.com", Description: "Staging"},
        {URL: "https://api.example.com", Description: "Production"},
    },
})
```

### With Authentication Protection

```go
docs := openswag.New(openswag.Config{
    Info: openswag.Info{
        Title:   "Protected API",
        Version: "1.0.0",
    },
    DocsAuth: &openswag.DocsAuth{
        Enabled: true,
        // Option 1: API Key (recommended)
        APIKey: "your-secret-api-key",
        // Option 2: Basic Auth
        Username: "admin",
        Password: "password",
        Realm:    "API Documentation",
    },
})
```

Access protected docs:
- API Key: `http://localhost:8080/docs/?key=your-secret-api-key`
- Basic Auth: Browser will prompt for username/password

### UI Customization

```go
docs := openswag.New(openswag.Config{
    Info: openswag.Info{
        Title:   "My API",
        Version: "1.0.0",
    },
    UI: openswag.UIConfig{
        Theme:       "purple",  // purple, blue, green, orange, red
        Layout:      "modern",  // modern, classic
        DarkMode:    true,
        ShowSidebar: true,
        CustomCSS:   "body { font-family: 'Inter', sans-serif; }",
    },
})
```

---

## Common Patterns

### Grouping Endpoints by Tags

```go
// User endpoints
docs.AddAll([]openswag.Endpoint{
    {Method: "GET", Path: "/api/users", Tags: []string{"Users"}, ...},
    {Method: "POST", Path: "/api/users", Tags: []string{"Users"}, ...},
})

// Product endpoints
docs.AddAll([]openswag.Endpoint{
    {Method: "GET", Path: "/api/products", Tags: []string{"Products"}, ...},
    {Method: "POST", Path: "/api/products", Tags: []string{"Products"}, ...},
})

// Order endpoints
docs.AddAll([]openswag.Endpoint{
    {Method: "GET", Path: "/api/orders", Tags: []string{"Orders"}, ...},
    {Method: "POST", Path: "/api/orders", Tags: []string{"Orders"}, ...},
})
```

### Nested/Related Resources

```go
docs.AddAll([]openswag.Endpoint{
    // User's orders
    {
        Method:   "GET",
        Path:     "/api/users/:userId/orders",
        Summary:  "Get user's orders",
        Tags:     []string{"Users", "Orders"},
        Response: []Order{},
    },
    // Order items
    {
        Method:   "GET",
        Path:     "/api/orders/:orderId/items",
        Summary:  "Get order items",
        Tags:     []string{"Orders"},
        Response: []OrderItem{},
    },
})
```

---

## Accessing Documentation

After mounting, documentation is available at:
- UI: `http://localhost:PORT/docs/`
- JSON Spec: `http://localhost:PORT/docs/openapi.json`

The JSON spec can be used for:
- Frontend TypeScript generation (openapi-typescript)
- API client generation
- Testing tools (Postman, Insomnia)
- API gateways

---

## TypeScript Generation for Frontend

After your Go server is running, generate TypeScript types:

```bash
# Install openapi-typescript
npm install -D openapi-typescript

# Add to package.json scripts
{
  "scripts": {
    "generate:api": "openapi-typescript http://localhost:8080/docs/openapi.json -o src/types/api.ts"
  }
}

# Run generation
npm run generate:api
```

This creates fully typed API interfaces from your Go structs.
