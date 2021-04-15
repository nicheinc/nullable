package nullable

import "encoding/json"

type Int struct {
	Set   bool
	Value *int
}

func (i *Int) SetValue(value int) {
	i.SetPtr(&value)
}

func (i *Int) SetPtr(value *int) {
	i.Set = true
	i.Value = value
}

func (i *Int) UnmarshalJSON(data []byte) error {
	i.Set = true
	return json.Unmarshal(data, &i.Value)
}

func (i *Int) Removed() bool {
	return i.Set && i.Value == nil
}

func (i *Int) IsZero() bool {
	return i.Set && i.Value != nil && *i.Value == 0
}

func (i *Int) IsNegative() bool {
	return i.Set && i.Value != nil && *i.Value < 0
}
