package main

import (
	"encoding/json"
	"log"
	"net/http"

	openswag "github.com/andrianprasetya/open-swag-go"
)

// ============================================
// DTOs
// ============================================

type CreateUserRequest struct {
	Name     string `json:"name" swagger:"required" example:"John Doe" description:"User's full name"`
	Email    string `json:"email" swagger:"required,format=email" example:"john@example.com"`
	Password string `json:"password" swagger:"required" example:"********" description:"Min 8 characters"`
	Age      int    `json:"age" example:"25" description:"Must be 18 or older"`
	Role     string `json:"role" example:"user" description:"user, admin, or moderator"`
}

type UpdateUserRequest struct {
	Name string `json:"name" example:"John Updated"`
	Age  int    `json:"age" example:"26"`
	Role string `json:"role" example:"admin"`
}

type UserResponse struct {
	ID        string `json:"id" format:"uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name      string `json:"name" example:"John Doe"`
	Email     string `json:"email" example:"john@example.com"`
	Age       int    `json:"age" example:"25"`
	Role      string `json:"role" example:"user"`
	CreatedAt string `json:"created_at" format:"date-time" example:"2024-01-15T10:30:00Z"`
	UpdatedAt string `json:"updated_at" format:"date-time" example:"2024-01-15T10:30:00Z"`
}

type UserListResponse struct {
	Data       []UserResponse `json:"data"`
	Total      int            `json:"total" example:"100"`
	Page       int            `json:"page" example:"1"`
	PerPage    int            `json:"per_page" example:"10"`
	TotalPages int            `json:"total_pages" example:"10"`
}

type LoginRequest struct {
	Email    string `json:"email" swagger:"required" example:"john@example.com"`
	Password string `json:"password" swagger:"required" example:"********"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIs..."`
	RefreshToken string `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIs..."`
	ExpiresIn    int    `json:"expires_in" example:"3600"`
	TokenType    string `json:"token_type" example:"Bearer"`
}

type ErrorResponse struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"Bad request"`
	Details string `json:"details,omitempty" example:"Invalid email format"`
}

type ValidationError struct {
	Code   int               `json:"code" example:"422"`
	Errors map[string]string `json:"errors"`
}

// ============================================
// HANDLERS
// ============================================

func createUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(UserResponse{
		ID:    "550e8400-e29b-41d4-a716-446655440000",
		Name:  "John Doe",
		Email: "john@example.com",
		Age:   25,
		Role:  "user",
	})
}

func getUser(w http.ResponseWriter, r *http.Request)    { json.NewEncoder(w).Encode(UserResponse{}) }
func updateUser(w http.ResponseWriter, r *http.Request) { json.NewEncoder(w).Encode(UserResponse{}) }
func deleteUser(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusNoContent) }
func listUsers(w http.ResponseWriter, r *http.Request)  { json.NewEncoder(w).Encode(UserListResponse{}) }
func login(w http.ResponseWriter, r *http.Request)      { json.NewEncoder(w).Encode(LoginResponse{}) }
func healthCheck(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// ============================================
// SWAGGER DOCS (Co-located with handlers)
// ============================================

var CreateUserDoc = openswag.Endpoint{
	Method:      "POST",
	Path:        "/api/v1/users",
	Handler:     createUser,
	Summary:     "Create a new user",
	Description: "Create a new user account with the provided information. Requires admin privileges.",
	Tags:        []string{"Users"},
	RequestBody: openswag.Body(CreateUserRequest{}),
	Responses: openswag.Responses{
		201: openswag.Response("User created successfully", UserResponse{}),
		400: openswag.Response("Invalid request body", ErrorResponse{}),
		409: openswag.Response("User already exists", ErrorResponse{}),
		422: openswag.Response("Validation failed", ValidationError{}),
	},
	Security: []string{"bearerAuth"},
}

var GetUserDoc = openswag.Endpoint{
	Method:      "GET",
	Path:        "/api/v1/users/{id}",
	Handler:     getUser,
	Summary:     "Get user by ID",
	Description: "Retrieve a single user by their unique identifier",
	Tags:        []string{"Users"},
	Parameters: []openswag.Parameter{
		openswag.PathParam("id", "User ID (UUID format)"),
	},
	Responses: openswag.Responses{
		200: openswag.Response("User found", UserResponse{}),
		404: openswag.Response("User not found", ErrorResponse{}),
	},
	Security: []string{"bearerAuth"},
}

var UpdateUserDoc = openswag.Endpoint{
	Method:      "PUT",
	Path:        "/api/v1/users/{id}",
	Handler:     updateUser,
	Summary:     "Update user",
	Description: "Update an existing user's information",
	Tags:        []string{"Users"},
	Parameters: []openswag.Parameter{
		openswag.PathParam("id", "User ID"),
	},
	RequestBody: openswag.Body(UpdateUserRequest{}),
	Responses: openswag.Responses{
		200: openswag.Response("User updated", UserResponse{}),
		404: openswag.Response("User not found", ErrorResponse{}),
		422: openswag.Response("Validation failed", ValidationError{}),
	},
	Security: []string{"bearerAuth"},
}

var DeleteUserDoc = openswag.Endpoint{
	Method:      "DELETE",
	Path:        "/api/v1/users/{id}",
	Handler:     deleteUser,
	Summary:     "Delete user",
	Description: "Permanently delete a user account. This action cannot be undone.",
	Tags:        []string{"Users"},
	Parameters: []openswag.Parameter{
		openswag.PathParam("id", "User ID"),
	},
	Responses: openswag.Responses{
		204: openswag.Response("User deleted", nil),
		404: openswag.Response("User not found", ErrorResponse{}),
	},
	Security: []string{"bearerAuth"},
}

var ListUsersDoc = openswag.Endpoint{
	Method:      "GET",
	Path:        "/api/v1/users",
	Handler:     listUsers,
	Summary:     "List all users",
	Description: "Get a paginated list of users with optional filtering",
	Tags:        []string{"Users"},
	Parameters: []openswag.Parameter{
		openswag.QueryParam("page", "Page number (default: 1)"),
		openswag.QueryParam("per_page", "Items per page (default: 10, max: 100)"),
		openswag.QueryParam("search", "Search by name or email"),
		openswag.QueryParam("role", "Filter by role (user, admin, moderator)"),
		openswag.QueryParam("sort", "Sort field (name, email, created_at)"),
		openswag.QueryParam("order", "Sort order (asc, desc)"),
	},
	Responses: openswag.Responses{
		200: openswag.Response("Users retrieved", UserListResponse{}),
	},
	Security: []string{"bearerAuth"},
}

var LoginDoc = openswag.Endpoint{
	Method:      "POST",
	Path:        "/api/v1/auth/login",
	Handler:     login,
	Summary:     "User login",
	Description: "Authenticate user and receive access tokens",
	Tags:        []string{"Authentication"},
	RequestBody: openswag.Body(LoginRequest{}),
	Responses: openswag.Responses{
		200: openswag.Response("Login successful", LoginResponse{}),
		401: openswag.Response("Invalid credentials", ErrorResponse{}),
		422: openswag.Response("Validation failed", ValidationError{}),
	},
}

var HealthCheckDoc = openswag.Endpoint{
	Method:      "GET",
	Path:        "/health",
	Handler:     healthCheck,
	Summary:     "Health check",
	Description: "Check if the API is running",
	Tags:        []string{"System"},
	Responses: openswag.Responses{
		200: openswag.Response("API is healthy", nil),
	},
}

// Export all endpoints
var UserEndpoints = []openswag.Endpoint{
	CreateUserDoc,
	GetUserDoc,
	UpdateUserDoc,
	DeleteUserDoc,
	ListUsersDoc,
}

var AuthEndpoints = []openswag.Endpoint{
	LoginDoc,
}

var SystemEndpoints = []openswag.Endpoint{
	HealthCheckDoc,
}

// ============================================
// MAIN
// ============================================

func main() {
	// Create documentation
	docs := openswag.New(openswag.Config{
		Info: openswag.Info{
			Title:       "My Awesome API",
			Version:     "2.0.0",
			Description: "A production-ready REST API with comprehensive documentation.\n\n## Features\n- User management\n- JWT authentication\n- Role-based access control",
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
			{URL: "https://staging-api.example.com", Description: "Staging"},
			{URL: "https://api.example.com", Description: "Production"},
		},
		Tags: []openswag.Tag{
			{Name: "Users", Description: "User management endpoints"},
			{Name: "Authentication", Description: "Authentication and authorization"},
			{Name: "System", Description: "System and health endpoints"},
		},
		UI: openswag.UIConfig{
			Theme:              "purple",
			DarkMode:           true,
			Layout:             "modern",
			ShowSidebar:        true,
			SidebarSearch:      true,
			TagGrouping:        true,
			CollapsibleSchemas: true,
		},
		Auth: openswag.AuthConfig{
			PersistCredentials: true,
			Schemes: []openswag.AuthScheme{
				openswag.BearerAuth("bearerAuth"),
				openswag.APIKeyAuth("apiKey", "X-API-Key"),
			},
		},
		TryIt: openswag.TryItConfig{
			Enabled:  true,
			Snippets: []string{"curl", "javascript", "go", "python"},
			Environments: []openswag.Environment{
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
	})

	// Register all endpoints
	docs.AddAll(UserEndpoints...)
	docs.AddAll(AuthEndpoints...)
	docs.AddAll(SystemEndpoints...)

	// Setup routes
	mux := http.NewServeMux()

	// Mount documentation
	docs.Mount(mux, "/docs")

	// API routes
	mux.HandleFunc("POST /api/v1/users", createUser)
	mux.HandleFunc("GET /api/v1/users", listUsers)
	mux.HandleFunc("GET /api/v1/users/{id}", getUser)
	mux.HandleFunc("PUT /api/v1/users/{id}", updateUser)
	mux.HandleFunc("DELETE /api/v1/users/{id}", deleteUser)
	mux.HandleFunc("POST /api/v1/auth/login", login)
	mux.HandleFunc("GET /health", healthCheck)

	log.Println("🚀 Server running on http://localhost:8080")
	log.Println("📚 API Docs: http://localhost:8080/docs/")
	log.Println("📋 OpenAPI Spec: http://localhost:8080/docs/openapi.json")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
