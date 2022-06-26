package nup

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

// Noop returns an update that does nothing. This is equivalent to the
// zero-valued Update.
func Noop[T comparable]() Update[T] {
	return Update[T]{
		op: OpNoop,
	}
}

// Remove returns an update that removes a field (sets it to the zero value).
func Remove[T comparable]() Update[T] {
	return Update[T]{
		op: OpRemove,
	}
}

// Set returns an update that sets a field's value to the given value.
func Set[T comparable](value T) Update[T] {
	return Update[T]{
		op:    OpSet,
		value: value,
	}
}

// RemoveOrSet returns an update that either removes or sets a field's value,
// depending on the given pointer. If the pointer is nil, it will remove;
// otherwise it will set to the pointer's value.
func RemoveOrSet[T comparable](ptr *T) Update[T] {
	if ptr == nil {
		return Remove[T]()
	}
	return Set(*ptr)
}

// ValueOperation returns the value this update sets fields to (if any) and the
// operation this update performs: no-op, remove, or set. If this update is not
// a set operation, then the returned value is T's zero value; i.e., the value
// is only meaningful if the operation is OpSet.
func (u Update[T]) ValueOperation() (value T, operation Operation) {
	return u.value, u.op
}

// Operation returns the operation this update performs: no-op, remove, or set.
func (u Update[T]) Operation() Operation {
	return u.op
}

// IsNoop returns whether this update is a no-op. IsNoop is equivalent to
// Operation() == OpNoop.
func (u Update[T]) IsNoop() bool {
	return u.op == OpNoop
}

// IsRemove returns whether this update is a remove operation. IsRemove is
// equivalent to Operation() == OpRemove.
func (u Update[T]) IsRemove() bool {
	return u.op == OpRemove
}

// IsSet returns whether this update is a set operation. IsSet is equivalent to
// Operation() == OpSet.
func (u Update[T]) IsSet() bool {
	return u.op == OpSet
}

// IsChange returns whether this update is either a set or remove operation
// (i.e., not a no-op). IsChange is equivalent to Operation() != OpNoop.
func (u Update[T]) IsChange() bool {
	return u.op != OpNoop
}

// Value returns the value this update sets fields to (if any) and an isSet flag
// indicating whether the update is a set operation. If the flag is false
// (because the update is actually a no-op or removal), then the returned value
// is T's zero value.
func (u Update[T]) Value() (value T, isSet bool) {
	return u.value, u.op == OpSet
}

// ValueOrNil returns this update's value if it's a set operation or else nil.
func (u Update[T]) ValueOrNil() *T {
	if u.op != OpSet {
		return nil
	}
	// Copy the update value so it can't be mutated via the returned pointer.
	value := u.value
	return &value
}

// Apply returns the result of applying the update to the given value. The
// result is the given value if the update is a no-op, the zero value if it's a
// removal, or the update's contained value if it's a set operation.
func (u Update[T]) Apply(value T) T {
	switch u.op {
	case OpNoop:
		return value
	case OpRemove:
		var zero T
		return zero
	default: // Set
		return u.value
	}
}

// ApplyPtr returns the result of applying the update to the given pointer
// value. The result is the given value if the update is a no-op, nil if it's a
// removal, or a copy of the update's contained value if it's a set operation.
func (u Update[T]) ApplyPtr(value *T) *T {
	switch u.op {
	case OpNoop:
		return value
	case OpRemove:
		return nil
	default: // Set
		// Copy the update value so it can't be mutated via the returned pointer.
		value := u.value
		return &value
	}
}

// Diff returns the update itself if Apply(value) != value; otherwise it returns
// a no-op update. Diff can be used to omit extraneous updates when applying
// them would have no effect.
func (u Update[T]) Diff(value T) Update[T] {
	if u.Apply(value) == value {
		return Noop[T]()
	}
	return u
}

// DiffPtr returns the update itself if ApplyPtr(value) does not contain a value
// equal to the given value; otherwise it returns a no-op update. DiffPtr can be
// used to omit extraneous updates when applying them would have no effect.
func (u Update[T]) DiffPtr(value *T) Update[T] {
	applied := u.ApplyPtr(value)
	if applied == nil || value == nil {
		if applied == value {
			return Noop[T]()
		}
		return u
	}
	return u.Diff(*value)
}

// UnmarshalJSON implements json.Unmarshaler.
func (u *Update[T]) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		u.op = OpRemove
		return nil
	}
	u.op = OpSet
	return json.Unmarshal(data, &u.value)
}

// IsSetTo returns whether the update sets to the given value.
func (u Update[T]) IsSetTo(value T) bool {
	return u.op == OpSet && u.value == value
}

// IsSetSuchThat returns whether the update is a set operation to a value that
// satisfies the given predicate.
func (u Update[T]) IsSetSuchThat(predicate func(T) bool) bool {
	return u.op == OpSet && predicate(u.value)
}

// String implements fmt.Stringer. It returns "<unset>", "<removed>", or a
// string representation of the updated value.
func (u Update[T]) String() string {
	switch u.op {
	case OpNoop:
		return "<no-op>"
	case OpRemove:
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

// Equal compares u with other using the == operator. This method is a
// quasi-standard mechanism to define custom equality. For instance, the time
// package defines a similar method
// (https://pkg.go.dev/github.com/google/go-cmp/cmp#Equal), and
// https://github.com/google/go-cmp respects methods of this form.
func (u Update[T]) Equal(other Update[T]) bool {
	return u == other
}

// interfaceValue, along with IsChange, implements updateMarshaller, which
// nup.MarshalJSON uses to detect update types and marshal them correctly.
func (u Update[T]) interfaceValue() interface{} {
	if u.op == OpSet {
		return u.value
	}
	return nil
}
