package jsm07

import (
	"encoding/json"

	"github.com/genelet/determined/dethcl"
)

// StringOrStringArray represents a value that can be either
// a String or an Array of Strings.
type StringOrStringArray struct {
	String      *string
	StringArray *[]string
}

// UnmarshalJSON implements the json.Unmarshaler interface for StringOrStringArray.
func (s *StringOrStringArray) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		s.String = &str
		return nil
	}

	var arr []string
	if err := json.Unmarshal(data, &arr); err == nil {
		s.StringArray = &arr
		return nil
	}

	return json.Unmarshal(data, &s.String) // Fallback to String if both fail
}

func (s *StringOrStringArray) MarshalJSON() ([]byte, error) {
	if s.String != nil {
		return json.Marshal(*s.String)
	}
	if s.StringArray != nil {
		return json.Marshal(*s.StringArray)
	}
	return nil, nil // Return nil if both are nil
}

func (s *StringOrStringArray) UnmarshalHCL(data []byte) error {
	var str string
	if err := dethcl.Unmarshal(data, &str); err == nil {
		s.String = &str
		return nil
	}
	var arr []string
	if err := dethcl.Unmarshal(data, &arr); err == nil {
		s.StringArray = &arr
		return nil
	}
	return dethcl.Unmarshal(data, &s.String) // Fallback to String if both fail
}

func (s *StringOrStringArray) MarshalHCL() ([]byte, error) {
	if s.String != nil {
		return dethcl.Marshal(*s.String)
	}
	if s.StringArray != nil {
		return dethcl.Marshal(*s.StringArray)
	}
	return nil, nil // Return nil if both are nil
}

// NewStringOrStringArrayWithString creates and returns a new object
func NewStringOrStringArrayWithString(s string) *StringOrStringArray {
	result := &StringOrStringArray{}
	result.String = &s
	return result
}

// NewStringOrStringArrayWithStringArray creates and returns a new object
func NewStringOrStringArrayWithStringArray(a []string) *StringOrStringArray {
	result := &StringOrStringArray{}
	result.StringArray = &a
	return result
}
