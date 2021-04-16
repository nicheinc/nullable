package nullable

import (
	"encoding/json"
	"strings"
)

type String struct {
	set   bool
	value *string
}

func (s *String) SetValue(value string) {
	s.SetPtr(&value)
}

func (s *String) SetPtr(value *string) {
	s.set = true
	s.value = value
}

func (s String) Value() *string {
	return s.value
}

func (s *String) UnmarshalJSON(data []byte) error {
	s.set = true
	return json.Unmarshal(data, &s.value)
}

func (s String) IsSet() bool {
	return s.set
}

func (s String) Removed() bool {
	return s.set && s.value == nil
}

func (s String) interfaceValue() interface{} {
	return s.value
}

func (s String) IsEmpty() bool {
	return s.set && s.value != nil && strings.TrimSpace(*s.value) == ""
}
