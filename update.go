package nullable

import (
	"encoding/json"
	"fmt"
)

// Update represents an update that can be applied to a value field. It may set,
// remove, or have no effect on a field's value. For updates to slice fields,
// see SliceUpdate.
type Update[T comparable] struct {
	op    Operation
	value T
}

// NewNoop returns an update that does nothing. This is equivalent to the
// zero-valued Update.
func NewNoop[T comparable]() Update[T] {
	return Update[T]{
		op: Noop,
	}
}

// NewRemove returns an update that removes a field (sets it to the zero value).
func NewRemove[T comparable]() Update[T] {
	return Update[T]{
		op: Remove,
	}
}

// NewSet returns an update that sets a field's value to the given value.
func NewSet[T comparable](value T) Update[T] {
	return Update[T]{
		op:    Set,
		value: value,
	}
}

// Operation returns the operation this update performs: no-op, remove, or set.
func (u Update[T]) Operation() Operation {
	return u.op
}

// IsNoop is shorthand for Operation() == Noop.
func (u Update[T]) IsNoop() bool {
	return u.op == Noop
}

// IsRemove is shorthand for Operation() == Remove.
func (u Update[T]) IsRemove() bool {
	return u.op == Remove
}

// IsSet is shorthand for Operation() == Set.
func (u Update[T]) IsSet() bool {
	return u.op == Set
}

// Value returns the update's value and an "ok" flag indicating whether the
// update is a set operation. If the flag is false (because the update is
// actually a no-op or removal), then the returned value is T's zero value.
func (u Update[T]) Value() (value T, ok bool) {
	return u.value, u.op == Set
}

// Apply returns the result of applying the update to the given value. The
// result is the given value if the update is a no-op, the zero value if it's a
// removal, or the update's contained value is if it's a set operation.
func (u Update[T]) Apply(value T) T {
	switch u.op {
	case Noop:
		return value
	case Remove:
		var zero T
		return zero
	default: // Set
		return u.value
	}
}

// Diff returns the update itself if Apply(value) != value; otherwise it returns
// a no-op update. Diff can be used to omit extraneous updates when applying
// them would have no effect.
func (u Update[T]) Diff(value T) Update[T] {
	if u.Apply(value) == value {
		return NewNoop[T]()
	}
	return u
}

// UnmarshalJSON implements json.Unmarshaler.
func (u *Update[T]) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		u.op = Remove
		return nil
	}
	u.op = Set
	return json.Unmarshal(data, &u.value)
}

// IsSetTo returns whether the update sets to the given value.
func (u Update[T]) IsSetTo(value T) bool {
	return u.op == Set && u.value == value
}

// IsSetSuchThat returns whether the update is a set operation to a value that
// satisfies the given predicate.
func (u Update[T]) IsSetSuchThat(predicate func(T) bool) bool {
	return u.op == Set && predicate(u.value)
}

// String implements fmt.Stringer. It returns "<unset>", "<removed>", or a
// string representation of the updated value.
func (u Update[T]) String() string {
	switch u.op {
	case Noop:
		return "<no-op>"
	case Remove:
		return "<remove>"
	}
	switch value := interface{}(u.value).(type) {
	case string:
		return value
	case fmt.Stringer:
		return value.String()
	default:
		return fmt.Sprintf("%v", value)
	}
}

// shouldBeMarshalled partially implements updateMarshaller.
func (u Update[T]) shouldBeMarshalled() bool {
	return u.op != Noop
}

// interfaceValue partially implements updateMarshaller.
func (u Update[T]) interfaceValue() interface{} {
	if u.op == Set {
		return u.value
	}
	return nil
}
