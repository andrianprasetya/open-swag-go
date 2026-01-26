package schema

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

// FromType converts a Go type to JSON Schema
func FromType(t interface{}) *Schema {
	return fromReflectType(reflect.TypeOf(t))
}

func fromReflectType(t reflect.Type) *Schema {
	if t == nil {
		return &Schema{Type: "object"}
	}

	// Handle pointer types
	if t.Kind() == reflect.Ptr {
		return fromReflectType(t.Elem())
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
			Items: fromReflectType(t.Elem()),
		}
	case reflect.Struct:
		return fromStruct(t)
	case reflect.Map:
		return &Schema{Type: "object"}
	default:
		return &Schema{Type: "object"}
	}
}

func fromStruct(t reflect.Type) *Schema {
	schema := &Schema{
		Type:       "object",
		Properties: make(map[string]*Schema),
		Required:   []string{},
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		if !field.IsExported() {
			continue
		}

		jsonTag := field.Tag.Get("json")
		if jsonTag == "-" {
			continue
		}

		name := strings.Split(jsonTag, ",")[0]
		if name == "" {
			name = field.Name
		}

		fieldSchema := fromReflectType(field.Type)
		ParseFieldTags(field, fieldSchema)

		schema.Properties[name] = fieldSchema

		if IsRequired(field) {
			schema.Required = append(schema.Required, name)
		}
	}

	return schema
}
