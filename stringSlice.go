package nullable

import "encoding/json"

type StringSlice struct {
	set   bool
	value []string
}

// MakeStringSlice returns a StringSlice set to the given value.
func MakeStringSlice(v []string) StringSlice {
	return StringSlice{
		set:   true,
		value: v,
	}
}

func (s *StringSlice) SetValue(value []string) {
	s.set = true
	s.value = value
}

func (s StringSlice) Value() []string {
	return s.value
}

func (s *StringSlice) UnmarshalJSON(data []byte) error {
	s.set = true
	return json.Unmarshal(data, &s.value)
}

func (s StringSlice) IsSet() bool {
	return s.set
}

func (s StringSlice) Removed() bool {
	return s.set && s.value == nil
}

func (s StringSlice) interfaceValue() interface{} {
	return s.value
}

func (s *StringSlice) IsEmpty() bool {
	return s.set && s.value != nil && len(s.value) == 0
}
