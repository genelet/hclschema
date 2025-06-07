package marshal07

import "encoding/json"

// CombinedOrStringArray represents a value that can be either
// a Combined or an Array of Strings.
type CombinedOrStringArray struct {
	Combined    *Combined
	StringArray *[]string
}

// UnmarshalJSON implements the json.Unmarshaler interface for CombinedOrStringArray.
func (s *CombinedOrStringArray) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		return nil // Handle empty data gracefully
	}

	var combined Combined
	if err := json.Unmarshal(data, &combined); err == nil {
		s.Combined = &combined
		return nil
	}

	var arr []string
	if err := json.Unmarshal(data, &arr); err == nil {
		s.StringArray = &arr
		return nil
	}

	return json.Unmarshal(data, &s.Combined) // Fallback to Combined if both fail
}

func (s *CombinedOrStringArray) MarshalJSON() ([]byte, error) {
	if s.Combined != nil {
		return json.Marshal(s.Combined)
	}
	if s.StringArray != nil {
		return json.Marshal(*s.StringArray)
	}
	return nil, nil // Return nil if both are nil
}

func NewCombinedOrStringArrayWithCombined(s *Combined) *CombinedOrStringArray {
	result := &CombinedOrStringArray{}
	result.Combined = s
	return result
}
func NewCombinedOrStringArrayWithStringArray(a []string) *CombinedOrStringArray {
	result := &CombinedOrStringArray{}
	result.StringArray = &a
	return result
}
