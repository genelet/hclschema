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

//go:generate go run generate-base.go

package jsonschema07

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

// This is a global map of all known Schemas.
// It is initialized when the first Schema is created and inserted.
var schemas map[string]*Schema

// NewBaseSchema builds a schema object from an embedded json representation.
func NewBaseSchema() (*Schema, error) {
	b, err := baseSchemaBytes()
	if err != nil {
		return nil, err
	}
	var node yaml.Node
	err = yaml.Unmarshal(b, &node)
	if err != nil {
		return nil, err
	}
	return NewSchemaFromObject(&node), nil
}

// NewSchemaFromFile reads a schema from a file.
// Currently this assumes that schemas are stored in the source distribution of this project.
func NewSchemaFromFile(filename string) (*Schema, error) {
	file, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var node yaml.Node
	err = yaml.Unmarshal(file, &node)
	if err != nil {
		return nil, err
	}
	return NewSchemaFromObject(&node), nil
}

// NewSchemaFromObject constructs a schema from a parsed JSON object.
// Due to the complexity of the schema representation, this is a
// custom reader and not the standard Go JSON reader (encoding/json).
func NewSchemaFromObject(jsonData *yaml.Node) *Schema {
	switch jsonData.Kind {
	case yaml.DocumentNode:
		return NewSchemaFromObject(jsonData.Content[0])
	case yaml.MappingNode:
		schema := &Schema{}

		for i := 0; i < len(jsonData.Content); i += 2 {
			k := jsonData.Content[i].Value
			v := jsonData.Content[i+1]

			switch k {
			case "$id":
				schema.ID = stringValue(v)
			case "$schema":
				schema.Schema = stringValue(v)
			case "$ref": // uri-reference
				schema.Ref = stringValue(v)
			case "$comment":
				schema.Comment = stringValue(v)
			case "title":
				schema.Title = stringValue(v)
			case "description":
				schema.Description = stringValue(v)

			case "default":
				schema.Default = v
			case "readOnly":
				schema.ReadOnly = boolValue(v)
			case "writeOnly":
				schema.WriteOnly = boolValue(v)
			case "examples":
				schema.Examples = arrayOfCombinedsValue(v)

			case "multipleOf":
				schema.MultipleOf = numberValue(v)
			case "maximum":
				schema.Maximum = numberValue(v)
			case "exclusiveMaximum":
				schema.ExclusiveMaximum = numberValue(v)
			case "minimum":
				schema.Minimum = numberValue(v)
			case "exclusiveMinimum":
				schema.ExclusiveMinimum = numberValue(v)

			case "maxLength":
				schema.MaxLength = intValue(v)
			case "minLength":
				schema.MinLength = intValue(v)
			case "pattern": // regex
				schema.Pattern = stringValue(v)

			case "additionalItems":
				schema.AdditionalItems = combinedFromObject(v)
			case "items":
				schema.Items = schemaOrCombinedArrayValue(v)
			case "maxItems":
				schema.MaxItems = intValue(v)
			case "minItems":
				schema.MinItems = intValue(v)
			case "uniqueItems":
				schema.UniqueItems = boolValue(v)

			case "contains":
				schema.Contains = combinedFromObject(v)
			case "maxProperties":
				schema.MaxProperties = intValue(v)
			case "minProperties":
				schema.MinProperties = intValue(v)
			case "required":
				schema.Required = arrayOfStringsValue(v)
			case "additionalProperties":
				schema.AdditionalProperties = combinedFromObject(v)
			case "definitions":
				schema.Definitions = mapOfCombinedsValue(v)
			case "properties":
				schema.Properties = mapOfCombinedsValue(v)
			case "patternProperties":
				schema.PatternProperties = mapOfCombinedsValue(v)
			case "dependencies":
				schema.Dependencies = mapOfCombinedsOrStringArraysValue(v)
			case "propertyNames":
				schema.PropertyNames = combinedFromObject(v)

			case "const":
				schema.Const = v
			case "enum":
				schema.Enumeration = arrayOfEnumValuesValue(v)
			case "type":
				schema.Type = stringOrStringArrayValue(v)
			case "format":
				schema.Format = stringValue(v)
			case "contentMediaType":
				schema.ContentMediaType = stringValue(v)
			case "contentEncoding":
				schema.ContentEncoding = stringValue(v)

			case "if":
				schema.If = combinedFromObject(v)
			case "then":
				schema.Then = combinedFromObject(v)
			case "else":
				schema.Else = combinedFromObject(v)
			case "allOf":
				schema.AllOf = arrayOfCombinedsValue(v)
			case "anyOf":
				schema.AnyOf = arrayOfCombinedsValue(v)
			case "oneOf":
				schema.OneOf = arrayOfCombinedsValue(v)
			case "not":
				schema.Not = combinedFromObject(v)

			default:
				fmt.Printf("UNSUPPORTED (%s)\n", k)
				//fmt.Printf("%s=>%#v", jsonData.Tag, v.Value)
				//panic(nil)
			}
		}

		// insert schema in global map
		if schema.ID != nil {
			if schemas == nil {
				schemas = make(map[string]*Schema, 0)
			}
			schemas[*(schema.ID)] = schema
		}
		return schema

	default:
		fmt.Printf("schemaValue: unexpected node %+v\n", jsonData)
	}

	return nil
}

//
// BUILDERS
// The following methods build elements of Combineds from interface{} values.
// Each returns nil if it is unable to build the desired element.
//

// combinedFromObject constructs a schema or a boolean value from a parsed JSON object.
func combinedFromObject(v *yaml.Node) *Combined {
	schemaOrBoolean := &Combined{}
	switch v.Kind {
	case yaml.ScalarNode:
		v2, _ := strconv.ParseBool(v.Value)
		schemaOrBoolean.Boolean = &v2
	case yaml.MappingNode:
		schemaOrBoolean.Schema = NewSchemaFromObject(v)
	default:
		fmt.Printf("NewSchemaFromObject: unexpected node %+v\n", v)
	}
	return schemaOrBoolean
}

// Gets the string value of an interface{} value if possible.
func stringValue(v *yaml.Node) *string {
	switch v.Kind {
	case yaml.ScalarNode:
		str := strings.Replace(v.Value, "\n", "\\n", -1)
		return &str
	default:
		fmt.Printf("stringValue: unexpected node %+v\n", v)
	}
	return nil
}

// Gets the numeric value of an interface{} value if possible.
func numberValue(v *yaml.Node) *IntegerOrFloat {
	number := &IntegerOrFloat{}
	switch v.Kind {
	case yaml.ScalarNode:
		switch v.Tag {
		case "!!float":
			v2, _ := strconv.ParseFloat(v.Value, 64)
			number.Float = &v2
			return number
		case "!!int":
			v2, _ := strconv.ParseInt(v.Value, 10, 64)
			number.Integer = &v2
			return number
		default:
			fmt.Printf("stringValue: unexpected node %+v\n", v)
		}
	default:
		fmt.Printf("stringValue: unexpected node %+v\n", v)
	}
	return nil
}

// Gets the integer value of an interface{} value if possible.
func intValue(v *yaml.Node) *int64 {
	switch v.Kind {
	case yaml.ScalarNode:
		switch v.Tag {
		case "!!float":
			v2, _ := strconv.ParseFloat(v.Value, 64)
			v3 := int64(v2)
			return &v3
		case "!!int":
			v2, _ := strconv.ParseInt(v.Value, 10, 64)
			return &v2
		default:
			fmt.Printf("intValue: unexpected node %+v\n", v)
		}
	default:
		fmt.Printf("intValue: unexpected node %+v\n", v)
	}
	return nil
}

// Gets the bool value of an interface{} value if possible.
func boolValue(v *yaml.Node) *bool {
	switch v.Kind {
	case yaml.ScalarNode:
		switch v.Tag {
		case "!!bool":
			v2, _ := strconv.ParseBool(v.Value)
			return &v2
		default:
			fmt.Printf("boolValue: unexpected node %+v\n", v)
		}
	default:
		fmt.Printf("boolValue: unexpected node %+v\n", v)
	}
	return nil
}

// Gets array of interfaces from an interface{} value if possible.
func arrayValue(v *yaml.Node) *[]interface{} {
	switch v.Kind {
	case yaml.ScalarNode:
		a := make([]interface{}, 0)
		a = append(a, v.Value)
		return &a
	case yaml.SequenceNode:
		a := make([]interface{}, 0)
		for _, v2 := range v.Content {
			switch v2.Kind {
			case yaml.ScalarNode:
				a = append(a, v2.Value)
			default:
				fmt.Printf("arrayValue: unexpected node %+v\n", v2)
			}
		}
		return &a
	default:
		fmt.Printf("arrayValue: unexpected node %+v\n", v)
	}
	return nil
}

// Gets a map of Combineds from an interface{} value if possible.
func mapOfCombinedsValue(v *yaml.Node) []*NamedCombined {
	switch v.Kind {
	case yaml.MappingNode:
		m := make([]*NamedCombined, 0)
		for i := 0; i < len(v.Content); i += 2 {
			k2 := v.Content[i].Value
			v2 := v.Content[i+1]
			pair := &NamedCombined{Name: k2, Value: combinedFromObject(v2)}
			m = append(m, pair)
		}
		return m
	default:
		fmt.Printf("mapOfCombinedsValue: unexpected node %+v\n", v)
	}
	return nil
}

// Gets an array of Combineds from an interface{} value if possible.
func arrayOfCombinedsValue(v *yaml.Node) []*Combined {
	switch v.Kind {
	case yaml.SequenceNode:
		m := make([]*Combined, 0)
		for _, v2 := range v.Content {
			switch v2.Kind {
			case yaml.MappingNode:
				s := combinedFromObject(v2)
				m = append(m, s)
			default:
				fmt.Printf("arrayOfCombinedOrBooleansValue: unexpected node %+v\n", v2)
			}
		}
		return m
	case yaml.MappingNode:
		m := make([]*Combined, 0)
		s := combinedFromObject(v)
		m = append(m, s)
		return m
	default:
		fmt.Printf("arrayOfCombinedOrBooleansValue: unexpected node %+v\n", v)
	}
	return nil
}

// Gets a Combined or an array of Combineds from an interface{} value if possible.
func schemaOrCombinedArrayValue(v *yaml.Node) *CombinedOrCombinedArray {
	switch v.Kind {
	case yaml.SequenceNode:
		m := make([]*Combined, 0)
		for _, v2 := range v.Content {
			switch v2.Kind {
			case yaml.MappingNode:
				s := combinedFromObject(v2)
				m = append(m, s)
			default:
				fmt.Printf("schemaOrCombinedArrayValue: unexpected node %+v\n", v2)
			}
		}
		return &CombinedOrCombinedArray{CombinedArray: &m}
	case yaml.MappingNode:
		s := combinedFromObject(v)
		return &CombinedOrCombinedArray{Combined: s}
	default:
		fmt.Printf("schemaOrCombinedArrayValue: unexpected node %+v\n", v)
	}
	return nil
}

// Gets an array of strings from an interface{} value if possible.
func arrayOfStringsValue(v *yaml.Node) []string {
	switch v.Kind {
	case yaml.ScalarNode:
		return []string{v.Value}
	case yaml.SequenceNode:
		a := make([]string, 0)
		for _, v2 := range v.Content {
			switch v2.Kind {
			case yaml.ScalarNode:
				a = append(a, v2.Value)
			default:
				fmt.Printf("arrayOfStringsValue: unexpected node %+v\n", v2)
			}
		}
		return a
	default:
		fmt.Printf("arrayOfStringsValue: unexpected node %+v\n", v)
	}
	return nil
}

// Gets a string or an array of strings from an interface{} value if possible.
func stringOrStringArrayValue(v *yaml.Node) *StringOrStringArray {
	switch v.Kind {
	case yaml.ScalarNode:
		s := &StringOrStringArray{}
		s.String = &v.Value
		return s
	case yaml.SequenceNode:
		a := make([]string, 0)
		for _, v2 := range v.Content {
			switch v2.Kind {
			case yaml.ScalarNode:
				a = append(a, v2.Value)
			default:
				fmt.Printf("arrayOfStringsValue: unexpected node %+v\n", v2)
			}
		}
		s := &StringOrStringArray{}
		s.StringArray = &a
		return s
	default:
		fmt.Printf("arrayOfStringsValue: unexpected node %+v\n", v)
	}
	return nil
}

// Gets an array of enum values from an interface{} value if possible.
func arrayOfEnumValuesValue(v *yaml.Node) []SchemaEnumValue {
	a := make([]SchemaEnumValue, 0)
	switch v.Kind {
	case yaml.SequenceNode:
		for _, v2 := range v.Content {
			switch v2.Kind {
			case yaml.ScalarNode:
				switch v2.Tag {
				case "!!str":
					a = append(a, SchemaEnumValue{String: &v2.Value})
				case "!!bool":
					v3, _ := strconv.ParseBool(v2.Value)
					a = append(a, SchemaEnumValue{Bool: &v3})
				default:
					fmt.Printf("arrayOfEnumValuesValue: unexpected type %s\n", v2.Tag)
				}
			default:
				fmt.Printf("arrayOfEnumValuesValue: unexpected node %+v\n", v2)
			}
		}
	default:
		fmt.Printf("arrayOfEnumValuesValue: unexpected node %+v\n", v)
	}
	return a
}

// Gets a map of schemas or string arrays from an interface{} value if possible.
func mapOfCombinedsOrStringArraysValue(v *yaml.Node) []*NamedCombinedOrStringArray {
	m := make([]*NamedCombinedOrStringArray, 0)
	switch v.Kind {
	case yaml.MappingNode:
		for i := 0; i < len(v.Content); i += 2 {
			k2 := v.Content[i].Value
			v2 := v.Content[i+1]
			switch v2.Kind {
			case yaml.SequenceNode:
				a := make([]string, 0)
				for _, v3 := range v2.Content {
					switch v3.Kind {
					case yaml.ScalarNode:
						a = append(a, v3.Value)
					default:
						fmt.Printf("mapOfCombinedsOrStringArraysValue: unexpected node %+v\n", v3)
					}
				}
				s := &CombinedOrStringArray{}
				s.StringArray = &a
				pair := &NamedCombinedOrStringArray{Name: k2, Value: s}
				m = append(m, pair)
			default:
				fmt.Printf("mapOfCombinedsOrStringArraysValue: unexpected node %+v\n", v2)
			}
		}
	default:
		fmt.Printf("mapOfCombinedsOrStringArraysValue: unexpected node %+v\n", v)
	}
	return m
}
