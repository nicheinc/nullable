package nullable

import (
	"encoding/json"
	"time"
)

// Time implements Nullable for time.Time fields.
type Time struct {
	set   bool
	value *time.Time
}

// NewTime returns a Time set to the given value.
func NewTime(v time.Time) Time {
	return Time{
		set:   true,
		value: &v,
	}
}

// NewTimePtr returns a Time set to the given pointer.
func NewTimePtr(v *time.Time) Time {
	return Time{
		set:   true,
		value: v,
	}
}

// SetValue modifies the receiver to be an update to the given value.
func (t *Time) SetValue(value time.Time) {
	t.SetPtr(&value)
}

// SetPtr modifies the receiver to be an update to the given value. If
// the value is nil, the receiver will be removed.
func (t *Time) SetPtr(value *time.Time) {
	t.set = true
	t.value = value
}

// Value returns nil if the receiver is unset/removed or else the updated value.
func (t Time) Value() *time.Time {
	return t.value
}

// Equals returns whether the receiver is set to the given value.
func (t Time) Equals(value time.Time) bool {
	return t.value != nil && t.value.Equal(value)
}

// Apply returns the given value, the zero value (0001-01-01 00:00:00
// +0000 UTC), or t's value, depending on whether t is unset, removed,
// or set, respectively.
func (t Time) Apply(value time.Time) time.Time {
	if !t.set {
		return value
	}
	if t.value == nil {
		return time.Time{}
	}
	return *t.value
}

// ApplyPtr returns the given value, nil, or t's value, depending on whether t
// is unset, removed, or set, respectively.
func (t Time) ApplyPtr(value *time.Time) *time.Time {
	if t.set {
		return t.value
	}
	return value
}

// Diff returns t if t.Apply(value) != value; otherwise it returns an
// unset Time. This can be used to omit extraneous updates when
// applying the update would have no effect.
func (t Time) Diff(value time.Time) Time {
	if t.Apply(value) == value {
		return Time{}
	}
	return t
}

// UnmarshalJSON implements json.Unmarshaler.
func (t *Time) UnmarshalJSON(data []byte) error {
	t.set = true
	return json.Unmarshal(data, &t.value)
}

// IsSet returns true if the receiver has been set/removed.
func (t Time) IsSet() bool {
	return t.set
}

// Removed returns whether the receiver has been removed (value set to nil).
func (t Time) Removed() bool {
	return t.set && t.value == nil
}

// InterfaceValue returns value as an interface{}.
func (t Time) InterfaceValue() interface{} {
	return t.value
}

// IsZero returns whether the receiver is set to (0001-01-01 00:00:00
// +0000 UTC).
func (t Time) IsZero() bool {
	return t.set && t.value != nil && t.value.IsZero()
}

// String implements fmt.Stringer. It returns "<removed>" if the
// receiver is removed, "<unset>" if it's unset, or else the contained
// time.
func (t Time) String() string {
	if t.Removed() {
		return "<removed>"
	}
	if !t.IsSet() {
		return "<unset>"
	}
	return t.value.String()
}
