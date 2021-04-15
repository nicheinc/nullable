package nullable

import "encoding/json"

type Bool struct {
	Set   bool
	Value *bool
}

func (b *Bool) SetValue(value bool) {
	b.SetPtr(&value)
}

func (b *Bool) SetPtr(value *bool) {
	b.Set = true
	b.Value = value
}

func (b *Bool) UnmarshalJSON(data []byte) error {
	b.Set = true
	return json.Unmarshal(data, &b.Value)
}

func (b *Bool) Removed() bool {
	return b.Set && b.Value == nil
}
