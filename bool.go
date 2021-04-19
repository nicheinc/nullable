package nullable

import "encoding/json"

type Bool struct {
	set   bool
	value *bool
}

// MakeBool returns a Bool set to the given value.
func MakeBool(v bool) Bool {
	return Bool{
		set:   true,
		value: &v,
	}
}

// MakeBoolPtr returns a Bool set to the given pointer.
func MakeBoolPtr(v *bool) Bool {
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

func (b Bool) interfaceValue() interface{} {
	return b.value
}
