package openswag

// AuthConfig configures authentication
type AuthConfig struct {
	PersistCredentials bool         `json:"persistCredentials"`
	Schemes            []AuthScheme `json:"schemes"`
	DefaultScheme      string       `json:"defaultScheme,omitempty"`
}

// AuthScheme represents a security scheme
type AuthScheme struct {
	Name        string        `json:"name"`
	Type        string        `json:"type"`
	In          string        `json:"in,omitempty"`
	HeaderName  string        `json:"headerName,omitempty"`
	Description string        `json:"description,omitempty"`
	OAuth2      *OAuth2Config `json:"oauth2,omitempty"`
}

// OAuth2Config for OAuth2 authentication
type OAuth2Config struct {
	AuthorizationURL string            `json:"authorizationUrl"`
	TokenURL         string            `json:"tokenUrl"`
	Scopes           map[string]string `json:"scopes"`
}

// ToOpenAPI converts to OpenAPI security scheme format
func (s AuthScheme) ToOpenAPI() map[string]interface{} {
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

// BearerAuth creates a bearer token auth scheme
func BearerAuth(name string) AuthScheme {
	return AuthScheme{
		Name:        name,
		Type:        "bearer",
		Description: "JWT Bearer token authentication",
	}
}

// APIKeyAuth creates an API key auth scheme
func APIKeyAuth(name, headerName string) AuthScheme {
	return AuthScheme{
		Name:        name,
		Type:        "apiKey",
		In:          "header",
		HeaderName:  headerName,
		Description: "API Key authentication",
	}
}

// CookieAuth creates a cookie-based auth scheme
func CookieAuth(name, cookieName string) AuthScheme {
	return AuthScheme{
		Name:        name,
		Type:        "apiKey",
		In:          "cookie",
		HeaderName:  cookieName,
		Description: "Cookie-based authentication",
	}
}

// BasicAuth creates a basic auth scheme
func BasicAuth(name string) AuthScheme {
	return AuthScheme{
		Name:        name,
		Type:        "basic",
		Description: "HTTP Basic authentication",
	}
}
