// Copyright 2017 Google LLC. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package jsonschema supports the reading, writing, and manipulation
// of JSON Schemas.
package jsonschema07

import "gopkg.in/yaml.v3"

// Schema represents a value that can be either a Absolute or a Boolean.
type Schema struct {
	Absolute *Absolute
	Boolean  *bool
}

// NewSchemaWithAbsolute creates and returns a new object
func NewSchemaWithAbsolute(s *Absolute) *Schema {
	result := &Schema{}
	result.Absolute = s
	return result
}

// NewSchemaWithBoolean creates and returns a new object
func NewSchemaWithBoolean(b bool) *Schema {
	result := &Schema{}
	result.Boolean = &b
	return result
}

// The Absolute struct models a JSON Schema and, because schemas are
// defined hierarchically, contains many references to itself.
// All fields are pointers and are nil if the associated values
// are not specified.
type Absolute struct {
	ID          *string `json:"$id,omitempty"`
	Schema      *string `json:"$schema,omitempty"`
	Ref         *string `json:"$ref,omitempty"`
	Comment     *string `json:"$comment,omitempty"`
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`

	Default   *yaml.Node `json:"default,omitempty"`
	ReadOnly  *bool      `json:"readOnly,omitempty"`
	WriteOnly *bool      `json:"writeOnly,omitempty"`
	Examples  *[]*Schema `json:"examples,omitempty"`

	MultipleOf       *SchemaNumber `json:"multipleOf,omitempty"`
	Maximum          *SchemaNumber `json:"maximum,omitempty"`
	ExclusiveMaximum *SchemaNumber `json:"exclusiveMaximum,omitempty"`
	Minimum          *SchemaNumber `json:"minimum,omitempty"`
	ExclusiveMinimum *SchemaNumber `json:"exclusiveMinimum,omitempty"`

	MaxLength *int64  `json:"maxLength,omitempty"`
	MinLength *int64  `json:"minLength,omitempty"`
	Pattern   *string `json:"pattern,omitempty"`

	AdditionalItems *Schema              `json:"additionalItems,omitempty"`
	Items           *SchemaOrSchemaArray `json:"items,omitempty"`
	MaxItems        *int64               `json:"maxItems,omitempty"`
	MinItems        *int64               `json:"minItems,omitempty"`
	UniqueItems     *bool                `json:"uniqueItems,omitempty"`

	Contains             *Schema                      `json:"contains,omitempty"`
	MaxProperties        *int64                       `json:"maxProperties,omitempty"`
	MinProperties        *int64                       `json:"minProperties,omitempty"`
	Required             *[]string                    `json:"required,omitempty"`
	AdditionalProperties *Schema                      `json:"additionalProperties,omitempty"`
	Definitions          *[]*NamedSchema              `json:"definitions,omitempty"`
	Properties           *[]*NamedSchema              `json:"properties,omitempty"`
	PatternProperties    *[]*NamedSchema              `json:"patternProperties,omitempty"`
	Dependencies         *[]*NamedSchemaOrStringArray `json:"dependencies,omitempty"`
	PropertyNames        *Schema                      `json:"propertyNames,omitempty"`

	Const            *yaml.Node `json:"const,omitempty"`
	Enumeration      *[]SchemaEnumValue
	Type             *StringOrStringArray `json:"type,omitempty"`
	Format           *string              `json:"format,omitempty"`
	ContentMediaType *string              `json:"contentMediaType,omitempty"`
	ContentEncoding  *string              `json:"contentEncoding,omitempty"`

	If    *Schema    `json:"if,omitempty"`
	Then  *Schema    `json:"then,omitempty"`
	Else  *Schema    `json:"else,omitempty"`
	AllOf *[]*Schema `json:"allOf,omitempty"`
	AnyOf *[]*Schema `json:"anyOf,omitempty"`
	OneOf *[]*Schema `json:"oneOf,omitempty"`
	Not   *Schema    `json:"not,omitempty"`
}

// These helper structs represent "combination" types that generally can
// have values of one type or another. All are used to represent parts
// of Schemas.

// SchemaNumber represents a value that can be either an Integer or a Float.
type SchemaNumber struct {
	Integer *int64
	Float   *float64
}

// NewSchemaNumberWithInteger creates and returns a new object
func NewSchemaNumberWithInteger(i int64) *SchemaNumber {
	result := &SchemaNumber{}
	result.Integer = &i
	return result
}

// NewSchemaNumberWithFloat creates and returns a new object
func NewSchemaNumberWithFloat(f float64) *SchemaNumber {
	result := &SchemaNumber{}
	result.Float = &f
	return result
}

// StringOrStringArray represents a value that can be either
// a String or an Array of Strings.
type StringOrStringArray struct {
	String      *string
	StringArray *[]string
}

// NewStringOrStringArrayWithString creates and returns a new object
func NewStringOrStringArrayWithString(s string) *StringOrStringArray {
	result := &StringOrStringArray{}
	result.String = &s
	return result
}

// NewStringOrStringArrayWithStringArray creates and returns a new object
func NewStringOrStringArrayWithStringArray(a []string) *StringOrStringArray {
	result := &StringOrStringArray{}
	result.StringArray = &a
	return result
}

// SchemaOrStringArray represents a value that can be either
// a Schema or an Array of Strings.
type SchemaOrStringArray struct {
	Schema      *Schema
	StringArray *[]string
}

// SchemaOrSchemaArray represents a value that can be either
// a Schema or an Array of Schemas.
type SchemaOrSchemaArray struct {
	Schema      *Schema
	SchemaArray *[]*Schema
}

// NewSchemaOrSchemaArrayWithSchema creates and returns a new object
func NewSchemaOrSchemaArrayWithSchema(s *Schema) *SchemaOrSchemaArray {
	result := &SchemaOrSchemaArray{}
	result.Schema = s
	return result
}

// NewSchemaOrSchemaArrayWithSchemaArray creates and returns a new object
func NewSchemaOrSchemaArrayWithSchemaArray(a []*Schema) *SchemaOrSchemaArray {
	result := &SchemaOrSchemaArray{}
	result.SchemaArray = &a
	return result
}

// SchemaEnumValue represents a value that can be part of an
// enumeration in a Schema.
type SchemaEnumValue struct {
	String *string
	Bool   *bool
}

// NamedSchema is a name-value pair that is used to emulate maps
// with ordered keys.
type NamedSchema struct {
	Name  string
	Value *Schema
}

// NewNamedSchema creates and returns a new object
func NewNamedSchema(name string, value *Schema) *NamedSchema {
	return &NamedSchema{Name: name, Value: value}
}

// NamedSchemaOrStringArray is a name-value pair that is used
// to emulate maps with ordered keys.
type NamedSchemaOrStringArray struct {
	Name  string
	Value *SchemaOrStringArray
}

// Access named subschemas by name

func namedSchemaArrayElementWithName(array *[]*NamedSchema, name string) *Schema {
	if array == nil {
		return nil
	}
	for _, pair := range *array {
		if pair.Name == name {
			return pair.Value
		}
	}
	return nil
}

// PropertyWithName returns the selected element.
func (s *Absolute) PropertyWithName(name string) *Schema {
	return namedSchemaArrayElementWithName(s.Properties, name)
}

// PatternPropertyWithName returns the selected element.
func (s *Absolute) PatternPropertyWithName(name string) *Schema {
	return namedSchemaArrayElementWithName(s.PatternProperties, name)
}

// DefinitionWithName returns the selected element.
func (s *Absolute) DefinitionWithName(name string) *Schema {
	return namedSchemaArrayElementWithName(s.Definitions, name)
}

// AddProperty adds a named property.
func (s *Absolute) AddProperty(name string, property *Schema) {
	*s.Properties = append(*s.Properties, NewNamedSchema(name, property))
}
