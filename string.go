package nullable

import (
	"encoding/json"
	"strings"
)

type String struct {
	set   bool
	value *string
}

// NewString returns a String set to the given value.
func NewString(v string) String {
	return String{
		set:   true,
		value: &v,
	}
}

// NewStringPtr returns a String set to the given pointer.
func NewStringPtr(v *string) String {
	return String{
		set:   true,
		value: v,
	}
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

func (s String) Equals(value string) bool {
	return s.value != nil && *s.value == value
}

// Apply returns the given value, the zero value (""), or s's value, depending
// on whether s is unset, removed, or set, respectively.
func (s String) Apply(value string) string {
	if !s.set {
		return value
	}
	if s.value == nil {
		return ""
	}
	return *s.value
}

// ApplyPtr returns the given value, nil, or s's value, depending on whether s
// is unset, removed, or set, respectively.
func (s String) ApplyPtr(value *string) *string {
	if s.set {
		return s.value
	}
	return value
}

// Diff returns s if s.Apply(value) != value; otherwise it returns an unset
// String. This can be used to omit extraneous updates when applying the update
// would have no effect.
func (s String) Diff(value string) String {
	if s.Apply(value) == value {
		return String{}
	}
	return s
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

func (s String) InterfaceValue() interface{} {
	return s.value
}

// IsEmpty checks whether the String is set to a string containing only whitespace.
func (s String) IsEmpty() bool {
	return s.set && s.value != nil && strings.TrimSpace(*s.value) == ""
}

func (s String) String() string {
	if s.Removed() {
		return "<removed>"
	}
	if !s.IsSet() {
		return "<unset>"
	}
	return *s.Value()
}
