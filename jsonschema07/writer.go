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

	"gopkg.in/yaml.v3"
)

const indentation = "  "

func renderMappingNode(node *yaml.Node, indent string) (result string) {
	result = "{\n"
	innerIndent := indent + indentation
	for i := 0; i < len(node.Content); i += 2 {
		// first print the key
		key := node.Content[i].Value
		result += fmt.Sprintf("%s\"%+v\": ", innerIndent, key)
		// then the value
		value := node.Content[i+1]
		switch value.Kind {
		case yaml.ScalarNode:
			if value.Tag == "!!bool" {
				result += value.Value
			} else {
				result += "\"" + value.Value + "\""
			}
		case yaml.MappingNode:
			result += renderMappingNode(value, innerIndent)
		case yaml.SequenceNode:
			result += renderSequenceNode(value, innerIndent)
		default:
			result += fmt.Sprintf("???MapItem(Key:%+v, Value:%T)", value, value)
		}
		if i < len(node.Content)-2 {
			result += ","
		}
		result += "\n"
	}

	result += indent + "}"
	return result
}

func renderSequenceNode(node *yaml.Node, indent string) (result string) {
	result = "[\n"
	innerIndent := indent + indentation
	for i := 0; i < len(node.Content); i++ {
		item := node.Content[i]
		switch item.Kind {
		case yaml.ScalarNode:
			if item.Tag == "!!bool" {
				result += innerIndent + item.Value
			} else {
				result += innerIndent + "\"" + item.Value + "\""
			}
		case yaml.MappingNode:
			result += innerIndent + renderMappingNode(item, innerIndent) + ""
		default:
			result += innerIndent + fmt.Sprintf("???ArrayItem(%+v)", item)
		}
		if i < len(node.Content)-1 {
			result += ","
		}
		result += "\n"
	}
	result += indent + "]"
	return result
}

func renderStringArray(array []string, indent string) (result string) {
	result = "[\n"
	innerIndent := indent + indentation
	for i, item := range array {
		result += innerIndent + "\"" + item + "\""
		if i < len(array)-1 {
			result += ","
		}
		result += "\n"
	}
	result += indent + "]"
	return result
}

// Render renders a yaml.Node as JSON
func Render(node *yaml.Node) string {
	if node.Kind == yaml.DocumentNode {
		if len(node.Content) == 1 {
			return Render(node.Content[0])
		}
	} else if node.Kind == yaml.MappingNode {
		return renderMappingNode(node, "") + "\n"
	} else if node.Kind == yaml.SequenceNode {
		return renderSequenceNode(node, "") + "\n"
	}
	return ""
}

func (object *SchemaNumber) nodeValue() *yaml.Node {
	if object.Integer != nil {
		return nodeForInt64(*object.Integer)
	} else if object.Float != nil {
		return nodeForFloat64(*object.Float)
	} else {
		return nil
	}
}

func (object *Schema) nodeValue() *yaml.Node {
	if object.Absolute != nil {
		return object.Absolute.nodeValue()
	} else if object.Boolean != nil {
		return nodeForBoolean(*object.Boolean)
	} else {
		return nil
	}
}

func nodeForStringArray(array []string) *yaml.Node {
	content := make([]*yaml.Node, 0)
	for _, item := range array {
		content = append(content, nodeForString(item))
	}
	return nodeForSequence(content)
}

func nodeForSchemaArray(array []*Schema) *yaml.Node {
	content := make([]*yaml.Node, 0)
	for _, item := range array {
		content = append(content, item.nodeValue())
	}
	return nodeForSequence(content)
}

func (object *StringOrStringArray) nodeValue() *yaml.Node {
	if object.String != nil {
		return nodeForString(*object.String)
	} else if object.StringArray != nil {
		return nodeForStringArray(*(object.StringArray))
	} else {
		return nil
	}
}

func (object *SchemaOrStringArray) nodeValue() *yaml.Node {
	if object.Schema != nil {
		return object.Schema.nodeValue()
	} else if object.StringArray != nil {
		return nodeForStringArray(*(object.StringArray))
	} else {
		return nil
	}
}

func (object *SchemaOrSchemaArray) nodeValue() *yaml.Node {
	if object.Schema != nil {
		return object.Schema.nodeValue()
	} else if object.SchemaArray != nil {
		return nodeForSchemaArray(*(object.SchemaArray))
	} else {
		return nil
	}
}

func (object *SchemaEnumValue) nodeValue() *yaml.Node {
	if object.String != nil {
		return nodeForString(*object.String)
	} else if object.Bool != nil {
		return nodeForBoolean(*object.Bool)
	} else {
		return nil
	}
}

func nodeForNamedSchemaArray(array *[]*NamedSchema) *yaml.Node {
	content := make([]*yaml.Node, 0)
	for _, pair := range *(array) {
		content = appendPair(content, pair.Name, pair.Value.nodeValue())
	}
	return nodeForMapping(content)
}

func nodeForNamedSchemaOrStringArray(array *[]*NamedSchemaOrStringArray) *yaml.Node {
	content := make([]*yaml.Node, 0)
	for _, pair := range *(array) {
		content = appendPair(content, pair.Name, pair.Value.nodeValue())
	}
	return nodeForMapping(content)
}

func nodeForSchemaEnumArray(array *[]SchemaEnumValue) *yaml.Node {
	content := make([]*yaml.Node, 0)
	for _, item := range *array {
		content = append(content, item.nodeValue())
	}
	return nodeForSequence(content)
}

func nodeForMapping(content []*yaml.Node) *yaml.Node {
	return &yaml.Node{
		Kind:    yaml.MappingNode,
		Content: content,
	}
}

func nodeForSequence(content []*yaml.Node) *yaml.Node {
	return &yaml.Node{
		Kind:    yaml.SequenceNode,
		Content: content,
	}
}

func nodeForString(value string) *yaml.Node {
	return &yaml.Node{
		Kind:  yaml.ScalarNode,
		Tag:   "!!str",
		Value: value,
	}
}

func nodeForBoolean(value bool) *yaml.Node {
	return &yaml.Node{
		Kind:  yaml.ScalarNode,
		Tag:   "!!bool",
		Value: fmt.Sprintf("%t", value),
	}
}

func nodeForInt64(value int64) *yaml.Node {
	return &yaml.Node{
		Kind:  yaml.ScalarNode,
		Tag:   "!!int",
		Value: fmt.Sprintf("%d", value),
	}
}

func nodeForNull(value interface{}) *yaml.Node {
	return &yaml.Node{
		Kind:  yaml.ScalarNode,
		Tag:   "!!null",
		Value: "null",
	}
}

func nodeForNumber(value *SchemaNumber) *yaml.Node {
	if value.Integer != nil {
		return nodeForInt64(*value.Integer)
	} else if value.Float != nil {
		return nodeForFloat64(*value.Float)
	} else {
		return nil
	}
}

func nodeForFloat64(value float64) *yaml.Node {
	return &yaml.Node{
		Kind:  yaml.ScalarNode,
		Tag:   "!!float",
		Value: fmt.Sprintf("%f", value),
	}
}

func nodeForArray(value []interface{}) *yaml.Node {
	content := make([]*yaml.Node, 0)
	for _, item := range value {
		switch item.(type) {
		case string:
			content = append(content, nodeForString(item.(string)))
		case int:
			content = append(content, nodeForInt64(int64(item.(int))))
		case int64:
			content = append(content, nodeForInt64(item.(int64)))
		case bool:
			content = append(content, nodeForBoolean(item.(bool)))
		case float64:
			content = append(content, nodeForFloat64(item.(float64)))
		case nil:
			content = append(content, nodeForNull(item))
		default:
			fmt.Printf("nodeForArray: unexpected type %T\n", item)
		}
	}
	return &yaml.Node{
		Kind:    yaml.SequenceNode,
		Content: content,
	}
}

func nodeForMap(value map[string]interface{}) *yaml.Node {
	content := make([]*yaml.Node, 0)
	for key, item := range value {
		switch item.(type) {
		case string:
			content = append(content, nodeForString(key))
			content = append(content, nodeForString(item.(string)))
		case int:
			content = append(content, nodeForString(key))
			content = append(content, nodeForInt64(int64(item.(int))))
		case int64:
			content = append(content, nodeForString(key))
			content = append(content, nodeForInt64(item.(int64)))
		case bool:
			content = append(content, nodeForString(key))
			content = append(content, nodeForBoolean(item.(bool)))
		case float64:
			content = append(content, nodeForString(key))
			content = append(content, nodeForFloat64(item.(float64)))
		case nil:
			content = append(content, nodeForString(key))
			content = append(content, nodeForNull(item))
		default:
			fmt.Printf("nodeForMap: unexpected type %T\n", item)
		}
	}
	return &yaml.Node{
		Kind:    yaml.MappingNode,
		Content: content,
	}
}

func appendPair(nodes []*yaml.Node, name string, value *yaml.Node) []*yaml.Node {
	nodes = append(nodes, nodeForString(name))
	nodes = append(nodes, value)
	return nodes
}

func (absolute *Absolute) nodeValue() *yaml.Node {
	n := &yaml.Node{Kind: yaml.MappingNode}
	content := make([]*yaml.Node, 0)
	if absolute.ID != nil {
		switch strings.TrimSuffix(*absolute.Schema, "#") {
		case "http://json-schema.org/draft-07/schema":
			fallthrough
		case "#":
			fallthrough
		default:
			content = appendPair(content, "$id", nodeForString(*absolute.ID))
		}
	}
	if absolute.Schema != nil {
		content = appendPair(content, "$schema", nodeForString(*absolute.Schema))
	}
	if absolute.Ref != nil {
		content = appendPair(content, "$ref", nodeForString(*absolute.Ref))
	}
	if absolute.Comment != nil {
		content = appendPair(content, "$comment", nodeForString(*absolute.Comment))
	}
	if absolute.Title != nil {
		content = appendPair(content, "title", nodeForString(*absolute.Title))
	}
	if absolute.Description != nil {
		content = appendPair(content, "description", nodeForString(*absolute.Description))
	}

	if absolute.Default != nil {
		content = appendPair(content, "default", absolute.Default)
	}
	if absolute.ReadOnly != nil && *absolute.ReadOnly {
		content = appendPair(content, "readOnly", nodeForBoolean(*absolute.ReadOnly))
	}
	if absolute.WriteOnly != nil && *absolute.WriteOnly {
		content = appendPair(content, "writeOnly", nodeForBoolean(*absolute.WriteOnly))
	}
	if absolute.Examples != nil {
		content = appendPair(content, "examples", nodeForSchemaArray(*absolute.Examples))
	}

	if absolute.MultipleOf != nil {
		content = appendPair(content, "multipleOf", absolute.MultipleOf.nodeValue())
	}
	if absolute.Maximum != nil {
		content = appendPair(content, "maximum", absolute.Maximum.nodeValue())
	}
	if absolute.ExclusiveMaximum != nil {
		content = appendPair(content, "exclusiveMaximum", absolute.ExclusiveMaximum.nodeValue())
	}
	if absolute.Minimum != nil {
		content = appendPair(content, "minimum", absolute.Minimum.nodeValue())
	}
	if absolute.ExclusiveMinimum != nil {
		content = appendPair(content, "exclusiveMinimum", absolute.ExclusiveMinimum.nodeValue())
	}

	if absolute.MaxLength != nil {
		content = appendPair(content, "maxLength", nodeForInt64(*absolute.MaxLength))
	}
	if absolute.MinLength != nil {
		content = appendPair(content, "minLength", nodeForInt64(*absolute.MinLength))
	}
	if absolute.Pattern != nil {
		content = appendPair(content, "pattern", nodeForString(*absolute.Pattern))
	}

	if absolute.AdditionalItems != nil {
		content = appendPair(content, "additionalItems", absolute.AdditionalItems.nodeValue())
	}
	if absolute.Items != nil {
		content = appendPair(content, "items", absolute.Items.nodeValue())
	}
	if absolute.MaxItems != nil {
		content = appendPair(content, "maxItems", nodeForInt64(*absolute.MaxItems))
	}
	if absolute.MinItems != nil {
		content = appendPair(content, "minItems", nodeForInt64(*absolute.MinItems))
	}
	if absolute.UniqueItems != nil {
		content = appendPair(content, "uniqueItems", nodeForBoolean(*absolute.UniqueItems))
	}

	if absolute.Contains != nil {
		content = appendPair(content, "contains", absolute.Contains.nodeValue())
	}
	if absolute.MaxProperties != nil {
		content = appendPair(content, "maxProperties", nodeForInt64(*absolute.MaxProperties))
	}
	if absolute.MinProperties != nil {
		content = appendPair(content, "minProperties", nodeForInt64(*absolute.MinProperties))
	}
	if absolute.Required != nil {
		content = appendPair(content, "required", nodeForStringArray(*absolute.Required))
	}
	if absolute.AdditionalProperties != nil {
		content = appendPair(content, "additionalProperties", absolute.AdditionalProperties.nodeValue())
	}
	if absolute.Definitions != nil {
		content = appendPair(content, "definitions", nodeForNamedSchemaArray(absolute.Definitions))
	}
	if absolute.Properties != nil {
		content = appendPair(content, "properties", nodeForNamedSchemaArray(absolute.Properties))
	}
	if absolute.PatternProperties != nil {
		content = appendPair(content, "patternProperties", nodeForNamedSchemaArray(absolute.PatternProperties))
	}
	if absolute.Dependencies != nil {
		content = appendPair(content, "dependencies", nodeForNamedSchemaOrStringArray(absolute.Dependencies))
	}
	if absolute.PropertyNames != nil {
		content = appendPair(content, "propertyNames", absolute.PropertyNames.nodeValue())
	}

	if absolute.Const != nil {
		content = appendPair(content, "const", absolute.Const)
	}
	if absolute.Enumeration != nil {
		content = appendPair(content, "enum", nodeForSchemaEnumArray(absolute.Enumeration))
	}
	if absolute.Type != nil {
		content = appendPair(content, "type", absolute.Type.nodeValue())
	}
	if absolute.Format != nil {
		content = appendPair(content, "format", nodeForString(*absolute.Format))
	}
	if absolute.ContentMediaType != nil {
		content = appendPair(content, "contentMediaType", nodeForString(*absolute.ContentMediaType))
	}
	if absolute.ContentEncoding != nil {
		content = appendPair(content, "contentEncoding", nodeForString(*absolute.ContentEncoding))
	}

	if absolute.If != nil {
		content = appendPair(content, "if", absolute.If.nodeValue())
	}
	if absolute.Then != nil {
		content = appendPair(content, "then", absolute.Then.nodeValue())
	}
	if absolute.Else != nil {
		content = appendPair(content, "else", absolute.Else.nodeValue())
	}
	if absolute.AllOf != nil {
		content = appendPair(content, "allOf", nodeForSchemaArray(*absolute.AllOf))
	}
	if absolute.AnyOf != nil {
		content = appendPair(content, "anyOf", nodeForSchemaArray(*absolute.AnyOf))
	}
	if absolute.OneOf != nil {
		content = appendPair(content, "oneOf", nodeForSchemaArray(*absolute.OneOf))
	}
	if absolute.Not != nil {
		content = appendPair(content, "not", absolute.Not.nodeValue())
	}

	n.Content = content
	return n
}

// JSONString returns a json representation of a schema.
func (schema *Schema) JSONString() string {
	node := schema.nodeValue()
	return Render(node)
}
