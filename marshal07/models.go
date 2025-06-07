package marshal07

import "encoding/json"

type Common struct {
	ID               *string           `json:"$id,omitempty" hcl:"$id,omitempty"`
	Ref              *string           `json:"$ref,omitempty" hcl:"$ref,omitempty"`
	Schema           *string           `json:"$schema,omitempty" hcl:"$schema,omitempty"`
	Format           *string           `json:"format,omitempty" hcl:"format,omitempty"`
	ContentMediaType *string           `json:"contentMediaType,omitempty" hcl:"contentMediaType,omitempty"`
	ContentEncoding  *string           `json:"contentEncoding,omitempty" hcl:"contentEncoding,omitempty"`
	Comment          *string           `json:"$comment,omitempty" hcl:"$comment,omitempty"`
	Title            *string           `json:"title,omitempty" hcl:"title,omitempty"`
	Description      *string           `json:"description,omitempty" hcl:"description,omitempty"`
	Const            *json.RawMessage  `json:"const,omitempty" hcl:"const,omitempty"`
	Enumeration      []SchemaEnumValue `json:"enum,omitempty" hcl:"enum,omitempty"`
	Default          *json.RawMessage  `json:"default,omitempty" hcl:"default,omitempty"`
}

type SchemaNumber struct {
	MultipleOf       *IntegerOrFloat `json:"multipleOf,omitempty" hcl:"multipleOf,omitempty"`
	Maximum          *IntegerOrFloat `json:"maximum,omitempty" hcl:"maximum,omitempty"`
	ExclusiveMaximum *IntegerOrFloat `json:"exclusiveMaximum,omitempty" hcl:"exclusiveMaximum,omitempty"`
	Minimum          *IntegerOrFloat `json:"minimum,omitempty" hcl:"minimum,omitempty"`
	ExclusiveMinimum *IntegerOrFloat `json:"exclusiveMinimum,omitempty" hcl:"exclusiveMinimum,omitempty"`
}

type SchemaString struct {
	MaxLength *int64  `json:"maxLength,omitempty" hcl:"maxLength,omitempty"`
	MinLength *int64  `json:"minLength,omitempty" hcl:"minLength,omitempty"`
	Pattern   *string `json:"pattern,omitempty" hcl:"pattern,omitempty"`
}

type SchemaArray struct {
	AdditionalItems *Combined                `json:"additionalItems,omitempty" hcl:"additionalItems,block"`
	Items           *CombinedOrCombinedArray `json:"items,omitempty" hcl:"items,block"`
	MaxItems        *int64                   `json:"maxItems,omitempty" hcl:"maxItems,omitempty"`
	MinItems        *int64                   `json:"minItems,omitempty" hcl:"minItems,omitempty"`
	UniqueItems     *bool                    `json:"uniqueItems,omitempty" hcl:"uniqueItems,omitempty"`
	Contains        *Combined                `json:"contains,omitempty" hcl:"contains,block"`
}

type SchemaObject struct {
	MaxProperties        *int64                            `json:"maxProperties,omitempty" hcl:"maxProperties,omitempty"`
	MinProperties        *int64                            `json:"minProperties,omitempty" hcl:"minProperties,omitempty"`
	Required             []string                          `json:"required,omitempty" hcl:"required,omitempty"`
	AdditionalProperties *Combined                         `json:"additionalProperties,omitempty" hcl:"additionalProperties,block"`
	PropertyNames        *Combined                         `json:"propertyNames,omitempty" hcl:"propertyNames,block"`
	Properties           map[string]*Combined              `json:"properties,omitempty" hcl:"properties,block"`
	PatternProperties    map[string]*Combined              `json:"patternProperties,omitempty" hcl:"patternProperties,block"`
	Dependencies         map[string]*CombinedOrStringArray `json:"dependencies,omitempty" hcl:"dependencies,block"`
}

// The Schema struct models a JSON Combined and, because schemas are
// defined hierarchically, contains many references to itself.
// All fields are pointers and are nil if the associated values
// are not specified.
type Schema struct {
	Type *StringOrStringArray `json:"type,omitempty" hcl:"type,omitempty"`
	Common
	Ref       *string           `json:"$ref,omitempty" hcl:"$ref,omitempty"`
	ReadOnly  *bool             `json:"readOnly,omitempty" hcl:"readOnly,omitempty"`
	WriteOnly *bool             `json:"writeOnly,omitempty" hcl:"writeOnly,omitempty"`
	Examples  []json.RawMessage `json:"examples,omitempty" hcl:"examples,block"`

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

// IsIntegerOrFloat returns true if the Schema is a number
func (s *Schema) IsIntegerOrFloat() bool {
	return s.MultipleOf != nil || s.Maximum != nil || s.ExclusiveMaximum != nil || s.Minimum != nil || s.ExclusiveMinimum != nil
}

// IsString returns true if the Schema is a string
func (s *Schema) IsString() bool {
	return s.MaxLength != nil || s.MinLength != nil || s.Pattern != nil
}

// IsArray returns true if the Schema is an array
func (s *Schema) IsArray() bool {
	return s.AdditionalItems != nil || s.Items != nil || s.MaxItems != nil || s.MinItems != nil || s.UniqueItems != nil || s.Contains != nil
}

// IsObject returns true if the Schema is an object
func (s *Schema) IsObject() bool {
	return s.MaxProperties != nil || s.MinProperties != nil || s.Required != nil || s.AdditionalProperties != nil || s.PropertyNames != nil ||
		s.Properties != nil || s.PatternProperties != nil || s.Dependencies != nil
}

// IsReference returns true if the Schema is a reference
func (s *Schema) IsReference() bool {
	return s.Ref != nil
}

// IsOnlyReference returns true if the Schema is a reference
// and does not have any other properties set.
func (s *Schema) IsOnlyReference() bool {
	return s.Ref != nil &&
		s.MultipleOf == nil && s.Maximum == nil && s.ExclusiveMaximum == nil && s.Minimum == nil && s.ExclusiveMinimum == nil &&
		s.MaxLength == nil && s.MinLength == nil && s.Pattern == nil &&
		s.AdditionalItems == nil && s.Items == nil && s.MaxItems == nil && s.MinItems == nil && s.UniqueItems == nil && s.Contains == nil &&
		s.MaxProperties == nil && s.MinProperties == nil && s.Required == nil && s.AdditionalProperties == nil && s.PropertyNames == nil &&
		s.Properties == nil && s.PatternProperties == nil && s.Dependencies == nil &&
		s.Default == nil && s.ReadOnly == nil && s.WriteOnly == nil &&
		s.Definitions == nil
}
