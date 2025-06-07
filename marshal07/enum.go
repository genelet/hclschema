package marshal07

import "encoding/json"

// SchemaEnumValue represents a value that can be part of an
// enumeration in a Combined.
type SchemaEnumValue struct {
	String *string
	Bool   *bool
	Number *IntegerOrFloat
	Null   *bool
}

func (s *SchemaEnumValue) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		s.Null = new(bool)
		*s.Null = true
		return nil
	}

	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		s.String = &str
		return nil
	}

	var boolean bool
	if err := json.Unmarshal(data, &boolean); err == nil {
		s.Bool = &boolean
		return nil
	}

	var number IntegerOrFloat
	if err := json.Unmarshal(data, &number); err == nil {
		s.Number = &number
		return nil
	}

	return json.Unmarshal(data, &s.String) // Fallback to String if all fail
}

func (s *SchemaEnumValue) MarshalJSON() ([]byte, error) {
	if s.Null != nil && *s.Null {
		return []byte("null"), nil
	}

	if s.String != nil {
		return json.Marshal(*s.String)
	}
	if s.Bool != nil {
		return json.Marshal(*s.Bool)
	}
	if s.Number != nil {
		return json.Marshal(s.Number)
	}

	return nil, nil // Return nil if all are nil
}
