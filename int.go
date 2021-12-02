package nullable

import "encoding/json"

// Int implements Nullable for int fields.
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

// SetValue modifies the receiver to be an update to the given value.
func (i *Int) SetValue(value int) {
	i.SetPtr(&value)
}

// SetPtr modifies the receiver to be an update to the given value. If the value
// is nil, the receiver will be removed.
func (i *Int) SetPtr(value *int) {
	i.set = true
	i.value = value
}

// Value returns nil if the receiver is unset/removed or else the updated value.
func (i Int) Value() *int {
	return i.value
}

// Equals returns whether the receiver is set to the given value.
func (i Int) Equals(value int) bool {
	return i.value != nil && *i.value == value
}

// Apply returns the given value, the zero value (0), or i's value, depending on
// whether i is unset, removed, or set, respectively.
func (i Int) Apply(value int) int {
	if !i.set {
		return value
	}
	if i.value == nil {
		return 0
	}
	return *i.value
}

// ApplyPtr returns the given value, nil, or i's value, depending on whether i
// is unset, removed, or set, respectively.
func (i Int) ApplyPtr(value *int) *int {
	if i.set {
		return i.value
	}
	return value
}

// Diff returns i if i.Apply(value) != value; otherwise it returns an unset Int.
// This can be used to omit extraneous updates when applying the update would
// have no effect.
func (i Int) Diff(value int) Int {
	if i.Apply(value) == value {
		return Int{}
	}
	return i
}

func (i *Int) UnmarshalJSON(data []byte) error {
	i.set = true
	return json.Unmarshal(data, &i.value)
}

// IsSet returns true if the receiver has been set/removed.
func (i Int) IsSet() bool {
	return i.set
}

// Removed returns whether the receiver has been removed (value set to nil).
func (i Int) Removed() bool {
	return i.set && i.value == nil
}

func (i Int) InterfaceValue() interface{} {
	return i.value
}

// IsZero returns whether the receiver is set to 0.
func (i Int) IsZero() bool {
	return i.set && i.value != nil && *i.value == 0
}

// IsNegative returns whether the receiver is set to a negative value.
func (i Int) IsNegative() bool {
	return i.set && i.value != nil && *i.value < 0
}
