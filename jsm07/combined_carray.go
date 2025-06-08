package jsm07

import (
	"encoding/json"

	"github.com/genelet/determined/dethcl"
)

// CombinedOrCombinedArray represents a value that can be either
// a Combined or an Array of Combineds.
type CombinedOrCombinedArray struct {
	Combined      *Combined
	CombinedArray *[]*Combined
}

// UnmarshalJSON implements the json.Unmarshaler interface for CombinedOrCombinedArray.
func (s *CombinedOrCombinedArray) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		return nil // Handle empty data gracefully
	}

	var combined Combined
	if err := json.Unmarshal(data, &combined); err == nil {
		s.Combined = &combined
		return nil
	}

	var arr []*Combined
	if err := json.Unmarshal(data, &arr); err == nil {
		s.CombinedArray = &arr
		return nil
	}

	return json.Unmarshal(data, &s.Combined) // Fallback to Combined if both fail
}

func (s *CombinedOrCombinedArray) MarshalJSON() ([]byte, error) {
	if s.Combined != nil {
		return json.Marshal(s.Combined)
	}
	if s.CombinedArray != nil {
		return json.Marshal(*s.CombinedArray)
	}
	return nil, nil // Return nil if both are nil
}

func (s *CombinedOrCombinedArray) UnmarshalHCL(data []byte) error {
	if len(data) == 0 {
		return nil // Handle empty data gracefully
	}
	var combined Combined
	if err := dethcl.Unmarshal(data, &combined); err == nil {
		s.Combined = &combined
		return nil
	}
	var arr []*Combined
	if err := dethcl.Unmarshal(data, &arr); err == nil {
		s.CombinedArray = &arr
		return nil
	}
	return dethcl.Unmarshal(data, &s.Combined) // Fallback to Combined if both fail
}

func (s *CombinedOrCombinedArray) MarshalHCL() ([]byte, error) {
	if s.Combined != nil {
		return dethcl.Marshal(s.Combined)
	}
	if s.CombinedArray != nil {
		return dethcl.Marshal(*s.CombinedArray)
	}
	return nil, nil // Return nil if both are nil
}

// NewCombinedOrCombinedArrayWithCombined creates and returns a new object
func NewCombinedOrCombinedArrayWithCombined(s *Combined) *CombinedOrCombinedArray {
	result := &CombinedOrCombinedArray{}
	result.Combined = s
	return result
}

// NewCombinedOrCombinedArrayWithCombinedArray creates and returns a new object
func NewCombinedOrCombinedArrayWithCombinedArray(a []*Combined) *CombinedOrCombinedArray {
	result := &CombinedOrCombinedArray{}
	result.CombinedArray = &a
	return result
}
