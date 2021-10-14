package nullable

import "encoding/json"

type Float64 struct {
	set   bool
	value *float64
}

// NewInt returns a Float64 set to the given value.
func NewFloat64(v float64) Float64 {
	return Float64{
		set:   true,
		value: &v,
	}
}

// NewIntPtr returns a Float64 set to the given pointer.
func NewFloat64Ptr(v *float64) Float64 {
	return Float64{
		set:   true,
		value: v,
	}
}

func (i *Float64) SetValue(value float64) {
	i.SetPtr(&value)
}

func (i *Float64) SetPtr(value *float64) {
	i.set = true
	i.value = value
}

func (i Float64) Value() *float64 {
	return i.value
}

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

// Diff returns f if f.Apply(value) != value or else an unset Float64. This can
// be used to avoid extraneous updates when the update would have no effect.
func (f Float64) Diff(value float64) Float64 {
	if f.Apply(value) == value {
		return Float64{}
	}
	return f
}

func (i *Float64) UnmarshalJSON(data []byte) error {
	i.set = true
	return json.Unmarshal(data, &i.value)
}

func (i Float64) IsSet() bool {
	return i.set
}

func (i Float64) Removed() bool {
	return i.set && i.value == nil
}

func (i Float64) InterfaceValue() interface{} {
	return i.value
}

func (i Float64) IsZero() bool {
	return i.set && i.value != nil && *i.value == 0.0
}

func (i Float64) IsNegative() bool {
	return i.set && i.value != nil && *i.value < 0.0
}
