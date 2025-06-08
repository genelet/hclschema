package jsm07

import (
	"encoding/json"

	"github.com/genelet/determined/dethcl"
)

// Combined represents a value that can be either a Schema or a Boolean.
type Combined struct {
	Schema  *Schema
	Boolean *bool
}

func (self *Combined) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		return nil // Handle empty data gracefully
	}

	var schema Schema
	if err := json.Unmarshal(data, &schema); err == nil {
		self.Schema = &schema
		return nil
	}

	var boolean bool
	if err := json.Unmarshal(data, &boolean); err == nil {
		self.Boolean = &boolean
		return nil
	}

	return json.Unmarshal(data, &self.Schema) // Fallback to Schema if both fail
}

func (self *Combined) MarshalJSON() ([]byte, error) {
	if self.Schema != nil {
		return json.Marshal(self.Schema)
	}
	if self.Boolean != nil {
		return json.Marshal(*self.Boolean)
	}
	return nil, nil // Return nil if both are nil
}

func (self *Combined) UnmarshalHCL(data []byte) error {
	if len(data) == 0 {
		return nil
	}

	var schema Schema
	if err := dethcl.Unmarshal(data, &schema); err == nil {
		self.Schema = &schema
		return nil
	}

	var boolean bool
	if err := dethcl.Unmarshal(data, &boolean); err == nil {
		self.Boolean = &boolean
		return nil
	}

	return dethcl.Unmarshal(data, &self.Schema) // Fallback to Schema if both fail
}

func (self *Combined) MarshalHCL() ([]byte, error) {
	if self.Boolean != nil {
		if *self.Boolean {
			return []byte("true"), nil
		} else {
			return []byte("false"), nil
		}
	}

	if self.Schema != nil {
		return dethcl.Marshal(self.Schema)
	}

	return nil, nil // Return nil if both are nil
}

// NewCombinedWithSchema creates and returns a new object
func NewCombinedWithSchema(s *Schema) *Combined {
	if s == nil {
		return nil
	}
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
