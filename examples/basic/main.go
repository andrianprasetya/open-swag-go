package main

import (
	"encoding/json"
	"log"
	"net/http"

	openswag "github.com/andrianprasetya/open-swag-go"
)

// DTO types
type CreateUserRequest struct {
	Name  string `json:"name" swagger:"required" example:"John Doe" description:"User's full name"`
	Email string `json:"email" swagger:"required" example:"john@example.com"`
	Age   int    `json:"age" example:"25"`
}

type UserResponse struct {
	ID        string `json:"id" format:"uuid"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Age       int    `json:"age"`
	CreatedAt string `json:"created_at" format:"date-time"`
}

type ErrorResponse struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"Bad request"`
}

// Handlers
func createUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(UserResponse{
		ID:    "123",
		Name:  "John Doe",
		Email: "john@example.com",
	})
}

func getUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(UserResponse{
		ID:    "123",
		Name:  "John Doe",
		Email: "john@example.com",
	})
}

func listUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode([]UserResponse{})
}

// Swagger docs (co-located style)
var CreateUserDoc = openswag.Endpoint{
	Method:      "POST",
	Path:        "/users",
	Handler:     createUser,
	Summary:     "Create a new user",
	Description: "Create a new user account with the provided information",
	Tags:        []string{"Users"},
	RequestBody: openswag.Body(CreateUserRequest{}),
	Responses: openswag.Responses{
		201: openswag.Response("User created successfully", UserResponse{}),
		400: openswag.Response("Invalid request", ErrorResponse{}),
	},
	Security: []string{"bearerAuth"},
}

var GetUserDoc = openswag.Endpoint{
	Method:      "GET",
	Path:        "/users/{id}",
	Handler:     getUser,
	Summary:     "Get user by ID",
	Description: "Retrieve a user by their unique identifier",
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

var ListUsersDoc = openswag.Endpoint{
	Method:      "GET",
	Path:        "/users",
	Handler:     listUsers,
	Summary:     "List all users",
	Description: "Get a paginated list of users",
	Tags:        []string{"Users"},
	Parameters: []openswag.Parameter{
		openswag.QueryParam("page", "Page number"),
		openswag.QueryParam("limit", "Items per page"),
	},
	Responses: openswag.Responses{
		200: openswag.Response("Users retrieved", []UserResponse{}),
	},
	Security: []string{"bearerAuth"},
}

func main() {
	// Create documentation
	docs := openswag.New(openswag.Config{
		Info: openswag.Info{
			Title:       "My API",
			Version:     "1.0.0",
			Description: "A sample API demonstrating open-swag-go",
			Contact: &openswag.Contact{
				Name:  "API Support",
				Email: "support@example.com",
			},
		},
		Servers: []openswag.Server{
			{URL: "http://localhost:8080", Description: "Development"},
		},
		Tags: []openswag.Tag{
			{Name: "Users", Description: "User management endpoints"},
		},
		UI: openswag.UIConfig{
			Theme:       "purple",
			DarkMode:    true,
			ShowSidebar: true,
		},
		Auth: openswag.AuthConfig{
			PersistCredentials: true,
			Schemes: []openswag.AuthScheme{
				openswag.BearerAuth("bearerAuth"),
				openswag.APIKeyAuth("apiKey", "X-API-Key"),
			},
		},
	})

	// Register endpoints
	docs.AddAll(CreateUserDoc, GetUserDoc, ListUsersDoc)

	// Setup routes
	mux := http.NewServeMux()
	docs.Mount(mux, "/docs")

	// API routes
	mux.HandleFunc("POST /users", createUser)
	mux.HandleFunc("GET /users", listUsers)
	mux.HandleFunc("GET /users/{id}", getUser)

	log.Println("🚀 Server running on http://localhost:8080")
	log.Println("📚 Docs available at http://localhost:8080/docs/")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
