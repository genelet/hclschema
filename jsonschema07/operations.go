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
// The following methods perform operations on Schemas.
//

// IsEmpty returns true if no members of the Schema are specified.
func (absolute *Absolute) IsEmpty() bool {
	return (absolute.Schema == nil) &&
		(absolute.ID == nil) &&
		(absolute.MultipleOf == nil) &&
		(absolute.Maximum == nil) &&
		(absolute.ExclusiveMaximum == nil) &&
		(absolute.Minimum == nil) &&
		(absolute.ExclusiveMinimum == nil) &&
		(absolute.MaxLength == nil) &&
		(absolute.MinLength == nil) &&
		(absolute.Pattern == nil) &&
		(absolute.AdditionalItems == nil) &&
		(absolute.Items == nil) &&
		(absolute.MaxItems == nil) &&
		(absolute.MinItems == nil) &&
		(absolute.UniqueItems == nil) &&
		(absolute.MaxProperties == nil) &&
		(absolute.MinProperties == nil) &&
		(absolute.Required == nil) &&
		(absolute.AdditionalProperties == nil) &&
		(absolute.Properties == nil) &&
		(absolute.PatternProperties == nil) &&
		(absolute.Dependencies == nil) &&
		(absolute.Enumeration == nil) &&
		(absolute.Type == nil) &&
		(absolute.AllOf == nil) &&
		(absolute.AnyOf == nil) &&
		(absolute.OneOf == nil) &&
		(absolute.Not == nil) &&
		(absolute.Definitions == nil) &&
		(absolute.Title == nil) &&
		(absolute.Description == nil) &&
		(absolute.Default == nil) &&
		(absolute.Format == nil) &&
		(absolute.Ref == nil)
}

// IsEqual returns true if two schemas are equal.
func (schema *Schema) IsEqual(schema2 *Schema) bool {
	return schema.String() == schema2.String()
}

// SchemaOperation represents a function that can be applied to a Schema.
type SchemaOperation func(schema *Schema, context string)

// Applies a specified function to a Schema and all of the Schemas that it contains.
func (schema *Schema) applyToSchemas(operation SchemaOperation, context string) {
	if schema.Boolean != nil {
		return
	}

	absolute := schema.Absolute
	if absolute.AdditionalItems != nil {
		s := absolute.AdditionalItems
		if s != nil {
			s.applyToSchemas(operation, "AdditionalItems")
		}
	}

	if absolute.Items != nil {
		if absolute.Items.SchemaArray != nil {
			for _, s := range *(absolute.Items.SchemaArray) {
				s.applyToSchemas(operation, "Items.SchemaArray")
			}
		} else if absolute.Items.Schema != nil {
			absolute.Items.Schema.applyToSchemas(operation, "Items.Schema")
		}
	}

	if absolute.AdditionalProperties != nil {
		s := absolute.AdditionalProperties
		if s != nil {
			s.applyToSchemas(operation, "AdditionalProperties")
		}
	}

	if absolute.Properties != nil {
		for _, pair := range *(absolute.Properties) {
			s := pair.Value
			s.applyToSchemas(operation, "Properties")
		}
	}
	if absolute.PatternProperties != nil {
		for _, pair := range *(absolute.PatternProperties) {
			s := pair.Value
			s.applyToSchemas(operation, "PatternProperties")
		}
	}

	if absolute.Dependencies != nil {
		for _, pair := range *(absolute.Dependencies) {
			schemaOrStringArray := pair.Value
			s := schemaOrStringArray.Schema
			if s != nil {
				s.applyToSchemas(operation, "Dependencies")
			}
		}
	}

	if absolute.AllOf != nil {
		for _, s := range *(absolute.AllOf) {
			s.applyToSchemas(operation, "AllOf")
		}
	}
	if absolute.AnyOf != nil {
		for _, s := range *(absolute.AnyOf) {
			s.applyToSchemas(operation, "AnyOf")
		}
	}
	if absolute.OneOf != nil {
		for _, s := range *(absolute.OneOf) {
			s.applyToSchemas(operation, "OneOf")
		}
	}
	if absolute.Not != nil {
		absolute.Not.applyToSchemas(operation, "Not")
	}

	if absolute.Definitions != nil {
		for _, pair := range *(absolute.Definitions) {
			s := pair.Value
			s.applyToSchemas(operation, "Definitions")
		}
	}

	operation(schema, context)
}

// CopyProperties copies all non-nil properties from the source Schema to the schema Schema.
func (schema *Schema) CopyProperties(source *Schema) {
	if schema.Boolean != nil {
		return
	}
	schema.Absolute.CopyProperties(source.Absolute)
}

func (absolute *Absolute) CopyProperties(source *Absolute) {
	if source.Schema != nil {
		absolute.Schema = source.Schema
	}
	if source.ID != nil {
		absolute.ID = source.ID
	}
	if source.MultipleOf != nil {
		absolute.MultipleOf = source.MultipleOf
	}
	if source.Maximum != nil {
		absolute.Maximum = source.Maximum
	}
	if source.ExclusiveMaximum != nil {
		absolute.ExclusiveMaximum = source.ExclusiveMaximum
	}
	if source.Minimum != nil {
		absolute.Minimum = source.Minimum
	}
	if source.ExclusiveMinimum != nil {
		absolute.ExclusiveMinimum = source.ExclusiveMinimum
	}
	if source.MaxLength != nil {
		absolute.MaxLength = source.MaxLength
	}
	if source.MinLength != nil {
		absolute.MinLength = source.MinLength
	}
	if source.Pattern != nil {
		absolute.Pattern = source.Pattern
	}
	if source.AdditionalItems != nil {
		absolute.AdditionalItems = source.AdditionalItems
	}
	if source.Items != nil {
		absolute.Items = source.Items
	}
	if source.MaxItems != nil {
		absolute.MaxItems = source.MaxItems
	}
	if source.MinItems != nil {
		absolute.MinItems = source.MinItems
	}
	if source.UniqueItems != nil {
		absolute.UniqueItems = source.UniqueItems
	}
	if source.MaxProperties != nil {
		absolute.MaxProperties = source.MaxProperties
	}
	if source.MinProperties != nil {
		absolute.MinProperties = source.MinProperties
	}
	if source.Required != nil {
		absolute.Required = source.Required
	}
	if source.AdditionalProperties != nil {
		absolute.AdditionalProperties = source.AdditionalProperties
	}
	if source.Properties != nil {
		absolute.Properties = source.Properties
	}
	if source.PatternProperties != nil {
		absolute.PatternProperties = source.PatternProperties
	}
	if source.Dependencies != nil {
		absolute.Dependencies = source.Dependencies
	}
	if source.Enumeration != nil {
		absolute.Enumeration = source.Enumeration
	}
	if source.Type != nil {
		absolute.Type = source.Type
	}
	if source.AllOf != nil {
		absolute.AllOf = source.AllOf
	}
	if source.AnyOf != nil {
		absolute.AnyOf = source.AnyOf
	}
	if source.OneOf != nil {
		absolute.OneOf = source.OneOf
	}
	if source.Not != nil {
		absolute.Not = source.Not
	}
	if source.Definitions != nil {
		absolute.Definitions = source.Definitions
	}
	if source.Title != nil {
		absolute.Title = source.Title
	}
	if source.Description != nil {
		absolute.Description = source.Description
	}
	if source.Default != nil {
		absolute.Default = source.Default
	}
	if source.Format != nil {
		absolute.Format = source.Format
	}
	if source.Ref != nil {
		absolute.Ref = source.Ref
	}
}

// TypeIs returns true if the Type of a Schema includes the specified type
func (absolute *Absolute) TypeIs(typeName string) bool {
	if absolute.Type != nil {
		// the absolute Type is either a string or an array of strings
		if absolute.Type.String != nil {
			return (*(absolute.Type.String) == typeName)
		} else if absolute.Type.StringArray != nil {
			for _, n := range *(absolute.Type.StringArray) {
				if n == typeName {
					return true
				}
			}
		}
	}
	return false
}

// ResolveRefs resolves "$ref" elements in a Schema and its children.
// But if a reference refers to an object type, is inside a oneOf, or contains a oneOf,
// the reference is kept and we expect downstream tools to separately model these
// referenced schemas.
func (schema *Schema) ResolveRefs() {
	rootSchema := schema
	count := 1
	for count > 0 {
		count = 0
		schema.applyToSchemas(
			func(schema *Schema, context string) {
				if schema.Boolean != nil {
					return
				}
				absolute := schema.Absolute
				if absolute.Ref != nil {
					resolvedRef, err := rootSchema.resolveJSONPointer(*(absolute.Ref))
					if err != nil {
						log.Printf("%+v", err)
					} else if resolvedRef.Boolean != nil {
						// don't substitute for booleans, we'll model the referenced schema with a class
					} else if resolvedRef.Absolute.TypeIs("object") {
						// don't substitute for objects, we'll model the referenced schema with a class
					} else if context == "OneOf" {
						// don't substitute for references inside oneOf declarations
					} else if resolvedRef.Absolute.OneOf != nil {
						// don't substitute for references that contain oneOf declarations
					} else if resolvedRef.Absolute.AdditionalProperties != nil {
						// don't substitute for references that look like objects
					} else {
						schema.Absolute.Ref = nil
						schema.CopyProperties(resolvedRef)
						count++
					}
				}
			}, "")
	}
}

// resolveJSONPointer resolves JSON pointers.
// This current implementation is very crude and custom for OpenAPI 2.0 schemas.
// It panics for any pointer that it is unable to resolve.
func (schema *Schema) resolveJSONPointer(ref string) (result *Schema, err error) {
	parts := strings.Split(ref, "#")
	if len(parts) == 2 {
		documentName := parts[0] + "#"
		if documentName == "#" && schema.Absolute.ID != nil {
			documentName = *(schema.Absolute.ID)
		}
		path := parts[1]
		document := schemas[documentName]
		pathParts := strings.Split(path, "/")

		// we currently do a very limited (hard-coded) resolution of certain paths and log errors for missed cases
		if len(pathParts) == 1 {
			return document, nil
		} else if len(pathParts) == 3 {
			switch pathParts[1] {
			case "definitions":
				dictionary := document.Absolute.Definitions
				for _, pair := range *dictionary {
					if pair.Name == pathParts[2] {
						result = pair.Value
					}
				}
			case "properties":
				dictionary := document.Absolute.Properties
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

// ResolveAllOfs replaces "allOf" elements by merging their properties into the parent Schema.
func (schema *Schema) ResolveAllOfs() {
	schema.applyToSchemas(
		func(schema *Schema, context string) {
			if schema.Boolean != nil {
				return
			}
			if schema.Absolute.AllOf != nil {
				for _, allOf := range *(schema.Absolute.AllOf) {
					schema.CopyProperties(allOf)
				}
				schema.Absolute.AllOf = nil
			}
		}, "resolveAllOfs")
}

// ResolveAnyOfs replaces all "anyOf" elements with "oneOf".
func (schema *Schema) ResolveAnyOfs() {
	schema.applyToSchemas(
		func(schema *Schema, context string) {
			if schema.Boolean != nil {
				return
			}
			if schema.Absolute.AnyOf != nil {
				schema.Absolute.OneOf = schema.Absolute.AnyOf
				schema.Absolute.AnyOf = nil
			}
		}, "resolveAnyOfs")
}

// return a pointer to a copy of a passed-in string
func stringptr(input string) (output *string) {
	return &input
}

// CopyOfficialSchemaProperty copies a named property from the official JSON Schema definition
func (schema *Schema) CopyOfficialSchemaProperty(name string) {
	if schema.Boolean != nil {
		return
	}
	*schema.Absolute.Properties = append(*schema.Absolute.Properties,
		NewNamedSchema(name, &Schema{Absolute: &Absolute{Ref: stringptr("http://json-schema.org/draft-04/schema#/properties/" + name)}}))
}

// CopyOfficialSchemaProperties copies named properties from the official JSON Schema definition
func (schema *Schema) CopyOfficialSchemaProperties(names []string) {
	for _, name := range names {
		schema.CopyOfficialSchemaProperty(name)
	}
}
