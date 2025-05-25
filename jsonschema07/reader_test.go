package jsonschema07

import (
	"reflect"
	"testing"

	"gopkg.in/yaml.v3"
)

// TestParseSchema tests the parsing of a JSON schema file.
func TestParseSchema(t *testing.T) {
	schema, err := NewSchemaFromFile("samples/mcp.json")
	if err != nil {
		t.Fatalf("Error parsing schema: %v", err)
	}

	//for _, def := range *schema.Definitions {
	//	if def.Name == "CallToolResult" {
	//		for _, prop := range *def.Value.Schema.Properties {
	//			if prop.Name == "content" {
	//				t.Errorf("Ref: %s => %#v", prop.Name, prop.Value.Schema.Items.Combined.Schema.AnyOf)
	//			}
	//		}
	//	}
	//}
	//jstring, err := json.MarshalIndent(schema, "", "  ")
	//if err != nil {
	//	t.Fatalf("Error marshaling schema to JSON: %v", err)
	//}
	//t.Errorf("%s", string(jstring))

	str := []byte(schema.String())
	t.Errorf("%s", str)
	var node yaml.Node
	err = yaml.Unmarshal(str, &node)
	if err != nil {
		t.Fatal(err)
	}
	schema1 := NewSchemaFromObject(&node)
	if schema1 == nil {
		t.Fatal("schema1 is nil")
	}
	if reflect.DeepEqual(schema, schema1) {
		t.Fatal("schema and schema1 are equal")
	}
}
