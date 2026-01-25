package snippets

import (
	"strings"
	"testing"
)

func TestCurlSnippet(t *testing.T) {
	gen := New()
	req := Request{
		Method: "POST",
		URL:    "https://api.example.com/users",
		Headers: map[string]string{
			"Content-Type":  "application/json",
			"Authorization": "Bearer token123",
		},
		Body: map[string]interface{}{
			"name":  "John",
			"email": "john@example.com",
		},
	}

	snippet := gen.Curl(req)

	if !strings.Contains(snippet, "curl -X POST") {
		t.Error("expected curl command")
	}
	if !strings.Contains(snippet, "https://api.example.com/users") {
		t.Error("expected URL in snippet")
	}
	if !strings.Contains(snippet, "Authorization") {
		t.Error("expected Authorization header")
	}
}

func TestJavaScriptSnippet(t *testing.T) {
	gen := New()
	req := Request{
		Method: "GET",
		URL:    "https://api.example.com/users",
		Headers: map[string]string{
			"Authorization": "Bearer token123",
		},
	}

	snippet := gen.JavaScript(req)

	if !strings.Contains(snippet, "fetch(") {
		t.Error("expected fetch call")
	}
	if !strings.Contains(snippet, "method: 'GET'") {
		t.Error("expected GET method")
	}
}

func TestGoSnippet(t *testing.T) {
	gen := New()
	req := Request{
		Method: "POST",
		URL:    "https://api.example.com/users",
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: map[string]string{"name": "John"},
	}

	snippet := gen.Go(req)

	if !strings.Contains(snippet, "http.NewRequest") {
		t.Error("expected http.NewRequest")
	}
	if !strings.Contains(snippet, "POST") {
		t.Error("expected POST method")
	}
}

func TestPythonSnippet(t *testing.T) {
	gen := New()
	req := Request{
		Method: "GET",
		URL:    "https://api.example.com/users",
		Headers: map[string]string{
			"Authorization": "Bearer token123",
		},
	}

	snippet := gen.Python(req)

	if !strings.Contains(snippet, "import requests") {
		t.Error("expected requests import")
	}
	if !strings.Contains(snippet, "requests.get") {
		t.Error("expected requests.get call")
	}
}

func TestGenerateAll(t *testing.T) {
	gen := New()
	req := Request{
		Method:  "GET",
		URL:     "https://api.example.com/users",
		Headers: map[string]string{},
	}

	snippets := gen.GenerateAll(req)

	expected := []string{"curl", "javascript", "go", "python", "php"}
	for _, lang := range expected {
		if _, ok := snippets[lang]; !ok {
			t.Errorf("expected %s snippet", lang)
		}
	}
}
