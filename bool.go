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

// Scan implements the sql.Scanner interface (https://pkg.go.dev/database/sql#Scanner).
func (b *Bool) Scan(src interface{}) error {
	switch value := src.(type) {
	case nil:
		*b = Bool{}
	case bool:
		b.SetValue(value)
	default:
		return &ScanTypeError{
			Src:  src,
			Dest: b,
		}
	}
	return nil
}
