package jsm07

import (
	jsonschema "github.com/genelet/hclschema/jsonschema07"
	"gopkg.in/yaml.v3"
)

type Reference struct {
	Ref *string
}

type Common struct {
	Type             *jsonschema.StringOrStringArray
	Format           *string
	Default          *yaml.Node
	Enumeration      []jsonschema.SchemaEnumValue
	Const            *yaml.Node
	ContentMediaType *string
	ContentEncoding  *string
	ID               *string
	Schema           *string
	Comment          *string
	Title            *string
	Description      *string
}

type CommonBoolean struct {
	Common
}

type SchemaNumber struct {
	MultipleOf       *jsonschema.SchemaNumber
	Maximum          *jsonschema.SchemaNumber
	ExclusiveMaximum *bool
	Minimum          *jsonschema.SchemaNumber
	ExclusiveMinimum *bool
}

type CommonNumber struct {
	Common
	SchemaNumber
}

type SchemaString struct {
	MaxLength *int64
	MinLength *int64
	Pattern   *string
}

type CommonString struct {
	Common
	SchemaString
}

type SchemaArray struct {
	Items       *SchemaOrSchemaArray
	MaxItems    *int64
	MinItems    *int64
	UniqueItems *bool
}

type CommonArray struct {
	Common
	SchemaArray
}

type SchemaOrStringArray struct {
	Schema      *Schema
	StringArray []string
}

type SchemaOrSchemaArray struct {
	Schema      *Schema
	SchemaArray []*Schema
}

type SchemaMap struct {
	AdditionalProperties *SchemaOrBoolean
}

type CommonMap struct {
	Common
	SchemaMap
}

type SchemaOrBoolean struct {
	Schema  *Schema
	Boolean *bool
}

type SchemaObject struct {
	MaxProperties *int64
	MinProperties *int64
	Required      []string
	Properties    map[string]*Schema
}

type CommonObject struct {
	Common
	SchemaObject
}

type CommonAllOf struct {
	Common
	SchemaObject
	AllOf []*Schema
}

type CommonAnyOf struct {
	Common
	SchemaObject
	AnyOf []*Schema
}

type CommonOneOf struct {
	Common
	SchemaObject
	OneOf []*Schema
}

type CommonIfThenElse struct {
	Common
	SchemaObject
	If   *Schema
	Then *Schema
	Else *Schema
}

type SchemaFull struct {
	Schema *string
	ID     *string
	*Reference
	ReadOnly  *bool
	WriteOnly *bool

	*Common
	*SchemaNumber
	*SchemaString
	*SchemaArray
	AdditionalItems *SchemaOrBoolean
	*SchemaMap
	*SchemaObject
	PatternProperties map[string]*Schema
	Dependencies      map[string]*SchemaOrStringArray

	AllOf       []*Schema
	AnyOf       []*Schema
	OneOf       []*Schema
	Not         *Schema
	Definitions map[string]*Schema

	Title       *string
	Description *string
}

type Schema struct {
	*Reference
	*CommonBoolean
	*CommonNumber
	*CommonString
	*CommonArray
	*CommonMap
	*CommonObject
	*CommonAllOf
	*CommonAnyOf
	*CommonOneOf
	*CommonIfThenElse
	*SchemaFull
}

func (s *Schema) IsFull() bool {
	return s != nil && s.SchemaFull != nil
}

func isCommon(s *jsonschema.Schema) bool {
	return s.Type != nil || s.Format != nil || s.Default != nil || s.Enumeration != nil
}
func isNumber(s *jsonschema.Schema) bool {
	return s.MultipleOf != nil || s.Maximum != nil || s.ExclusiveMaximum != nil || s.Minimum != nil || s.ExclusiveMinimum != nil
}
func isString(s *jsonschema.Schema) bool {
	return s.MaxLength != nil || s.MinLength != nil || s.Pattern != nil
}
func isArray(s *jsonschema.Schema) bool {
	return s.Items != nil || s.MaxItems != nil || s.MinItems != nil || s.UniqueItems != nil
}
func isObject(s *jsonschema.Schema) bool {
	return s.MaxProperties != nil || s.MinProperties != nil || (s.Properties != nil && len(*s.Properties) > 0) || (s.Required != nil && len(*s.Required) > 0)
}
func isMap(s *jsonschema.Schema) bool {
	return s.AdditionalProperties != nil
}
func isReference(s *jsonschema.Schema) bool {
	return s.Ref != nil
}
func isOnlyReference(s *jsonschema.Schema) bool {
	return s.Ref != nil && !isCommon(s) && !isNumber(s) && !isString(s) && !isArray(s) && !isObject(s) && !isMap(s)
}

func isRest(s *jsonschema.Schema) bool {
	return s.Schema != nil || s.ID != nil || s.ReadOnly != nil || s.WriteOnly != nil || (s.PatternProperties != nil && len(*s.PatternProperties) > 0) || (s.Dependencies != nil && len(*s.Dependencies) > 0) || (s.AllOf != nil && len(*s.AllOf) > 0) || (s.AnyOf != nil && len(*s.AnyOf) > 0) || (s.OneOf != nil && len(*s.OneOf) > 0) || s.Not != nil || (s.Definitions != nil && len(*s.Definitions) > 0) || s.Title != nil || s.Description != nil
}

func isFull(s *jsonschema.Schema) bool {
	if isRest(s) {
		return true
	}

	if isReference(s) {
		return !isOnlyReference(s)
	}

	if s.Type == nil || s.Type.String == nil {
		return true
	}

	switch *s.Type.String {
	case "boolean":
		return isString(s) || isNumber(s) || isArray(s) || isObject(s) || isMap(s)
	case "number", "integer":
		return isString(s) || isArray(s) || isObject(s) || isMap(s)
	case "string":
		return isNumber(s) || isArray(s) || isObject(s) || isMap(s)
	case "array":
		return isString(s) || isNumber(s) || isObject(s) || isMap(s)
	case "object":
		if s.AdditionalProperties != nil {
			return isString(s) || isNumber(s) || isArray(s) || isObject(s)
		} else {
			return isString(s) || isNumber(s) || isArray(s) || isMap(s)
		}
	default:
	}
	return true
}
