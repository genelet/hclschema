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
// of JSON Combineds.
package jsonschema07

import "gopkg.in/yaml.v3"

// The Schema struct models a JSON Combined and, because schemas are
// defined hierarchically, contains many references to itself.
// All fields are pointers and are nil if the associated values
// are not specified.
type Schema struct {
	ID          *string `json:"$id,omitempty"`
	Schema      *string `json:"$schema,omitempty"`
	Ref         *string `json:"$ref,omitempty"`
	Comment     *string `json:"$comment,omitempty"`
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`

	Default   *yaml.Node   `json:"default,omitempty"`
	ReadOnly  *bool        `json:"readOnly,omitempty"`
	WriteOnly *bool        `json:"writeOnly,omitempty"`
	Examples  *[]*Combined `json:"examples,omitempty"`

	MultipleOf       *SchemaNumber `json:"multipleOf,omitempty"`
	Maximum          *SchemaNumber `json:"maximum,omitempty"`
	ExclusiveMaximum *SchemaNumber `json:"exclusiveMaximum,omitempty"`
	Minimum          *SchemaNumber `json:"minimum,omitempty"`
	ExclusiveMinimum *SchemaNumber `json:"exclusiveMinimum,omitempty"`

	MaxLength *int64  `json:"maxLength,omitempty"`
	MinLength *int64  `json:"minLength,omitempty"`
	Pattern   *string `json:"pattern,omitempty"`

	AdditionalItems *Combined                `json:"additionalItems,omitempty"`
	Items           *CombinedOrCombinedArray `json:"items,omitempty"`
	MaxItems        *int64                   `json:"maxItems,omitempty"`
	MinItems        *int64                   `json:"minItems,omitempty"`
	UniqueItems     *bool                    `json:"uniqueItems,omitempty"`

	Contains             *Combined                      `json:"contains,omitempty"`
	MaxProperties        *int64                         `json:"maxProperties,omitempty"`
	MinProperties        *int64                         `json:"minProperties,omitempty"`
	Required             *[]string                      `json:"required,omitempty"`
	AdditionalProperties *Combined                      `json:"additionalProperties,omitempty"`
	Definitions          *[]*NamedCombined              `json:"definitions,omitempty"`
	Properties           *[]*NamedCombined              `json:"properties,omitempty"`
	PatternProperties    *[]*NamedCombined              `json:"patternProperties,omitempty"`
	Dependencies         *[]*NamedCombinedOrStringArray `json:"dependencies,omitempty"`
	PropertyNames        *Combined                      `json:"propertyNames,omitempty"`

	Const            *yaml.Node `json:"const,omitempty"`
	Enumeration      *[]SchemaEnumValue
	Type             *StringOrStringArray `json:"type,omitempty"`
	Format           *string              `json:"format,omitempty"`
	ContentMediaType *string              `json:"contentMediaType,omitempty"`
	ContentEncoding  *string              `json:"contentEncoding,omitempty"`

	If    *Combined    `json:"if,omitempty"`
	Then  *Combined    `json:"then,omitempty"`
	Else  *Combined    `json:"else,omitempty"`
	AllOf *[]*Combined `json:"allOf,omitempty"`
	AnyOf *[]*Combined `json:"anyOf,omitempty"`
	OneOf *[]*Combined `json:"oneOf,omitempty"`
	Not   *Combined    `json:"not,omitempty"`
}

// Combined represents a value that can be either a Schema or a Boolean.
type Combined struct {
	Schema  *Schema
	Boolean *bool
}

// NewCombinedWithSchema creates and returns a new object
func NewCombinedWithSchema(s *Schema) *Combined {
	result := &Combined{}
	result.Schema = s
	return result
}

// NewCombinedWithBoolean creates and returns a new object
func NewCombinedWithBoolean(b bool) *Combined {
	result := &Combined{}
	result.Boolean = &b
	return result
}

// These helper structs represent "combination" types that generally can
// have values of one type or another. All are used to represent parts
// of Combineds.

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

// CombinedOrStringArray represents a value that can be either
// a Combined or an Array of Strings.
type CombinedOrStringArray struct {
	Combined    *Combined
	StringArray *[]string
}

// CombinedOrCombinedArray represents a value that can be either
// a Combined or an Array of Combineds.
type CombinedOrCombinedArray struct {
	Combined      *Combined
	CombinedArray *[]*Combined
}

// NewCombinedOrCombinedArrayWithCombined creates and returns a new object
func NewCombinedOrCombinedArrayWithCombined(s *Combined) *CombinedOrCombinedArray {
	result := &CombinedOrCombinedArray{}
	result.Combined = s
	return result
}

// NewCombinedOrCombinedArrayWithCombinedArray creates and returns a new object
func NewCombinedOrCombinedArrayWithCombinedArray(a []*Combined) *CombinedOrCombinedArray {
	result := &CombinedOrCombinedArray{}
	result.CombinedArray = &a
	return result
}

// SchemaEnumValue represents a value that can be part of an
// enumeration in a Combined.
type SchemaEnumValue struct {
	String *string
	Bool   *bool
}

// NamedCombined is a name-value pair that is used to emulate maps
// with ordered keys.
type NamedCombined struct {
	Name  string
	Value *Combined
}

// NewNamedCombined creates and returns a new object
func NewNamedCombined(name string, value *Combined) *NamedCombined {
	return &NamedCombined{Name: name, Value: value}
}

// NamedCombinedOrStringArray is a name-value pair that is used
// to emulate maps with ordered keys.
type NamedCombinedOrStringArray struct {
	Name  string
	Value *CombinedOrStringArray
}

// Access named subschemas by name

func namedCombinedArrayElementWithName(array *[]*NamedCombined, name string) *Combined {
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
func (s *Schema) PropertyWithName(name string) *Combined {
	return namedCombinedArrayElementWithName(s.Properties, name)
}

// PatternPropertyWithName returns the selected element.
func (s *Schema) PatternPropertyWithName(name string) *Combined {
	return namedCombinedArrayElementWithName(s.PatternProperties, name)
}

// DefinitionWithName returns the selected element.
func (s *Schema) DefinitionWithName(name string) *Combined {
	return namedCombinedArrayElementWithName(s.Definitions, name)
}

// AddProperty adds a named property.
func (s *Schema) AddProperty(name string, property *Combined) {
	*s.Properties = append(*s.Properties, NewNamedCombined(name, property))
}
