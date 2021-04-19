package nullable

import "encoding/json"

type Int struct {
	set   bool
	value *int
}

// NewInt returns an Int set to the given value.
func NewInt(v int) Int {
	return Int{
		set:   true,
		value: &v,
	}
}

// NewIntPtr returns an Int set to the given pointer.
func NewIntPtr(v *int) Int {
	return Int{
		set:   true,
		value: v,
	}
}

func (i *Int) SetValue(value int) {
	i.SetPtr(&value)
}

func (i *Int) SetPtr(value *int) {
	i.set = true
	i.value = value
}

func (i Int) Value() *int {
	return i.value
}

func (i *Int) UnmarshalJSON(data []byte) error {
	i.set = true
	return json.Unmarshal(data, &i.value)
}

func (i Int) IsSet() bool {
	return i.set
}

func (i Int) Removed() bool {
	return i.set && i.value == nil
}

func (i Int) interfaceValue() interface{} {
	return i.value
}

func (i Int) IsZero() bool {
	return i.set && i.value != nil && *i.value == 0
}

func (i Int) IsNegative() bool {
	return i.set && i.value != nil && *i.value < 0
}
