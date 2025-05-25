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

package jsonschema07

import (
	"fmt"
	"log"
	"strings"
)

//
// OPERATIONS
// The following methods perform operations on Combineds.
//

// IsEmpty returns true if no members of the Combined are specified.
func (schema *Schema) IsEmpty() bool {
	return (schema.ID == nil) &&
		(schema.Comment == nil) &&
		(schema.Schema == nil) &&
		(schema.Ref == nil) &&
		(schema.Title == nil) &&
		(schema.Description == nil) &&

		(schema.Default == nil) &&
		(schema.ReadOnly == nil) &&
		(schema.WriteOnly == nil) &&
		(schema.Examples == nil) &&

		(schema.MultipleOf == nil) &&
		(schema.Maximum == nil) &&
		(schema.ExclusiveMaximum == nil) &&
		(schema.Minimum == nil) &&
		(schema.ExclusiveMinimum == nil) &&

		(schema.MaxLength == nil) &&
		(schema.MinLength == nil) &&
		(schema.Pattern == nil) &&

		(schema.AdditionalItems == nil) &&
		(schema.Items == nil) &&
		(schema.MaxItems == nil) &&
		(schema.MinItems == nil) &&
		(schema.UniqueItems == nil) &&

		(schema.Contains == nil) &&
		(schema.MaxProperties == nil) &&
		(schema.MinProperties == nil) &&
		(schema.Required == nil) &&
		(schema.AdditionalProperties == nil) &&
		(schema.Definitions == nil) &&
		(schema.Properties == nil) &&
		(schema.PatternProperties == nil) &&
		(schema.Dependencies == nil) &&
		(schema.PropertyNames == nil) &&

		(schema.Const == nil) &&
		(schema.Enumeration == nil) &&
		(schema.Type == nil) &&
		(schema.Format == nil) &&
		(schema.ContentMediaType == nil) &&
		(schema.ContentEncoding == nil) &&

		(schema.If == nil) &&
		(schema.Then == nil) &&
		(schema.Else == nil) &&
		(schema.AllOf == nil) &&
		(schema.AnyOf == nil) &&
		(schema.OneOf == nil) &&
		(schema.Not == nil)
}

// IsEqual returns true if two schemas are equal.
func (schema *Schema) IsEqual(schema2 *Schema) bool {
	return schema.String() == schema2.String()
}

func (combined *Combined) IsEqual(combined2 *Combined) bool {
	return combined.String() == combined2.String()
}

// CombinedOperation represents a function that can be applied to a Combined.
type CombinedOperation func(combined *Combined, context string)

// Applies a specified function to a Combined and all of the Combineds that it contains.
func (combined *Combined) applyToCombineds(operation CombinedOperation, context string) {
	if combined.Boolean != nil {
		return
	}

	schema := combined.Schema

	if schema.AdditionalItems != nil {
		schema.AdditionalItems.applyToCombineds(operation, "AdditionalItems")
	}

	if schema.Items != nil {
		if schema.Items.CombinedArray != nil {
			for _, s := range *(schema.Items.CombinedArray) {
				s.applyToCombineds(operation, "Items.CombinedArray")
			}
		} else if schema.Items.Combined != nil {
			schema.Items.Combined.applyToCombineds(operation, "Items.Combined")
		}
	}

	if schema.Contains != nil {
		schema.Contains.applyToCombineds(operation, "Contains")
	}

	if schema.AdditionalProperties != nil {
		schema.AdditionalProperties.applyToCombineds(operation, "AdditionalProperties")
	}

	if schema.Definitions != nil {
		for _, pair := range *(schema.Definitions) {
			s := pair.Value
			s.applyToCombineds(operation, "Definitions")
		}
	}
	if schema.Properties != nil {
		for _, pair := range *(schema.Properties) {
			s := pair.Value
			s.applyToCombineds(operation, "Properties")
		}
	}
	if schema.PatternProperties != nil {
		for _, pair := range *(schema.PatternProperties) {
			s := pair.Value
			s.applyToCombineds(operation, "PatternProperties")
		}
	}

	if schema.Dependencies != nil {
		for _, pair := range *(schema.Dependencies) {
			schemaOrStringArray := pair.Value
			s := schemaOrStringArray.Combined
			if s != nil {
				s.applyToCombineds(operation, "Dependencies")
			}
		}
	}

	if schema.PropertyNames != nil {
		schema.PropertyNames.applyToCombineds(operation, "PropertyNames")
	}

	if schema.If != nil {
		schema.If.applyToCombineds(operation, "If")
	}
	if schema.Then != nil {
		schema.Then.applyToCombineds(operation, "Then")
	}
	if schema.Else != nil {
		schema.Else.applyToCombineds(operation, "Else")
	}

	if schema.AllOf != nil {
		for _, s := range *(schema.AllOf) {
			s.applyToCombineds(operation, "AllOf")
		}
	}
	if schema.AnyOf != nil {
		for _, s := range *(schema.AnyOf) {
			s.applyToCombineds(operation, "AnyOf")
		}
	}
	if schema.OneOf != nil {
		for _, s := range *(schema.OneOf) {
			s.applyToCombineds(operation, "OneOf")
		}
	}
	if schema.Not != nil {
		schema.Not.applyToCombineds(operation, "Not")
	}

	operation(combined, context)
}

// CopyProperties copies all non-nil properties from the source Combined to the schema Combined.
func (combined *Combined) CopyProperties(source *Combined) {
	if source.Boolean != nil {
		combined.Boolean = source.Boolean
		return
	}
	combined.Schema.CopyProperties(source.Schema)
}

func (schema *Schema) CopyProperties(source *Schema) {
	if source.Schema != nil {
		schema.Schema = source.Schema
	}
	if source.ID != nil {
		schema.ID = source.ID
	}
	if source.Comment != nil {
		schema.Comment = source.Comment
	}
	if source.Ref != nil {
		schema.Ref = source.Ref
	}
	if source.Title != nil {
		schema.Title = source.Title
	}
	if source.Description != nil {
		schema.Description = source.Description
	}

	if source.Default != nil {
		schema.Default = source.Default
	}
	if source.ReadOnly != nil {
		schema.ReadOnly = source.ReadOnly
	}
	if source.WriteOnly != nil {
		schema.WriteOnly = source.WriteOnly
	}
	if source.Examples != nil {
		schema.Examples = source.Examples
	}

	if source.MultipleOf != nil {
		schema.MultipleOf = source.MultipleOf
	}
	if source.Maximum != nil {
		schema.Maximum = source.Maximum
	}
	if source.ExclusiveMaximum != nil {
		schema.ExclusiveMaximum = source.ExclusiveMaximum
	}
	if source.Minimum != nil {
		schema.Minimum = source.Minimum
	}
	if source.ExclusiveMinimum != nil {
		schema.ExclusiveMinimum = source.ExclusiveMinimum
	}

	if source.MaxLength != nil {
		schema.MaxLength = source.MaxLength
	}
	if source.MinLength != nil {
		schema.MinLength = source.MinLength
	}
	if source.Pattern != nil {
		schema.Pattern = source.Pattern
	}

	if source.AdditionalItems != nil {
		schema.AdditionalItems = source.AdditionalItems
	}
	if source.Items != nil {
		schema.Items = source.Items
	}
	if source.MaxItems != nil {
		schema.MaxItems = source.MaxItems
	}
	if source.MinItems != nil {
		schema.MinItems = source.MinItems
	}
	if source.UniqueItems != nil {
		schema.UniqueItems = source.UniqueItems
	}

	if source.Contains != nil {
		schema.Contains = source.Contains
	}
	if source.MaxProperties != nil {
		schema.MaxProperties = source.MaxProperties
	}
	if source.MinProperties != nil {
		schema.MinProperties = source.MinProperties
	}
	if source.Required != nil {
		schema.Required = source.Required
	}
	if source.AdditionalProperties != nil {
		schema.AdditionalProperties = source.AdditionalProperties
	}
	if source.Definitions != nil {
		schema.Definitions = source.Definitions
	}
	if source.Properties != nil {
		schema.Properties = source.Properties
	}
	if source.PatternProperties != nil {
		schema.PatternProperties = source.PatternProperties
	}
	if source.Dependencies != nil {
		schema.Dependencies = source.Dependencies
	}
	if source.PropertyNames != nil {
		schema.PropertyNames = source.PropertyNames
	}

	if source.Const != nil {
		schema.Const = source.Const
	}
	if source.Enumeration != nil {
		schema.Enumeration = source.Enumeration
	}
	if source.Type != nil {
		schema.Type = source.Type
	}
	if source.Format != nil {
		schema.Format = source.Format
	}
	if source.ContentMediaType != nil {
		schema.ContentMediaType = source.ContentMediaType
	}
	if source.ContentEncoding != nil {
		schema.ContentEncoding = source.ContentEncoding
	}

	if source.If != nil {
		schema.If = source.If
	}
	if source.Then != nil {
		schema.Then = source.Then
	}
	if source.Else != nil {
		schema.Else = source.Else
	}
	if source.AllOf != nil {
		schema.AllOf = source.AllOf
	}
	if source.AnyOf != nil {
		schema.AnyOf = source.AnyOf
	}
	if source.OneOf != nil {
		schema.OneOf = source.OneOf
	}
	if source.Not != nil {
		schema.Not = source.Not
	}
}

// TypeIs returns true if the Type of a Combined includes the specified type
func (schema *Schema) TypeIs(typeName string) bool {
	if schema.Type != nil {
		// the schema Type is either a string or an array of strings
		if schema.Type.String != nil {
			return (*(schema.Type.String) == typeName)
		} else if schema.Type.StringArray != nil {
			for _, n := range *(schema.Type.StringArray) {
				if n == typeName {
					return true
				}
			}
		}
	}
	return false
}

// ResolveRefs resolves "$ref" elements in a Combined and its children.
// But if a reference refers to an object type, is inside a oneOf, or contains a oneOf,
// the reference is kept and we expect downstream tools to separately model these
// referenced schemas.
func (combined *Combined) ResolveRefs() {
	rootCombined := combined
	count := 1
	for count > 0 {
		count = 0
		combined.applyToCombineds(
			func(combined *Combined, context string) {
				if combined.Boolean != nil {
					return
				}
				schema := combined.Schema
				if schema.Ref != nil {
					resolvedRef, err := rootCombined.resolveJSONPointer(*(schema.Ref))
					if err != nil {
						log.Printf("%+v", err)
					} else if resolvedRef.Boolean != nil {
						// don't substitute for booleans, we'll model the referenced schema with a class
					} else if resolvedRef.Schema.TypeIs("object") {
						// don't substitute for objects, we'll model the referenced schema with a class
					} else if context == "OneOf" {
						// don't substitute for references inside oneOf declarations
					} else if resolvedRef.Schema.OneOf != nil {
						// don't substitute for references that contain oneOf declarations
					} else if resolvedRef.Schema.AdditionalProperties != nil {
						// don't substitute for references that look like objects
					} else {
						schema.Ref = nil
						combined.CopyProperties(resolvedRef)
						count++
					}
				}
			}, "")
	}
}

// resolveJSONPointer resolves JSON pointers.
// This current implementation is very crude and custom for OpenAPI 2.0 schemas.
// It panics for any pointer that it is unable to resolve.
func (combined *Combined) resolveJSONPointer(ref string) (result *Combined, err error) {
	parts := strings.Split(ref, "#")
	if len(parts) == 2 {
		documentName := parts[0] + "#"
		if documentName == "#" && combined.Schema.ID != nil {
			documentName = *(combined.Schema.ID)
		}
		path := parts[1]
		document := schemas[documentName]
		pathParts := strings.Split(path, "/")

		// we currently do a very limited (hard-coded) resolution of certain paths and log errors for missed cases
		if len(pathParts) == 1 {
			return NewCombinedWithSchema(document), nil
		} else if len(pathParts) == 3 {
			switch pathParts[1] {
			case "definitions":
				dictionary := document.Definitions
				for _, pair := range *dictionary {
					if pair.Name == pathParts[2] {
						result = pair.Value
					}
				}
			case "properties":
				dictionary := document.Properties
				for _, pair := range *dictionary {
					if pair.Name == pathParts[2] {
						result = pair.Value
					}
				}
			default:
				break
			}
		}
	}
	if result == nil {
		return nil, fmt.Errorf("unresolved pointer: %+v", ref)
	}
	return result, nil
}

// ResolveAllOfs replaces "allOf" elements by merging their properties into the parent Combined.
func (combined *Combined) ResolveAllOfs() {
	combined.applyToCombineds(
		func(combined *Combined, context string) {
			if combined.Boolean != nil {
				return
			}
			if combined.Schema.AllOf != nil {
				for _, allOf := range *(combined.Schema.AllOf) {
					combined.CopyProperties(allOf)
				}
				combined.Schema.AllOf = nil
			}
		}, "resolveAllOfs")
}

// ResolveAnyOfs replaces all "anyOf" elements with "oneOf".
func (combined *Combined) ResolveAnyOfs() {
	combined.applyToCombineds(
		func(combined *Combined, context string) {
			if combined.Boolean != nil {
				return
			}
			if combined.Schema.AnyOf != nil {
				combined.Schema.OneOf = combined.Schema.AnyOf
				combined.Schema.AnyOf = nil
			}
		}, "resolveAnyOfs")
}

// return a pointer to a copy of a passed-in string
func stringptr(input string) (output *string) {
	return &input
}

// CopyOfficialCombinedProperty copies a named property from the official JSON Combined definition
func (combined *Combined) CopyOfficialCombinedProperty(name string) {
	if combined.Boolean != nil {
		return
	}
	*combined.Schema.Properties = append(*combined.Schema.Properties,
		NewNamedCombined(name, &Combined{
			Schema: &Schema{Ref: stringptr("http://json-schema.org/draft-07/schema#/properties/" + name)}}))
}

// CopyOfficialCombinedProperties copies named properties from the official JSON Combined definition
func (combined *Combined) CopyOfficialCombinedProperties(names []string) {
	for _, name := range names {
		combined.CopyOfficialCombinedProperty(name)
	}
}
