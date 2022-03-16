package nullable

import (
	"encoding/json"
	"fmt"
)

// SliceUpdate represents an update to a slice field. It may set, remove, or
// have no effect on a field's value. For updates to value fields, see Update.
type SliceUpdate[T comparable] struct {
	op    Operation
	value []T
}

// NewNoopSlice returns an update that does nothing. This is equivalent to the
// zero-valued SliceUpdate.
func NewNoopSlice[T comparable]() SliceUpdate[T] {
	return SliceUpdate[T]{
		op: Noop,
	}
}

// NewRemoveSlice returns an update that removes a field (sets it to nil).
func NewRemoveSlice[T comparable]() SliceUpdate[T] {
	return SliceUpdate[T]{
		op: Remove,
	}
}

// NewSetSlice returns an update that sets a field's value to the given value.
func NewSetSlice[T comparable](value []T) SliceUpdate[T] {
	return SliceUpdate[T]{
		op:    Set,
		value: value,
	}
}

// Operation returns the operation this update performs: no-op, remove, or set.
func (u SliceUpdate[T]) Operation() Operation {
	return u.op
}

// IsNoop is shorthand for Operation() == Noop.
func (u SliceUpdate[T]) IsNoop() bool {
	return u.op == Noop
}

// IsRemove is shorthand for Operation() == Remove.
func (u SliceUpdate[T]) IsRemove() bool {
	return u.op == Remove
}

// IsSet is shorthand for Operation() == Set.
func (u SliceUpdate[T]) IsSet() bool {
	return u.op == Set
}

// Value returns the update's value and an "ok" flag indicating whether the
// update is a set operation. If the flag is false (because the update is
// actually a no-op or removal), then the returned value is nil.
func (u SliceUpdate[T]) Value() (value []T, ok bool) {
	return u.value, u.op == Set
}

// Apply returns the result of applying the update to the given value. The
// result is the given value if the update is a no-op, nil if it's a removal, or
// the update's contained value is if it's a set operation.
func (u SliceUpdate[T]) Apply(value []T) []T {
	switch u.op {
	case Noop:
		return value
	case Remove:
		return nil
	default: // Set
		return u.value
	}
}

// Diff returns the update itself if Apply(value) != value; otherwise it returns
// a no-op update. Diff can be used to omit extraneous updates when applying
// them would have no effect.
func (u SliceUpdate[T]) Diff(value []T) SliceUpdate[T] {
	if sliceEquals(u.Apply(value), value) {
		return NewNoopSlice[T]()
	}
	return u
}

// UnmarshalJSON implements json.Unmarshaler.
func (u *SliceUpdate[T]) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		u.op = Remove
		return nil
	}
	u.op = Set
	return json.Unmarshal(data, &u.value)
}

// IsSetTo returns whether the update set to the given value.
func (u SliceUpdate[T]) IsSetTo(newValue []T) bool {
	return u.op == Set && sliceEquals(u.value, newValue)
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

// IsSetSuchThat returns whether the update is a set operation to a value that
// satisfies the given predicate.
func (u SliceUpdate[T]) IsSetSuchThat(predicate func([]T) bool) bool {
	return u.op == Set && predicate(u.value)
}

// String implements fmt.Stringer. It returns "<unset>", "<removed>", or a
// string representation of the updated value.
func (u SliceUpdate[T]) String() string {
	switch u.op {
	case Noop:
		return "<no-op>"
	case Remove:
		return "<remove>"
	}
	return fmt.Sprintf("%v", u.value)
}

// shouldBeMarshalled partially implements updateMarshaller.
func (u SliceUpdate[T]) shouldBeMarshalled() bool {
	return u.op != Noop
}

// interfaceValue partially implements updateMarshaller.
func (u SliceUpdate[T]) interfaceValue() interface{} {
	if u.op == Set {
		return u.value
	}
	return nil
}
