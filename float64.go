package nullable

import "encoding/json"

// Float64 implements Nullable for float64 fields.
type Float64 struct {
	set   bool
	value *float64
}

// NewFloat64 returns a Float64 set to the given value.
func NewFloat64(v float64) Float64 {
	return Float64{
		set:   true,
		value: &v,
	}
}

// NewFloat64Ptr returns a Float64 set to the given pointer.
func NewFloat64Ptr(v *float64) Float64 {
	return Float64{
		set:   true,
		value: v,
	}
}

// SetValue modifies the receiver to be an update to the given value.
func (f *Float64) SetValue(value float64) {
	f.SetPtr(&value)
}

// SetPtr modifies the receiver to be an update to the given value. If the value
// is nil, the receiver will be removed.
func (f *Float64) SetPtr(value *float64) {
	f.set = true
	f.value = value
}

// Value returns nil if the receiver is unset/removed or else the updated value.
func (f Float64) Value() *float64 {
	return f.value
}

// Equals returns whether the receiver is set to the given value.
func (f Float64) Equals(value float64) bool {
	return f.value != nil && *f.value == value
}

// Apply returns the given value, the zero value (0), or f's value, depending on
// whether f is unset, removed, or set, respectively.
func (f Float64) Apply(value float64) float64 {
	if !f.set {
		return value
	}
	if f.value == nil {
		return 0
	}
	return *f.value
}

// ApplyPtr returns the given value, nil, or f's value, depending on whether f
// is unset, removed, or set, respectively.
func (f Float64) ApplyPtr(value *float64) *float64 {
	if f.set {
		return f.value
	}
	return value
}

// Diff returns f if f.Apply(value) != value; otherwise it returns an unset
// Float64. This can be used to omit extraneous updates when applying the update
// would have no effect.
func (f Float64) Diff(value float64) Float64 {
	if f.Apply(value) == value {
		return Float64{}
	}
	return f
}

// UnmarshalJSON implements json.Unmarshaler.
func (f *Float64) UnmarshalJSON(data []byte) error {
	f.set = true
	return json.Unmarshal(data, &f.value)
}

// IsSet returns true if the receiver has been set/removed.
func (f Float64) IsSet() bool {
	return f.set
}

// Removed returns whether the receiver has been removed (value set to nil).
func (f Float64) Removed() bool {
	return f.set && f.value == nil
}

// InterfaceValue returns value as an interface{}.
func (f Float64) InterfaceValue() interface{} {
	return f.value
}

// IsZero returns whether the receiver is set to 0.
func (f Float64) IsZero() bool {
	return f.set && f.value != nil && *f.value == 0.0
}

// IsNegative returns whether the receiver is set to a negative value.
func (f Float64) IsNegative() bool {
	return f.set && f.value != nil && *f.value < 0.0
}
