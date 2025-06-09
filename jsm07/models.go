package jsm07

import (
	"encoding/json"
)

type Common struct {
	ID               *string           `json:"$id,omitempty" hcl:"_id,optional"`
	Schema           *string           `json:"$schema,omitempty" hcl:"_schema,optional"`
	Format           *string           `json:"format,omitempty" hcl:"format,optional"`
	ContentMediaType *string           `json:"contentMediaType,omitempty" hcl:"contentMediaType,optional"`
	ContentEncoding  *string           `json:"contentEncoding,omitempty" hcl:"contentEncoding,optional"`
	Comment          *string           `json:"$comment,omitempty" hcl:"_comment,optional"`
	Title            *string           `json:"title,omitempty" hcl:"title,optional"`
	Description      *string           `json:"description,omitempty" hcl:"description,optional"`
	Enumeration      []SchemaEnumValue `json:"enum,omitempty" hcl:"enum,optional"`
	Const            *json.RawMessage  `json:"const,omitempty" hcl:"const,optional"`
	Default          *json.RawMessage  `json:"default,omitempty" hcl:"default,optional"`
	Examples         *json.RawMessage  `json:"examples,omitempty" hcl:"examples,optional"`
}

type SchemaNumber struct {
	MultipleOf       *IntegerOrFloat `json:"multipleOf,omitempty" hcl:"multipleOf,optional"`
	Maximum          *IntegerOrFloat `json:"maximum,omitempty" hcl:"maximum,optional"`
	ExclusiveMaximum *IntegerOrFloat `json:"exclusiveMaximum,omitempty" hcl:"exclusiveMaximum,optional"`
	Minimum          *IntegerOrFloat `json:"minimum,omitempty" hcl:"minimum,optional"`
	ExclusiveMinimum *IntegerOrFloat `json:"exclusiveMinimum,omitempty" hcl:"exclusiveMinimum,optional"`
}

type SchemaString struct {
	MaxLength *int64  `json:"maxLength,omitempty" hcl:"maxLength,optional"`
	MinLength *int64  `json:"minLength,omitempty" hcl:"minLength,optional"`
	Pattern   *string `json:"pattern,omitempty" hcl:"pattern,optional"`
}

type SchemaArray struct {
	AdditionalItems *Combined                `json:"additionalItems,omitempty" hcl:"additionalItems,block"`
	Items           *CombinedOrCombinedArray `json:"items,omitempty" hcl:"items,block"`
	MaxItems        *int64                   `json:"maxItems,omitempty" hcl:"maxItems,optional"`
	MinItems        *int64                   `json:"minItems,omitempty" hcl:"minItems,optional"`
	UniqueItems     *bool                    `json:"uniqueItems,omitempty" hcl:"uniqueItems,optional"`
	Contains        *Combined                `json:"contains,omitempty" hcl:"contains,block"`
}

type SchemaObject struct {
	MaxProperties        *int64                            `json:"maxProperties,omitempty" hcl:"maxProperties,optional"`
	MinProperties        *int64                            `json:"minProperties,omitempty" hcl:"minProperties,optional"`
	Required             []string                          `json:"required,omitempty" hcl:"required,optional"`
	AdditionalProperties *Combined                         `json:"additionalProperties,omitempty" hcl:"additionalProperties,block"`
	PropertyNames        *Combined                         `json:"propertyNames,omitempty" hcl:"propertyNames,block"`
	Properties           map[string]*Combined              `json:"properties" hcl:"properties,block"`
	PatternProperties    map[string]*Combined              `json:"patternProperties,omitempty" hcl:"patternProperties,block"`
	Dependencies         map[string]*CombinedOrStringArray `json:"dependencies,omitempty" hcl:"dependencies,block"`
}

// The Schema struct models a JSON Combined and, because schemas are
// defined hierarchically, contains many references to itself.
// All fields are pointers and are nil if the associated values
// are not specified.
type Schema struct {
	Type *StringOrStringArray `json:"type,omitempty" hcl:"type,optional"`
	Common
	Ref       *string `json:"$ref,omitempty" hcl:"_ref,optional"`
	ReadOnly  *bool   `json:"readOnly,omitempty" hcl:"readOnly,optional"`
	WriteOnly *bool   `json:"writeOnly,omitempty" hcl:"writeOnly,optional"`

	SchemaNumber
	SchemaString
	SchemaArray
	SchemaObject

	Definitions map[string]*Combined `json:"definitions,omitempty" hcl:"definitions,block"`

	If    *Combined   `json:"if,omitempty" hcl:"if,block"`
	Then  *Combined   `json:"then,omitempty" hcl:"then,block"`
	Else  *Combined   `json:"else,omitempty" hcl:"else,block"`
	AllOf []*Combined `json:"allOf,omitempty" hcl:"allOf,block"`
	AnyOf []*Combined `json:"anyOf,omitempty" hcl:"anyOf,block"`
	OneOf []*Combined `json:"oneOf,omitempty" hcl:"oneOf,block"`
	Not   *Combined   `json:"not,omitempty" hcl:"not,block"`
}
