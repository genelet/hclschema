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
	"strings"
)

//
// DISPLAY
// The following methods display Schemas.
//

// Description returns a string representation of a string or string array.
func (s *StringOrStringArray) Description() string {
	if s.String != nil {
		return *s.String
	}
	if s.StringArray != nil {
		return strings.Join(*s.StringArray, ", ")
	}
	return ""
}

// Returns a string representation of a Schema.
func (schema *Schema) String() string {
	return schema.describeSchema("")
}

// Helper: Returns a string representation of a Schema indented by a specified string.
func (schema *Schema) describeSchema(indent string) string {
	if schema.Boolean != nil {
		return indent + fmt.Sprintf("%+v\n", *(schema.Boolean))
	}

	absolute := schema.Absolute
	result := ""
	if absolute.Schema != nil {
		result += indent + "$schema: " + *(absolute.Schema) + "\n"
	}
	if absolute.ID != nil {
		switch strings.TrimSuffix(*absolute.Schema, "#") {
		case "http://json-schema.org/draft-07/schema#":
			fallthrough
		case "#":
			fallthrough
		case "":
			result += indent + "id: " + *(absolute.ID) + "\n"
		default:
			result += indent + "$id: " + *(absolute.ID) + "\n"
		}
	}
	if absolute.Ref != nil {
		result += indent + "$ref: " + *(absolute.Ref) + "\n"
	}
	if absolute.Comment != nil {
		result += indent + "$comment: " + *(absolute.Comment) + "\n"
	}
	if absolute.Title != nil {
		result += indent + "title: " + *(absolute.Title) + "\n"
	}
	if absolute.Description != nil {
		result += indent + "description: " + *(absolute.Description) + "\n"
	}

	if absolute.Default != nil {
		result += indent + "default:\n"
		result += indent + fmt.Sprintf("  %+v\n", *(absolute.Default))
	}
	if absolute.ReadOnly != nil && *absolute.ReadOnly {
		result += indent + fmt.Sprintf("readOnly: %+v\n", *(absolute.ReadOnly))
	}
	if absolute.WriteOnly != nil && *absolute.WriteOnly {
		result += indent + fmt.Sprintf("writeOnly: %+v\n", *(absolute.WriteOnly))
	}
	if absolute.Examples != nil {
		result += indent + "examples:\n"
		for _, example := range *(absolute.Examples) {
			result += indent + "  " + fmt.Sprintf("%+v\n", *(example.Absolute))
		}
	}

	if absolute.MultipleOf != nil {
		result += indent + fmt.Sprintf("multipleOf: %+v\n", *(absolute.MultipleOf))
	}
	if absolute.Maximum != nil {
		result += indent + fmt.Sprintf("maximum: %+v\n", *(absolute.Maximum))
	}
	if absolute.ExclusiveMaximum != nil {
		result += indent + fmt.Sprintf("exclusiveMaximum: %+v\n", *(absolute.ExclusiveMaximum))
	}
	if absolute.Minimum != nil {
		result += indent + fmt.Sprintf("minimum: %+v\n", *(absolute.Minimum))
	}
	if absolute.ExclusiveMinimum != nil {
		result += indent + fmt.Sprintf("exclusiveMinimum: %+v\n", *(absolute.ExclusiveMinimum))
	}

	if absolute.MaxLength != nil {
		result += indent + fmt.Sprintf("maxLength: %+v\n", *(absolute.MaxLength))
	}
	if absolute.MinLength != nil {
		result += indent + fmt.Sprintf("minLength: %+v\n", *(absolute.MinLength))
	}
	if absolute.Pattern != nil {
		result += indent + fmt.Sprintf("pattern: %+v\n", *(absolute.Pattern))
	}

	if absolute.AdditionalItems != nil {
		s := absolute.AdditionalItems
		if s != nil {
			result += indent + "additionalItems:\n"
			result += s.describeSchema(indent + "  ")
		} else {
			b := *(absolute.AdditionalItems.Boolean)
			result += indent + fmt.Sprintf("additionalItems: %+v\n", b)
		}
	}
	if absolute.Items != nil {
		result += indent + "items:\n"
		items := absolute.Items
		if items.SchemaArray != nil {
			for i, s := range *(items.SchemaArray) {
				result += indent + "  " + fmt.Sprintf("%d", i) + ":\n"
				result += s.describeSchema(indent + "  " + "  ")
			}
		} else if items.Schema != nil {
			result += items.Schema.describeSchema(indent + "  " + "  ")
		}
	}
	if absolute.MaxItems != nil {
		result += indent + fmt.Sprintf("maxItems: %+v\n", *(absolute.MaxItems))
	}
	if absolute.MinItems != nil {
		result += indent + fmt.Sprintf("minItems: %+v\n", *(absolute.MinItems))
	}
	if absolute.UniqueItems != nil {
		result += indent + fmt.Sprintf("uniqueItems: %+v\n", *(absolute.UniqueItems))
	}

	if absolute.Contains != nil {
		result += indent + "contains:\n"
		result += absolute.Contains.describeSchema(indent + "  ")
	}
	if absolute.MaxProperties != nil {
		result += indent + fmt.Sprintf("maxProperties: %+v\n", *(absolute.MaxProperties))
	}
	if absolute.MinProperties != nil {
		result += indent + fmt.Sprintf("minProperties: %+v\n", *(absolute.MinProperties))
	}
	if absolute.Required != nil {
		result += indent + fmt.Sprintf("required: %+v\n", *(absolute.Required))
	}
	if absolute.AdditionalProperties != nil {
		s := absolute.AdditionalProperties
		if s != nil {
			result += indent + "additionalProperties:\n"
			result += s.describeSchema(indent + "  ")
		} else {
			b := *(absolute.AdditionalProperties.Boolean)
			result += indent + fmt.Sprintf("additionalProperties: %+v\n", b)
		}
	}
	if absolute.Definitions != nil {
		result += indent + "definitions:\n"
		for _, pair := range *(absolute.Definitions) {
			name := pair.Name
			s := pair.Value
			result += indent + "  " + name + ":\n"
			result += s.describeSchema(indent + "  " + "  ")
		}
	}
	if absolute.Properties != nil {
		result += indent + "properties:\n"
		for _, pair := range *(absolute.Properties) {
			name := pair.Name
			s := pair.Value
			result += indent + "  " + name + ":\n"
			result += s.describeSchema(indent + "  " + "  ")
		}
	}
	if absolute.PatternProperties != nil {
		result += indent + "patternProperties:\n"
		for _, pair := range *(absolute.PatternProperties) {
			name := pair.Name
			s := pair.Value
			result += indent + "  " + name + ":\n"
			result += s.describeSchema(indent + "  " + "  ")
		}
	}
	if absolute.Dependencies != nil {
		result += indent + "dependencies:\n"
		for _, pair := range *(absolute.Dependencies) {
			name := pair.Name
			schemaOrStringArray := pair.Value
			s := schemaOrStringArray.Schema
			if s != nil {
				result += indent + "  " + name + ":\n"
				result += s.describeSchema(indent + "  " + "  ")
			} else {
				a := schemaOrStringArray.StringArray
				if a != nil {
					result += indent + "  " + name + ":\n"
					for _, s2 := range *a {
						result += indent + "  " + "  " + s2 + "\n"
					}
				}
			}

		}
	}
	if absolute.PropertyNames != nil {
		result += indent + "propertyNames:\n"
		result += absolute.PropertyNames.describeSchema(indent + "  ")
	}

	if absolute.Const != nil {
		result += indent + "const:\n"
		result += indent + fmt.Sprintf("  %+v\n", *(absolute.Const))
	}
	if absolute.Enumeration != nil {
		result += indent + "enumeration:\n"
		for _, value := range *(absolute.Enumeration) {
			if value.String != nil {
				result += indent + "  " + fmt.Sprintf("%+v\n", *value.String)
			} else {
				result += indent + "  " + fmt.Sprintf("%+v\n", *value.Bool)
			}
		}
	}
	if absolute.Type != nil {
		result += indent + fmt.Sprintf("type: %+v\n", absolute.Type.Description())
	}
	if absolute.Format != nil {
		result += indent + "format: " + *(absolute.Format) + "\n"
	}
	if absolute.ContentMediaType != nil {
		result += indent + "contentMediaType: " + *(absolute.ContentMediaType) + "\n"
	}
	if absolute.ContentEncoding != nil {
		result += indent + "contentEncoding: " + *(absolute.ContentEncoding) + "\n"
	}

	if absolute.If != nil {
		result += indent + "if:\n"
		result += absolute.If.describeSchema(indent + "  ")
	}
	if absolute.Then != nil {
		result += indent + "then:\n"
		result += absolute.Then.describeSchema(indent + "  ")
	}
	if absolute.Else != nil {
		result += indent + "else:\n"
		result += absolute.Else.describeSchema(indent + "  ")
	}
	if absolute.AllOf != nil {
		result += indent + "allOf:\n"
		for _, s := range *(absolute.AllOf) {
			result += s.describeSchema(indent + "  ")
			result += indent + "-\n"
		}
	}
	if absolute.AnyOf != nil {
		result += indent + "anyOf:\n"
		for _, s := range *(absolute.AnyOf) {
			result += s.describeSchema(indent + "  ")
			result += indent + "-\n"
		}
	}
	if absolute.OneOf != nil {
		result += indent + "oneOf:\n"
		for _, s := range *(absolute.OneOf) {
			result += s.describeSchema(indent + "  ")
			result += indent + "-\n"
		}
	}
	if absolute.Not != nil {
		result += indent + "not:\n"
		result += absolute.Not.describeSchema(indent + "  ")
	}

	return result
}
