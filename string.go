package nullable

import (
	"encoding/json"
	"strings"
)

type String struct {
	Set   bool
	Value *string
}

func (s *String) SetValue(value string) {
	s.SetPtr(&value)
}

func (s *String) SetPtr(value *string) {
	s.Set = true
	s.Value = value
}

func (s *String) UnmarshalJSON(data []byte) error {
	s.Set = true
	return json.Unmarshal(data, &s.Value)
}

func (s *String) Removed() bool {
	return s.Set && s.Value == nil
}

func (s *String) IsEmpty() bool {
	return s.Set && s.Value != nil && strings.TrimSpace(*s.Value) == ""
}
