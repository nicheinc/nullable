package nullable

import "encoding/json"

type StringSlice struct {
	Set   bool
	Value []string
}

func (s *StringSlice) SetValue(value []string) {
	s.Set = true
	s.Value = value
}

func (s *StringSlice) UnmarshalJSON(data []byte) error {
	s.Set = true
	return json.Unmarshal(data, &s.Value)
}

func (s *StringSlice) Removed() bool {
	return s.Set && s.Value == nil
}

func (s *StringSlice) IsEmpty() bool {
	return s.Set && s.Value != nil && len(s.Value) == 0
}
