package nullable

import "encoding/json"

type Bool struct {
	set   bool
	value *bool
}

// NewBool returns a Bool set to the given value.
func NewBool(v bool) Bool {
	return Bool{
		set:   true,
		value: &v,
	}
}

// NewBoolPtr returns a Bool set to the given pointer.
func NewBoolPtr(v *bool) Bool {
	return Bool{
		set:   true,
		value: v,
	}
}

func (b *Bool) SetValue(value bool) {
	b.SetPtr(&value)
}

func (b *Bool) SetPtr(value *bool) {
	b.set = true
	b.value = value
}

func (b Bool) Value() *bool {
	return b.value
}

func (b Bool) Equals(value bool) bool {
	return b.value != nil && *b.value == value
}

// Apply returns the given value, the zero value (false), or b's value,
// depending on whether b is unset, removed, or set, respectively.
func (b Bool) Apply(value bool) bool {
	if !b.set {
		return value
	}
	if b.value == nil {
		return false
	}
	return *b.value
}

// ApplyPtr returns the given value, nil, or b's value, depending on whether b
// is unset, removed, or set, respectively.
func (b Bool) ApplyPtr(value *bool) *bool {
	if b.set {
		return b.value
	}
	return value
}

// Diff returns b if b.Apply(value) != value or else an unset Bool. This can be
// used to avoid extraneous updates when the update would have no effect.
func (b Bool) Diff(value bool) Bool {
	if b.Apply(value) == value {
		return Bool{}
	}
	return b
}

func (b *Bool) UnmarshalJSON(data []byte) error {
	b.set = true
	return json.Unmarshal(data, &b.value)
}

func (b Bool) IsSet() bool {
	return b.set
}

func (b Bool) Removed() bool {
	return b.set && b.value == nil
}

func (b Bool) InterfaceValue() interface{} {
	return b.value
}
