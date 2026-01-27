package schema

import (
	"reflect"
	"strings"
	"time"
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
	if t == nil {
		return &Schema{Type: "object"}
	}
	return fromReflectType(reflect.TypeOf(t))
}

// FromReflectType converts a reflect.Type to JSON Schema
func FromReflectType(t reflect.Type) *Schema {
	return fromReflectType(t)
}

func fromReflectType(t reflect.Type) *Schema {
	if t == nil {
		return &Schema{Type: "object"}
	}

	// Handle pointer types
	if t.Kind() == reflect.Ptr {
		return fromReflectType(t.Elem())
	}

	// Handle time.Time specially
	if t == reflect.TypeOf(time.Time{}) {
		return &Schema{Type: "string", Format: "date-time"}
	}

	switch t.Kind() {
	case reflect.String:
		return &Schema{Type: "string"}
	case reflect.Int:
		return &Schema{Type: "integer"}
	case reflect.Int8, reflect.Int16:
		return &Schema{Type: "integer"}
	case reflect.Int32:
		return &Schema{Type: "integer", Format: "int32"}
	case reflect.Int64:
		return &Schema{Type: "integer", Format: "int64"}
	case reflect.Uint, reflect.Uint8, reflect.Uint16:
		return &Schema{Type: "integer"}
	case reflect.Uint32:
		return &Schema{Type: "integer", Format: "int32"}
	case reflect.Uint64:
		return &Schema{Type: "integer", Format: "int64"}
	case reflect.Float32:
		return &Schema{Type: "number", Format: "float"}
	case reflect.Float64:
		return &Schema{Type: "number", Format: "double"}
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
		return &Schema{
			Type: "object",
			// AdditionalProperties could be set here if needed
		}
	case reflect.Interface:
		return &Schema{Type: "object"}
	default:
		return &Schema{Type: "string"}
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

		// Get field name from json tag first, then form tag
		jsonTag := field.Tag.Get("json")
		if jsonTag == "-" {
			continue
		}

		name := strings.Split(jsonTag, ",")[0]
		if name == "" {
			// Fallback to form tag
			formTag := field.Tag.Get("form")
			if formTag != "" && formTag != "-" {
				name = strings.Split(formTag, ",")[0]
			}
		}
		if name == "" {
			name = field.Name
		}

		// Build schema from field type
		fieldSchema := fromReflectType(field.Type)

		// Parse additional tags
		ParseFieldTags(field, fieldSchema)

		schema.Properties[name] = fieldSchema

		// Check if required
		if IsRequired(field) {
			schema.Required = append(schema.Required, name)
		}
	}

	// Remove empty required array
	if len(schema.Required) == 0 {
		schema.Required = nil
	}

	return schema
}
