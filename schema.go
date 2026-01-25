package openswag

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

// SchemaFromType converts a Go type to JSON Schema
func SchemaFromType(t interface{}) *Schema {
	return schemaFromReflectType(reflect.TypeOf(t))
}

func schemaFromReflectType(t reflect.Type) *Schema {
	if t == nil {
		return &Schema{Type: "object"}
	}

	// Handle pointer types
	if t.Kind() == reflect.Ptr {
		return schemaFromReflectType(t.Elem())
	}

	switch t.Kind() {
	case reflect.String:
		return &Schema{Type: "string"}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return &Schema{Type: "integer"}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return &Schema{Type: "integer"}
	case reflect.Float32, reflect.Float64:
		return &Schema{Type: "number"}
	case reflect.Bool:
		return &Schema{Type: "boolean"}
	case reflect.Slice, reflect.Array:
		return &Schema{
			Type:  "array",
			Items: schemaFromReflectType(t.Elem()),
		}
	case reflect.Struct:
		return schemaFromStruct(t)
	case reflect.Map:
		return &Schema{Type: "object"}
	default:
		return &Schema{Type: "object"}
	}
}

func schemaFromStruct(t reflect.Type) *Schema {
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

		// Parse field schema
		fieldSchema := schemaFromReflectType(field.Type)
		parseFieldTags(field, fieldSchema)

		schema.Properties[name] = fieldSchema

		// Check if required
		if isFieldRequired(field) {
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

	// Parse swagger tag
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

func isFieldRequired(field reflect.StructField) bool {
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
