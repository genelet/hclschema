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

func (object *Combined) nodeValue() *yaml.Node {
	if object.Schema != nil {
		return object.Schema.nodeValue()
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

func nodeForCombinedArray(array []*Combined) *yaml.Node {
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

func (object *CombinedOrStringArray) nodeValue() *yaml.Node {
	if object.Combined != nil {
		return object.Combined.nodeValue()
	} else if object.StringArray != nil {
		return nodeForStringArray(*(object.StringArray))
	} else {
		return nil
	}
}

func (object *CombinedOrCombinedArray) nodeValue() *yaml.Node {
	if object.Combined != nil {
		return object.Combined.nodeValue()
	} else if object.CombinedArray != nil {
		return nodeForCombinedArray(*(object.CombinedArray))
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

func nodeForNamedCombinedArray(array *[]*NamedCombined) *yaml.Node {
	content := make([]*yaml.Node, 0)
	for _, pair := range *(array) {
		content = appendPair(content, pair.Name, pair.Value.nodeValue())
	}
	return nodeForMapping(content)
}

func nodeForNamedCombinedOrStringArray(array *[]*NamedCombinedOrStringArray) *yaml.Node {
	content := make([]*yaml.Node, 0)
	for _, pair := range *(array) {
		content = appendPair(content, pair.Name, pair.Value.nodeValue())
	}
	return nodeForMapping(content)
}

func nodeForCombinedEnumArray(array *[]SchemaEnumValue) *yaml.Node {
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

func (schema *Schema) nodeValue() *yaml.Node {
	n := &yaml.Node{Kind: yaml.MappingNode}
	content := make([]*yaml.Node, 0)
	if schema.ID != nil {
		switch strings.TrimSuffix(*schema.Schema, "#") {
		case "http://json-schema.org/draft-07/schema":
			fallthrough
		case "#":
			fallthrough
		default:
			content = appendPair(content, "$id", nodeForString(*schema.ID))
		}
	}
	if schema.Schema != nil {
		content = appendPair(content, "$schema", nodeForString(*schema.Schema))
	}
	if schema.Ref != nil {
		content = appendPair(content, "$ref", nodeForString(*schema.Ref))
	}
	if schema.Comment != nil {
		content = appendPair(content, "$comment", nodeForString(*schema.Comment))
	}
	if schema.Title != nil {
		content = appendPair(content, "title", nodeForString(*schema.Title))
	}
	if schema.Description != nil {
		content = appendPair(content, "description", nodeForString(*schema.Description))
	}

	if schema.Default != nil {
		content = appendPair(content, "default", schema.Default)
	}
	if schema.ReadOnly != nil && *schema.ReadOnly {
		content = appendPair(content, "readOnly", nodeForBoolean(*schema.ReadOnly))
	}
	if schema.WriteOnly != nil && *schema.WriteOnly {
		content = appendPair(content, "writeOnly", nodeForBoolean(*schema.WriteOnly))
	}
	if schema.Examples != nil {
		content = appendPair(content, "examples", nodeForCombinedArray(*schema.Examples))
	}

	if schema.MultipleOf != nil {
		content = appendPair(content, "multipleOf", schema.MultipleOf.nodeValue())
	}
	if schema.Maximum != nil {
		content = appendPair(content, "maximum", schema.Maximum.nodeValue())
	}
	if schema.ExclusiveMaximum != nil {
		content = appendPair(content, "exclusiveMaximum", schema.ExclusiveMaximum.nodeValue())
	}
	if schema.Minimum != nil {
		content = appendPair(content, "minimum", schema.Minimum.nodeValue())
	}
	if schema.ExclusiveMinimum != nil {
		content = appendPair(content, "exclusiveMinimum", schema.ExclusiveMinimum.nodeValue())
	}

	if schema.MaxLength != nil {
		content = appendPair(content, "maxLength", nodeForInt64(*schema.MaxLength))
	}
	if schema.MinLength != nil {
		content = appendPair(content, "minLength", nodeForInt64(*schema.MinLength))
	}
	if schema.Pattern != nil {
		content = appendPair(content, "pattern", nodeForString(*schema.Pattern))
	}

	if schema.AdditionalItems != nil {
		content = appendPair(content, "additionalItems", schema.AdditionalItems.nodeValue())
	}
	if schema.Items != nil {
		content = appendPair(content, "items", schema.Items.nodeValue())
	}
	if schema.MaxItems != nil {
		content = appendPair(content, "maxItems", nodeForInt64(*schema.MaxItems))
	}
	if schema.MinItems != nil {
		content = appendPair(content, "minItems", nodeForInt64(*schema.MinItems))
	}
	if schema.UniqueItems != nil {
		content = appendPair(content, "uniqueItems", nodeForBoolean(*schema.UniqueItems))
	}

	if schema.Contains != nil {
		content = appendPair(content, "contains", schema.Contains.nodeValue())
	}
	if schema.MaxProperties != nil {
		content = appendPair(content, "maxProperties", nodeForInt64(*schema.MaxProperties))
	}
	if schema.MinProperties != nil {
		content = appendPair(content, "minProperties", nodeForInt64(*schema.MinProperties))
	}
	if schema.Required != nil {
		content = appendPair(content, "required", nodeForStringArray(*schema.Required))
	}
	if schema.AdditionalProperties != nil {
		content = appendPair(content, "additionalProperties", schema.AdditionalProperties.nodeValue())
	}
	if schema.Definitions != nil {
		content = appendPair(content, "definitions", nodeForNamedCombinedArray(schema.Definitions))
	}
	if schema.Properties != nil {
		content = appendPair(content, "properties", nodeForNamedCombinedArray(schema.Properties))
	}
	if schema.PatternProperties != nil {
		content = appendPair(content, "patternProperties", nodeForNamedCombinedArray(schema.PatternProperties))
	}
	if schema.Dependencies != nil {
		content = appendPair(content, "dependencies", nodeForNamedCombinedOrStringArray(schema.Dependencies))
	}
	if schema.PropertyNames != nil {
		content = appendPair(content, "propertyNames", schema.PropertyNames.nodeValue())
	}

	if schema.Const != nil {
		content = appendPair(content, "const", schema.Const)
	}
	if schema.Enumeration != nil {
		content = appendPair(content, "enum", nodeForCombinedEnumArray(schema.Enumeration))
	}
	if schema.Type != nil {
		content = appendPair(content, "type", schema.Type.nodeValue())
	}
	if schema.Format != nil {
		content = appendPair(content, "format", nodeForString(*schema.Format))
	}
	if schema.ContentMediaType != nil {
		content = appendPair(content, "contentMediaType", nodeForString(*schema.ContentMediaType))
	}
	if schema.ContentEncoding != nil {
		content = appendPair(content, "contentEncoding", nodeForString(*schema.ContentEncoding))
	}

	if schema.If != nil {
		content = appendPair(content, "if", schema.If.nodeValue())
	}
	if schema.Then != nil {
		content = appendPair(content, "then", schema.Then.nodeValue())
	}
	if schema.Else != nil {
		content = appendPair(content, "else", schema.Else.nodeValue())
	}
	if schema.AllOf != nil {
		content = appendPair(content, "allOf", nodeForCombinedArray(*schema.AllOf))
	}
	if schema.AnyOf != nil {
		content = appendPair(content, "anyOf", nodeForCombinedArray(*schema.AnyOf))
	}
	if schema.OneOf != nil {
		content = appendPair(content, "oneOf", nodeForCombinedArray(*schema.OneOf))
	}
	if schema.Not != nil {
		content = appendPair(content, "not", schema.Not.nodeValue())
	}

	n.Content = content
	return n
}

// JSONString returns a json representation of a schema.
func (combined *Combined) JSONString() string {
	node := combined.nodeValue()
	return Render(node)
}
