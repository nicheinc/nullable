package nullable

import (
	"encoding/json"
	"fmt"
)

// SliceUpdate represents an update to a slice. The update may represent a null
// update, a removal, or a change in value.
type SliceUpdate[T comparable] struct {
	set   bool
	value []T
}

// NewSliceUpdate returns an Update set to the given value.
func NewSliceUpdate[T comparable](value []T) SliceUpdate[T] {
	return SliceUpdate[T]{
		set:   true,
		value: value,
	}
}

// SetValue modifies the receiver to be an update to the given value.
func (u *SliceUpdate[T]) SetValue(value []T) {
	u.set = true
	u.value = value
}

// Value returns nil if the receiver is unset/removed or else the updated value.
func (u SliceUpdate[T]) Value() []T {
	return u.value
}

// Apply returns the given value if the receiver is unset or else returns the
// receiver's value.
func (u SliceUpdate[T]) Apply(value []T) []T {
	if !u.set {
		return value
	}
	return u.value
}

// Diff returns the update u if u.Apply(value) != value; otherwise it returns an
// unset update. This can be used to omit extraneous updates when applying the
// update would have no effect.
func (u SliceUpdate[T]) Diff(value []T) SliceUpdate[T] {
	if sliceEquals(u.Apply(value), value) {
		return SliceUpdate[T]{}
	}
	return u
}

// UnmarshalJSON implements json.Unmarshaler.
func (u *SliceUpdate[T]) UnmarshalJSON(data []byte) error {
	u.set = true
	return json.Unmarshal(data, &u.value)
}

// IsSet returns true if the receiver has been set/removed.
func (u SliceUpdate[T]) IsSet() bool {
	return u.set
}

// Removed returns whether the receiver has been removed (value set to nil).
func (u SliceUpdate[T]) Removed() bool {
	return u.set && u.value == nil
}

// IsSetTo returns whether the update is a change to the given new value.
func (u SliceUpdate[T]) IsSetTo(newValue []T) bool {
	return u.value != nil && sliceEquals(u.value, newValue)
}

func sliceEquals[T comparable](slice1 []T, slice2 []T) bool {
	if len(slice1) != len(slice2) {
		return false
	}
	for i := range slice1 {
		if slice1[i] != slice2[i] {
			return false
		}
	}
	return true
}

// IsSetSuchThat returns whether the update is a change whose new value
// satisfies the given predicate.
func (u SliceUpdate[T]) IsSetSuchThat(predicate func([]T) bool) bool {
	return u.value != nil && predicate(u.value)
}

// String implements fmt.Stringer. It returns "<unset>", "<removed>", or a
// string representation of the updated value.
func (u SliceUpdate[T]) String() string {
	if u.Removed() {
		return "<removed>"
	}
	if !u.IsSet() {
		return "<unset>"
	}
	return fmt.Sprintf("%v", u.value)
}

// include partially implements updateMarshaller.
func (u SliceUpdate[T]) include() bool {
	return u.set
}

// interfaceValue partially implements updateMarshaller.
func (u SliceUpdate[T]) interfaceValue() interface{} {
	return u.value
}
