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
	if s.value == nil {
		return false
	}
	return stringSliceEquals(s.value, value)
}

func stringSliceEquals(slice1 []string, slice2 []string) bool {
	if len(slice1) != len(slice2) {
		return false
	}
	for i := range slice1 {
		if slice1[i] != slice2[i] {
			return false
		}
	}
	return true
}

// Apply returns the given value, the zero value (nil), or s's value, depending
// on whether s is unset, removed, or set, respectively.
func (s StringSlice) Apply(value []string) []string {
	if !s.set {
		return value
	}
	return s.value
}

// Diff returns the "simplest" s2 such that s2.Apply(value) = s.Apply(value).
// "Simplest" means that if possible, the result will be unset.
func (s StringSlice) Diff(value []string) StringSlice {
	if stringSliceEquals(s.Apply(value), value) {
		return StringSlice{}
	}
	return s
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
