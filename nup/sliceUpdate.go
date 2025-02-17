package nup

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

// SliceNoop returns a slice update that does nothing. This is equivalent to the
// zero-valued SliceUpdate.
func SliceNoop[T comparable]() SliceUpdate[T] {
	return SliceUpdate[T]{
		op: OpNoop,
	}
}

// SliceRemove returns a slice update that removes a field (sets it to nil).
func SliceRemove[T comparable]() SliceUpdate[T] {
	return SliceUpdate[T]{
		op: OpRemove,
	}
}

// SliceRemoveOrSet returns a slice update that either removes or sets a field's
// value, depending on the given slice value. If the value is nil, it will
// remove; otherwise it will set to the given value. Note that a nil slice is
// different from an allocated but zero-length slice, []T{}.
func SliceRemoveOrSet[T comparable](value []T) SliceUpdate[T] {
	if value == nil {
		return SliceUpdate[T]{
			op: OpRemove,
		}
	}
	return SliceUpdate[T]{
		op:    OpSet,
		value: value,
	}
}

// ValueOperation returns a shallow copy of the value this update sets fields to
// (if any) and the operation this update performs: no-op, remove, or set. If
// this update is not a set operation, then the returned value is always nil;
// i.e., the value is only meaningful if the operation is OpSet.
func (u SliceUpdate[T]) ValueOperation() (value []T, operation Operation) {
	return u.value, u.op
}

// Operation returns the operation this update performs: no-op, remove, or set.
func (u SliceUpdate[T]) Operation() Operation {
	return u.op
}

// IsNoop returns whether this update is a no-op. IsNoop is equivalent to
// Operation() == OpNoop.
func (u SliceUpdate[T]) IsNoop() bool {
	return u.op == OpNoop
}

// IsZero is equivalent to IsNoop.
func (u SliceUpdate[T]) IsZero() bool {
	return u.IsNoop()
}

// IsRemove returns whether this update is a remove operation. IsRemove is
// equivalent to Operation() == OpRemove.
func (u SliceUpdate[T]) IsRemove() bool {
	return u.op == OpRemove
}

// IsSet returns whether this update is a set operation. IsSet is equivalent to
// Operation() == OpSet.
func (u SliceUpdate[T]) IsSet() bool {
	return u.op == OpSet
}

// IsChange returns whether this update is either a set or remove operation
// (i.e., not a no-op). IsChange is equivalent to Operation() != OpNoop.
func (u SliceUpdate[T]) IsChange() bool {
	return u.op != OpNoop
}

// Value returns a shallow copy of the value this update sets fields to (if any)
// and an isSet flag indicating whether the update is a set operation. If the
// flag is false (because the update is actually a no-op or removal), then the
// returned value is nil.
func (u SliceUpdate[T]) Value() (value []T, isSet bool) {
	return u.value, u.op == OpSet
}

// ValueOrNil returns a shallow copy of this update's value if it's a set
// operation or else nil.
func (u SliceUpdate[T]) ValueOrNil() []T {
	return u.value
}

// Apply returns the result of applying the update to the given value. The
// result is the given value if the update is a no-op, nil if it's a removal, or
// a shallow copy of the update's contained value if it's a set operation.
func (u SliceUpdate[T]) Apply(value []T) []T {
	switch u.op {
	case OpNoop:
		return value
	case OpRemove:
		return nil
	default: // Set
		return u.value
	}
}

// Diff returns the update itself if Apply(value) is not element-wise equal to
// value; otherwise it returns a no-op update. Diff can be used to omit
// extraneous updates when applying them would have no effect.
func (u SliceUpdate[T]) Diff(value []T) SliceUpdate[T] {
	if sliceEquals(u.Apply(value), value) {
		return SliceNoop[T]()
	}
	return u
}

// MarshalJSON implements json.Marshaler.
func (u SliceUpdate[T]) MarshalJSON() ([]byte, error) {
	if u.op == OpSet {
		return json.Marshal(u.value)
	}
	return []byte("null"), nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (u *SliceUpdate[T]) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		u.op = OpRemove
		return nil
	}
	u.op = OpSet
	return json.Unmarshal(data, &u.value)
}

// IsSetTo returns whether the update sets to a value that is element-wise equal
// to the given value.
func (u SliceUpdate[T]) IsSetTo(value []T) bool {
	return u.op == OpSet && sliceEquals(u.value, value)
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
	return u.op == OpSet && predicate(u.value)
}

// String implements fmt.Stringer. It returns "<unset>", "<removed>", or a
// string representation of the updated value.
func (u SliceUpdate[T]) String() string {
	switch u.op {
	case OpNoop:
		return "<no-op>"
	case OpRemove:
		return "<remove>"
	}
	return fmt.Sprintf("%v", u.value)
}

// Equal returns whether u and other perform the same type of operation and, if
// both are set operations, set to values that are element-wise equal using the
// == operator. This method is a quasi-standard mechanism to define custom
// equality. For instance, the time package defines a similar method
// (https://pkg.go.dev/github.com/google/go-cmp/cmp#Equal), and
// https://github.com/google/go-cmp respects methods of this form.
func (u SliceUpdate[T]) Equal(other SliceUpdate[T]) bool {
	return u.op == other.op && sliceEquals(u.value, other.value)
}

// interfaceValue, along with IsChange, implements updateMarshaller, which
// nup.MarshalJSON uses to detect update types and marshal them correctly.
func (u SliceUpdate[T]) interfaceValue() interface{} {
	if u.op == OpSet {
		return u.value
	}
	return nil
}
