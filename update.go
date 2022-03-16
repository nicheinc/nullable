package nullable

import (
	"encoding/json"
	"fmt"
)

// Update represents an update to a value. The update may represent a null
// update, a removal, or a change in value.
type Update[T comparable] struct {
	set   bool
	value *T
}

// NewUpdate returns an Update set to the given value.
func NewUpdate[T comparable](value T) Update[T] {
	return Update[T]{
		set:   true,
		value: &value,
	}
}

// NewUpdatePtr returns an Update set to the given pointer.
func NewUpdatePtr[T comparable](value *T) Update[T] {
	return Update[T]{
		set:   true,
		value: value,
	}
}

// SetValue modifies the receiver to be an update to the given value.
func (u *Update[T]) SetValue(value T) {
	u.SetPtr(&value)
}

// SetPtr modifies the receiver to be an update to the given value. If the value
// is nil, the receiver will be removed.
func (u *Update[T]) SetPtr(value *T) {
	u.set = true
	u.value = value
}

// Value returns nil if the receiver is unset/removed or else the updated value.
func (u Update[T]) Value() *T {
	return u.value
}

// Apply returns the given value, the zero value of T, or the receiver's value,
// depending on whether the receiver is unset, removed, or set, respectively.
func (u Update[T]) Apply(value T) T {
	if !u.set {
		return value
	}
	if u.value == nil {
		var zeroValue T
		return zeroValue
	}
	return *u.value
}

// ApplyPtr returns the given value, nil, or the receiver's value, depending on
// whether the receiver is unset, removed, or set, respectively.
func (u Update[T]) ApplyPtr(value *T) *T {
	if u.set {
		return u.value
	}
	return value
}

// Diff returns the update u if u.Apply(value) != value; otherwise it returns an
// unset update. This can be used to omit extraneous updates when applying the
// update would have no effect.
func (u Update[T]) Diff(value T) Update[T] {
	if u.Apply(value) == value {
		return Update[T]{}
	}
	return u
}

// UnmarshalJSON implements json.Unmarshaler.
func (u *Update[T]) UnmarshalJSON(data []byte) error {
	u.set = true
	return json.Unmarshal(data, &u.value)
}

// IsSet returns true if the receiver has been set/removed.
func (u Update[T]) IsSet() bool {
	return u.set
}

// Removed returns whether the receiver has been removed (value set to nil).
func (u Update[T]) Removed() bool {
	return u.set && u.value == nil
}

// IsSetTo returns whether the update is a change to the given new value.
func (u Update[T]) IsSetTo(newValue T) bool {
	return u.value != nil && *u.value == newValue
}

// IsSetSuchThat returns whether the update is a change whose new value
// satisfies the given predicate.
func (u Update[T]) IsSetSuchThat(predicate func(T) bool) bool {
	return u.value != nil && predicate(*u.value)
}

// String implements fmt.Stringer. It returns "<unset>", "<removed>", or a
// string representation of the updated value.
func (u Update[T]) String() string {
	if u.Removed() {
		return "<removed>"
	}
	if !u.IsSet() {
		return "<unset>"
	}
	switch value := interface{}(*u.value).(type) {
	case string:
		return value
	case fmt.Stringer:
		return value.String()
	default:
		return fmt.Sprintf("%v", value)
	}
}

// include partially implements updateMarshaller.
func (u Update[T]) include() bool {
	return u.set
}

// interfaceValue partially implements updateMarshaller.
func (u Update[T]) interfaceValue() interface{} {
	return u.value
}
