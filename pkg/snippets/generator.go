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

// Generator generates code snippets for various languages
type Generator struct{}

// New creates a new snippet generator
func New() *Generator {
	return &Generator{}
}

// GenerateAll creates snippets for all supported languages
func (g *Generator) GenerateAll(req Request) map[string]string {
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
	parts = append(parts, fmt.Sprintf("curl -X %s '%s'", req.Method, req.URL))

	for key, value := range req.Headers {
		parts = append(parts, fmt.Sprintf("  -H '%s: %s'", key, value))
	}

	if req.Body != nil {
		bodyJSON, _ := json.Marshal(req.Body)
		parts = append(parts, fmt.Sprintf("  -d '%s'", string(bodyJSON)))
	}

	return strings.Join(parts, " \\\n")
}

// JavaScript generates fetch code
func (g *Generator) JavaScript(req Request) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("const response = await fetch('%s', {\n", req.URL))
	sb.WriteString(fmt.Sprintf("  method: '%s',\n", req.Method))

	// Headers
	sb.WriteString("  headers: {\n")
	for key, value := range req.Headers {
		sb.WriteString(fmt.Sprintf("    '%s': '%s',\n", key, value))
	}
	sb.WriteString("  },\n")

	// Body
	if req.Body != nil {
		bodyJSON, _ := json.MarshalIndent(req.Body, "  ", "  ")
		sb.WriteString(fmt.Sprintf("  body: JSON.stringify(%s),\n", string(bodyJSON)))
	}

	sb.WriteString("});\n\n")
	sb.WriteString("const data = await response.json();\n")
	sb.WriteString("console.log(data);")

	return sb.String()
}

// Go generates Go http code
func (g *Generator) Go(req Request) string {
	var sb strings.Builder

	sb.WriteString("package main\n\n")
	sb.WriteString("import (\n")
	sb.WriteString("    \"encoding/json\"\n")
	sb.WriteString("    \"fmt\"\n")
	sb.WriteString("    \"log\"\n")
	sb.WriteString("    \"net/http\"\n")
	if req.Body != nil {
		sb.WriteString("    \"strings\"\n")
	}
	sb.WriteString(")\n\n")
	sb.WriteString("func main() {\n")

	// Body setup
	bodyArg := "nil"
	if req.Body != nil {
		bodyJSON, _ := json.MarshalIndent(req.Body, "    ", "  ")
		sb.WriteString(fmt.Sprintf("    body := strings.NewReader(`%s`)\n\n", string(bodyJSON)))
		bodyArg = "body"
	}

	sb.WriteString(fmt.Sprintf("    req, err := http.NewRequest(\"%s\", \"%s\", %s)\n", req.Method, req.URL, bodyArg))
	sb.WriteString("    if err != nil {\n")
	sb.WriteString("        log.Fatal(err)\n")
	sb.WriteString("    }\n\n")

	// Headers
	for key, value := range req.Headers {
		sb.WriteString(fmt.Sprintf("    req.Header.Set(\"%s\", \"%s\")\n", key, value))
	}

	sb.WriteString("\n    client := &http.Client{}\n")
	sb.WriteString("    resp, err := client.Do(req)\n")
	sb.WriteString("    if err != nil {\n")
	sb.WriteString("        log.Fatal(err)\n")
	sb.WriteString("    }\n")
	sb.WriteString("    defer resp.Body.Close()\n\n")
	sb.WriteString("    var result map[string]interface{}\n")
	sb.WriteString("    json.NewDecoder(resp.Body).Decode(&result)\n")
	sb.WriteString("    fmt.Println(result)\n")
	sb.WriteString("}")

	return sb.String()
}

// Python generates Python requests code
func (g *Generator) Python(req Request) string {
	var sb strings.Builder

	sb.WriteString("import requests\n\n")

	// Headers
	sb.WriteString("headers = {\n")
	for key, value := range req.Headers {
		sb.WriteString(fmt.Sprintf("    '%s': '%s',\n", key, value))
	}
	sb.WriteString("}\n\n")

	// Body
	if req.Body != nil {
		bodyJSON, _ := json.MarshalIndent(req.Body, "", "    ")
		sb.WriteString(fmt.Sprintf("data = %s\n\n", string(bodyJSON)))
	}

	// Request
	method := strings.ToLower(req.Method)
	if req.Body != nil {
		sb.WriteString(fmt.Sprintf("response = requests.%s(\n", method))
		sb.WriteString(fmt.Sprintf("    '%s',\n", req.URL))
		sb.WriteString("    headers=headers,\n")
		sb.WriteString("    json=data,\n")
		sb.WriteString(")\n\n")
	} else {
		sb.WriteString(fmt.Sprintf("response = requests.%s(\n", method))
		sb.WriteString(fmt.Sprintf("    '%s',\n", req.URL))
		sb.WriteString("    headers=headers,\n")
		sb.WriteString(")\n\n")
	}

	sb.WriteString("print(response.status_code)\n")
	sb.WriteString("print(response.json())")

	return sb.String()
}

// PHP generates PHP curl code
func (g *Generator) PHP(req Request) string {
	var sb strings.Builder

	sb.WriteString("<?php\n\n")
	sb.WriteString("$curl = curl_init();\n\n")
	sb.WriteString("curl_setopt_array($curl, [\n")
	sb.WriteString(fmt.Sprintf("    CURLOPT_URL => '%s',\n", req.URL))
	sb.WriteString("    CURLOPT_RETURNTRANSFER => true,\n")
	sb.WriteString(fmt.Sprintf("    CURLOPT_CUSTOMREQUEST => '%s',\n", req.Method))

	// Headers
	sb.WriteString("    CURLOPT_HTTPHEADER => [\n")
	for key, value := range req.Headers {
		sb.WriteString(fmt.Sprintf("        '%s: %s',\n", key, value))
	}
	sb.WriteString("    ],\n")

	// Body
	if req.Body != nil {
		bodyJSON, _ := json.Marshal(req.Body)
		sb.WriteString(fmt.Sprintf("    CURLOPT_POSTFIELDS => '%s',\n", string(bodyJSON)))
	}

	sb.WriteString("]);\n\n")
	sb.WriteString("$response = curl_exec($curl);\n")
	sb.WriteString("$httpCode = curl_getinfo($curl, CURLINFO_HTTP_CODE);\n")
	sb.WriteString("curl_close($curl);\n\n")
	sb.WriteString("echo \"Status: $httpCode\\n\";\n")
	sb.WriteString("echo $response;")

	return sb.String()
}

// Ruby generates Ruby code
func (g *Generator) Ruby(req Request) string {
	var sb strings.Builder

	sb.WriteString("require 'net/http'\n")
	sb.WriteString("require 'json'\n")
	sb.WriteString("require 'uri'\n\n")

	sb.WriteString(fmt.Sprintf("uri = URI('%s')\n", req.URL))
	sb.WriteString("http = Net::HTTP.new(uri.host, uri.port)\n")
	sb.WriteString("http.use_ssl = uri.scheme == 'https'\n\n")

	sb.WriteString(fmt.Sprintf("request = Net::HTTP::%s.new(uri)\n", strings.Title(strings.ToLower(req.Method))))

	for key, value := range req.Headers {
		sb.WriteString(fmt.Sprintf("request['%s'] = '%s'\n", key, value))
	}

	if req.Body != nil {
		bodyJSON, _ := json.Marshal(req.Body)
		sb.WriteString(fmt.Sprintf("request.body = '%s'\n", string(bodyJSON)))
	}

	sb.WriteString("\nresponse = http.request(request)\n")
	sb.WriteString("puts response.code\n")
	sb.WriteString("puts response.body")

	return sb.String()
}
