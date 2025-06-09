package jsm07

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/genelet/determined/dethcl"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestCombinedOrCombinedArrayUnmarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected *CombinedOrCombinedArray
	}{
		{
			name:  "Combined",
			input: `{"schema": {"type": "string"}}`,
			expected: &CombinedOrCombinedArray{
				Combined: &Combined{
					Schema: &Schema{Type: NewStringOrStringArrayWithString("string")},
				},
			},
		},
		{
			name:  "CombinedArray",
			input: `[{"schema": {"type": "integer"}}, {"schema": {"type": "boolean"}}]`,
			expected: &CombinedOrCombinedArray{
				CombinedArray: &[]*Combined{
					{Schema: &Schema{Type: NewStringOrStringArrayWithString("integer")}},
					{Schema: &Schema{Type: NewStringOrStringArrayWithString("boolean")}},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result CombinedOrCombinedArray
			if err := result.UnmarshalJSON([]byte(tt.input)); err != nil {
				t.Fatalf("UnmarshalJSON failed: %v", err)
			}
			if result.Combined == nil && result.CombinedArray == nil {
				t.Fatal("Expected either Combined or CombinedArray to be set")
			}
			if (result.Combined != nil && tt.expected.Combined == nil) ||
				(result.CombinedArray != nil && tt.expected.CombinedArray == nil) {
				t.Fatal("Unmarshalled result does not match expected")
			}
		})
	}
}

// TestMCPJSON tests the Marshal and Unmarshal functionality of the mcp schema in the samples directory.
func TestMCPJSON(t *testing.T) {
	bs, err := os.ReadFile("samples/mcp.json")
	if err != nil {
		t.Fatalf("Failed to read mcp.json: %v", err)
	}
	mcp := new(Schema)
	if err := json.Unmarshal(bs, mcp); err != nil {
		t.Fatalf("Failed to unmarshal mcp.json: %v", err)
	}

	bs1, err := json.Marshal(mcp)
	if err != nil {
		t.Fatalf("Failed to marshal mcp: %v", err)
	}

	mcp1 := new(Schema)
	if err := json.Unmarshal(bs1, mcp1); err != nil {
		t.Fatalf("Failed to unmarshal marshalled mcp: %v", err)
	}

	if diff := cmp.Diff(mcp, mcp1); diff != "" {
		t.Errorf("MCP schema mismatch (-want +got):\n%s", diff)
	}
}

// TestMCPHCL tests the Marshal and Unmarshal functionality of the mcp schema in the samples directory.
func TestMCPHCL(t *testing.T) {
	bs, err := os.ReadFile("samples/mcp.json")
	if err != nil {
		t.Fatalf("Failed to read mcp.json: %v", err)
	}
	mcp := new(Schema)
	if err := json.Unmarshal(bs, mcp); err != nil {
		t.Fatalf("Failed to unmarshal mcp.json: %v", err)
	}

	bs1, err := dethcl.Marshal(mcp)
	if err != nil {
		t.Fatalf("Failed to marshal mcp: %v", err)
	}

	mcp1 := new(Schema)
	if err := dethcl.Unmarshal(bs1, mcp1); err != nil {
		t.Fatalf("Failed to unmarshal mcp: %v", err)
	}

	if diff := cmp.Diff(mcp, mcp1, cmpopts.IgnoreFields(Schema{}, "Description")); diff != "" {
		t.Errorf("MCP schema mismatch (-want +got):\n%s", diff)
	}

	bs2, err := dethcl.Marshal(mcp)
	if err != nil {
		t.Fatalf("Failed to marshal mcp: %v", err)
	}

	mcp2 := new(Schema)
	if err := dethcl.Unmarshal(bs2, mcp2); err != nil {
		t.Fatalf("Failed to unmarshal mcp: %v", err)
	}

	if diff := cmp.Diff(mcp, mcp2, cmpopts.IgnoreFields(Schema{}, "Description")); diff != "" {
		t.Errorf("MCP schema mismatch (-want +got):\n%s", diff)
	}
}
