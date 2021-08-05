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

func (i Float64) interfaceValue() interface{} {
	return i.value
}

func (i Float64) IsZero() bool {
	return i.set && i.value != nil && *i.value == 0.0
}

func (i Float64) IsNegative() bool {
	return i.set && i.value != nil && *i.value < 0.0
}
