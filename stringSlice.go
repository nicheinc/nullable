package nullable

import "encoding/json"

// StringSlice implements Nullable for []string fields.
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

// SetValue modifies the receiver to be an update to the given value.
func (s *StringSlice) SetValue(value []string) {
	s.set = true
	s.value = value
}

// Value returns nil if the receiver is unset/removed or else the updated value.
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

// Diff returns s if s.Apply(value) is not pairwise equal to value; otherwise it
// returns an unset StringSlice. This can be used to omit extraneous updates
// when applying the update would have no effect.
func (s StringSlice) Diff(value []string) StringSlice {
	if stringSliceEquals(s.Apply(value), value) {
		return StringSlice{}
	}
	return s
}

// UnmarshalJSON implements json.Unmarshaler.
func (s *StringSlice) UnmarshalJSON(data []byte) error {
	s.set = true
	return json.Unmarshal(data, &s.value)
}

// IsSet returns true if the receiver has been set/removed.
func (s StringSlice) IsSet() bool {
	return s.set
}

// Removed returns whether the receiver has been removed (value set to nil).
func (s StringSlice) Removed() bool {
	return s.set && s.value == nil
}

// InterfaceValue returns value as an interface{}.
func (s StringSlice) InterfaceValue() interface{} {
	return s.value
}

// IsEmpty returns whether the receiver has been set to a non-nil empty slice.
func (s *StringSlice) IsEmpty() bool {
	return s.set && s.value != nil && len(s.value) == 0
}
