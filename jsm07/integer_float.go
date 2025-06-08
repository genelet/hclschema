package jsm07

import (
	"encoding/json"

	"github.com/genelet/determined/dethcl"
)

// IntegerOrFloat represents a value that can be either an Integer or a Float.
type IntegerOrFloat struct {
	Integer *int64
	Float   *float64
}

func (self *IntegerOrFloat) UnmarshalJSON(data []byte) error {
	var integer int64
	if err := json.Unmarshal(data, &integer); err == nil {
		self.Integer = &integer
		return nil
	}

	var float float64
	if err := json.Unmarshal(data, &float); err == nil {
		self.Float = &float
		return nil
	}

	return json.Unmarshal(data, &self.Integer) // Fallback to Integer if both fail
}

func (self *IntegerOrFloat) MarshalJSON() ([]byte, error) {
	if self.Integer != nil {
		return json.Marshal(*self.Integer)
	}
	if self.Float != nil {
		return json.Marshal(*self.Float)
	}
	return nil, nil // Return nil if both are nil
}

func (self *IntegerOrFloat) UnmarshalHCL(data []byte) error {
	var integer int64
	if err := dethcl.Unmarshal(data, &integer); err == nil {
		self.Integer = &integer
		return nil
	}

	var float float64
	if err := dethcl.Unmarshal(data, &float); err == nil {
		self.Float = &float
		return nil
	}

	return dethcl.Unmarshal(data, &self.Integer) // Fallback to Integer if both fail
}

func (self *IntegerOrFloat) MarshalHCL() ([]byte, error) {
	if self.Integer != nil {
		return dethcl.Marshal(*self.Integer)
	}
	if self.Float != nil {
		return dethcl.Marshal(*self.Float)
	}
	return nil, nil // Return nil if both are nil
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
