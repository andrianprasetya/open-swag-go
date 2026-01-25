package openswag

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNew(t *testing.T) {
	docs := New(Config{
		Info: Info{
			Title:   "Test API",
			Version: "1.0.0",
		},
	})

	if docs == nil {
		t.Fatal("expected docs to be created")
	}

	if docs.config.UI.Theme != "purple" {
		t.Errorf("expected default theme 'purple', got '%s'", docs.config.UI.Theme)
	}
}

func TestAddEndpoint(t *testing.T) {
	docs := New(Config{
		Info: Info{Title: "Test API", Version: "1.0.0"},
	})

	docs.Add(Endpoint{
		Method:  "GET",
		Path:    "/users",
		Summary: "List users",
		Tags:    []string{"Users"},
	})

	if len(docs.endpoints) != 1 {
		t.Errorf("expected 1 endpoint, got %d", len(docs.endpoints))
	}
}

func TestAddAllEndpoints(t *testing.T) {
	docs := New(Config{
		Info: Info{Title: "Test API", Version: "1.0.0"},
	})

	endpoints := []Endpoint{
		{Method: "GET", Path: "/users", Summary: "List users"},
		{Method: "POST", Path: "/users", Summary: "Create user"},
	}

	docs.AddAll(endpoints...)

	if len(docs.endpoints) != 2 {
		t.Errorf("expected 2 endpoints, got %d", len(docs.endpoints))
	}
}

func TestBuildSpec(t *testing.T) {
	docs := New(Config{
		Info: Info{Title: "Test API", Version: "1.0.0"},
		Tags: []Tag{{Name: "Users", Description: "User endpoints"}},
	})

	docs.Add(Endpoint{
		Method:  "GET",
		Path:    "/users",
		Summary: "List users",
		Tags:    []string{"Users"},
		Responses: Responses{
			200: Response("Success", nil),
		},
	})

	spec := docs.BuildSpec()

	if spec.OpenAPI != "3.1.0" {
		t.Errorf("expected OpenAPI 3.1.0, got %s", spec.OpenAPI)
	}

	if spec.Info.Title != "Test API" {
		t.Errorf("expected title 'Test API', got '%s'", spec.Info.Title)
	}

	if _, exists := spec.Paths["/users"]; !exists {
		t.Error("expected /users path to exist")
	}
}

func TestSpecHandler(t *testing.T) {
	docs := New(Config{
		Info: Info{Title: "Test API", Version: "1.0.0"},
	})

	docs.Add(Endpoint{
		Method:  "GET",
		Path:    "/health",
		Summary: "Health check",
	})

	req := httptest.NewRequest(http.MethodGet, "/openapi.json", nil)
	rec := httptest.NewRecorder()

	docs.SpecHandler()(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}

	contentType := rec.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("expected Content-Type 'application/json', got '%s'", contentType)
	}

	var spec map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &spec); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}
}

func TestSchemaFromType(t *testing.T) {
	type TestStruct struct {
		Name  string `json:"name" example:"John"`
		Age   int    `json:"age"`
		Email string `json:"email" swagger:"required"`
	}

	schema := SchemaFromType(TestStruct{})

	if schema.Type != "object" {
		t.Errorf("expected type 'object', got '%s'", schema.Type)
	}

	if _, exists := schema.Properties["name"]; !exists {
		t.Error("expected 'name' property to exist")
	}

	if _, exists := schema.Properties["age"]; !exists {
		t.Error("expected 'age' property to exist")
	}
}

func TestAuthSchemes(t *testing.T) {
	bearer := BearerAuth("bearerAuth")
	if bearer.Type != "bearer" {
		t.Errorf("expected type 'bearer', got '%s'", bearer.Type)
	}

	apiKey := APIKeyAuth("apiKey", "X-API-Key")
	if apiKey.HeaderName != "X-API-Key" {
		t.Errorf("expected header 'X-API-Key', got '%s'", apiKey.HeaderName)
	}

	basic := BasicAuth("basicAuth")
	if basic.Type != "basic" {
		t.Errorf("expected type 'basic', got '%s'", basic.Type)
	}
}

func TestParameterHelpers(t *testing.T) {
	path := PathParam("id", "User ID")
	if path.In != "path" || !path.Required {
		t.Error("PathParam should be in 'path' and required")
	}

	query := QueryParam("page", "Page number")
	if query.In != "query" || query.Required {
		t.Error("QueryParam should be in 'query' and not required")
	}

	header := HeaderParam("X-Request-ID", "Request ID")
	if header.In != "header" {
		t.Error("HeaderParam should be in 'header'")
	}
}

func TestBodyHelper(t *testing.T) {
	type Request struct {
		Name string `json:"name"`
	}

	body := Body(Request{})
	if body.ContentType != "application/json" {
		t.Errorf("expected content type 'application/json', got '%s'", body.ContentType)
	}

	if !body.Required {
		t.Error("Body should be required by default")
	}
}

func TestResponseHelper(t *testing.T) {
	resp := Response("Success", nil)
	if resp.Description != "Success" {
		t.Errorf("expected description 'Success', got '%s'", resp.Description)
	}
}
