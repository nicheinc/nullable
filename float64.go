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

func (f *Float64) SetPtr(value *float64) {
	f.set = true
	f.value = value
}

func (f Float64) Value() *float64 {
	return f.value
}

func (f *Float64) UnmarshalJSON(data []byte) error {
	f.set = true
	return json.Unmarshal(data, &f.value)
}

func (f Float64) IsSet() bool {
	return f.set
}

func (f Float64) Removed() bool {
	return f.set && f.value == nil
}

func (f Float64) InterfaceValue() interface{} {
	return f.value
}

func (f Float64) IsZero() bool {
	return f.set && f.value != nil && *f.value == 0.0
}

func (f Float64) IsNegative() bool {
	return f.set && f.value != nil && *f.value < 0.0
}

// Scan implements the sql.Scanner interface (https://pkg.go.dev/database/sql#Scanner).
func (f *Float64) Scan(src interface{}) error {
	switch value := src.(type) {
	case nil:
		f.SetPtr(nil)
	case float64:
		f.SetValue(value)
	case int64:
		f.SetValue(float64(value))
	default:
		return &ScanTypeError{
			Src:  src,
			Dest: f,
		}
	}
	return nil
}
