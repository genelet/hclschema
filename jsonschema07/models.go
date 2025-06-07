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

type Common struct {
	ID               *string              `json:"$id,omitempty" hcl:"$id,omitempty"`
	Ref              *string              `json:"$ref,omitempty" hcl:"$ref,omitempty"`
	Schema           *string              `json:"$schema,omitempty" hcl:"$schema,omitempty"`
	Type             *StringOrStringArray `json:"type,omitempty" hcl:"type,omitempty"`
	Format           *string              `json:"format,omitempty" hcl:"format,omitempty"`
	ContentMediaType *string              `json:"contentMediaType,omitempty" hcl:"contentMediaType,omitempty"`
	ContentEncoding  *string              `json:"contentEncoding,omitempty" hcl:"contentEncoding,omitempty"`
	Comment          *string              `json:"$comment,omitempty" hcl:"$comment,omitempty"`
	Title            *string              `json:"title,omitempty" hcl:"title,omitempty"`
	Description      *string              `json:"description,omitempty" hcl:"description,omitempty"`
	Const            *yaml.Node           `json:"const,omitempty" hcl:"const,omitempty"`
	Enumeration      []SchemaEnumValue    `json:"enum,omitempty" hcl:"enum,omitempty"`
	Default          *yaml.Node           `json:"default,omitempty" hcl:"default,omitempty"`
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
	MaxProperties        *int64                        `json:"maxProperties,omitempty" hcl:"maxProperties,omitempty"`
	MinProperties        *int64                        `json:"minProperties,omitempty" hcl:"minProperties,omitempty"`
	Required             []string                      `json:"required,omitempty" hcl:"required,omitempty"`
	AdditionalProperties *Combined                     `json:"additionalProperties,omitempty" hcl:"additionalProperties,block"`
	PropertyNames        *Combined                     `json:"propertyNames,omitempty" hcl:"propertyNames,block"`
	Properties           []*NamedCombined              `json:"properties,omitempty" hcl:"properties,block"`
	PatternProperties    []*NamedCombined              `json:"patternProperties,omitempty" hcl:"patternProperties,block"`
	Dependencies         []*NamedCombinedOrStringArray `json:"dependencies,omitempty" hcl:"dependencies,block"`
}

// The Schema struct models a JSON Combined and, because schemas are
// defined hierarchically, contains many references to itself.
// All fields are pointers and are nil if the associated values
// are not specified.
type Schema struct {
	Common
	Ref       *string     `json:"$ref,omitempty" hcl:"$ref,omitempty"`
	ReadOnly  *bool       `json:"readOnly,omitempty" hcl:"readOnly,omitempty"`
	WriteOnly *bool       `json:"writeOnly,omitempty" hcl:"writeOnly,omitempty"`
	Examples  []*Combined `json:"examples,omitempty" hcl:"examples,block"`

	SchemaNumber
	SchemaString
	SchemaArray
	SchemaObject

	Definitions []*NamedCombined `json:"definitions,omitempty" hcl:"definitions,block"`

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

// IntegerOrFloat represents a value that can be either an Integer or a Float.
type IntegerOrFloat struct {
	Integer *int64
	Float   *float64
}

// NewIntegerOrFloatWithInteger creates and returns a new object
func NewIntegerOrFloatWithInteger(i int64) *IntegerOrFloat {
	result := &IntegerOrFloat{}
	result.Integer = &i
	return result
}

// NewIntegerOrFloatWithFloat creates and returns a new object
func NewIntegerOrFloatWithFloat(f float64) *IntegerOrFloat {
	result := &IntegerOrFloat{}
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

func NewCombinedOrStringArrayWithCombined(s *Combined) *CombinedOrStringArray {
	result := &CombinedOrStringArray{}
	result.Combined = s
	return result
}
func NewCombinedOrStringArrayWithStringArray(a []string) *CombinedOrStringArray {
	result := &CombinedOrStringArray{}
	result.StringArray = &a
	return result
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
	Number *IntegerOrFloat
	Null   *SchemaNull
}

type SchemaNull struct {
	IsNull *bool `json:"isNull,omitempty" hcl:"isNull,omitempty"`
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

func namedCombinedArrayElementWithName(array []*NamedCombined, name string) *Combined {
	if array == nil {
		return nil
	}
	for _, pair := range array {
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
	s.Properties = append(s.Properties, NewNamedCombined(name, property))
}
