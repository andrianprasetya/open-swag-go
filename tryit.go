package openswag

// TryItConfig configures the Try-It console
type TryItConfig struct {
	Enabled        bool          `json:"enabled"`
	Snippets       []string      `json:"snippets"`
	SaveHistory    bool          `json:"saveHistory"`
	MaxHistorySize int           `json:"maxHistorySize"`
	Environments   []Environment `json:"environments"`
}

// Environment represents a set of variables for Try-It
type Environment struct {
	Name      string            `json:"name"`
	Variables map[string]string `json:"variables"`
	IsDefault bool              `json:"isDefault"`
}

// DefaultTryItConfig returns sensible defaults
func DefaultTryItConfig() TryItConfig {
	return TryItConfig{
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
func (c TryItConfig) ToScalarConfig() map[string]interface{} {
	return map[string]interface{}{
		"hiddenClients": getHiddenClients(c.Snippets),
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
