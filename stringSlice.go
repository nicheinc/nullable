package nullable

import "encoding/json"

type StringSlice struct {
	set   bool
	value []string
}

// NewStringSlice returns a StringSlice set to the given value.
func NewStringSlice(v []string) StringSlice {
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

// Equals returns whether s is set to a non-nil []string that is element-wise
// equal to the given []string.
func (s StringSlice) Equals(value []string) bool {
	if s.value == nil || len(s.value) != len(value) {
		return false
	}
	for i := range s.value {
		if s.value[i] != value[i] {
			return false
		}
	}
	return true
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

func (s StringSlice) InterfaceValue() interface{} {
	return s.value
}

func (s *StringSlice) IsEmpty() bool {
	return s.set && s.value != nil && len(s.value) == 0
}
