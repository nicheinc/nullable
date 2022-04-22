# Version 2 Features

Version 2 replaces the `Nullable` interface and various implementations
(including user-defined implementations!) with just two generic types:
`Update[T]` and `SliceUpdate[T]`.

`Update` uses a value internally, rather than a pointer. This may reduce the
need for allocations, and it also means that updates can now be meaningfully
compared using `==`.

Unlike `StringSlice`, `SliceUpdate` can distinguish between removed and set to
`nil`. This is more in line with non-slice updates, which could always
distinguish between removed and set to the zero value.

In `nullable` version 1, only `nullable.String` implements `fmt.Stringer`. With
go 1.18 generics, it's not possible to implement an interface only for
particular type parameters of a generic type, leaving us with the choice to
implement `fmt.Stringer` for all types or for none. We chose the former. The
`Update` implementation uses a type switch to check if `T` is itself a `string`
or `fmt.Stringer`, falling back to `fmt.Sprintf` with the `%v` verb for all
other types. (The `SliceUpdate[T]` implementation always uses `fmt.Sprintf` with
`%v` since a `[]T` is not `string` or `fmt.Stringer`.)

# Migration Guide

The package has been renamed from `nullable` to `nup` (for "nullable update").
It's 62.5% shorter and 87% more whimsical!

The old `Value` method on each `Nullable` type has been replaced with three
`Update` methods: `Value`, `ValueOperation`, and `ValueOrNil`:
- The new `Value` method returns two values: `value` and `isSet`. `isSet`
  indicates whether the update is a set operation. If `isSet` is true, then
  `value` is the update's set value. If `isSet` is false, then `value` is the
  zero value for the update's type parameter.
- `ValueOperation` is similar to the new `Value` method but returns `value` and
  `operation`. `value` is the zero value unless `operation == OpSet`. This
  method is useful for handling all three update types in one switch statement,
  for example.
- `ValueOrNil` behaves similarly to the old `Value` methods, returning a value
  pointer that is non-`nil` only if the update is a set operation. However,
  unlike the old methods, the returned pointer references a fresh copy of the
  contained value, so it can't be used to mutate the update.

The following table illustrates some of the other noteworthy differences between
version 1 and 2:

| Version 1               | Version 2                                            |
| :---------------------- | :--------------------------------------------------- |
| `nullable.Int`          | `nup.Update[int]`                                    |
| `nullable.NewIntPtr(p)` | `nup.RemoveOrSet(p)`                                 |
| `u.IsSet()`             | `u.IsChange()`                                       |
| `u.Removed()`           | `u.IsRemove()`                                       |
| `u.Equals(5)`           | `u.IsSetTo(5)`                                       |
| `u.IsNegative()`        | `u.IsSetSuchThat(func(v int) bool { return v < 0 })` |

The semantics of `IsSet()` have changed from the version 1 types. Previously,
`IsSet()` returned true if the update was a set _or removal_ update. In version
2, a given update is exactly one of three operations: no-op, remove, or set.
Thus `IsSet()` is now mutually exclusive to `IsRemove` (`Removed()`, in version
1). Suppose you had a variable `update` of type `nullable.Int`, and you used
`update.IsSet()` to check that `update` is not the zero value. Then you should
_not_ only change `update`'s type to `nup.Update[int]` because `IsSet()` will no
longer return true if `update` is a remove operation. You should instead use
`update.IsChange()`.

The following update-type-specific value checks have been removed: `Int.IsZero`,
`Int.IsNegative`, `Float64.IsZero`, `Float64.IsNegative`, `String.IsEmpty`, and
`StringSlice.IsEmpty`. If go generics are ever extended to support type
constraints on generic methods, we could restore some of these methods. For now,
use `IsSetTo` or `IsSetSuchThat` with a suitable argument (as shown in the table
above). Another option is to implement the removed type-specific methods you
care about as free functions, for instance:

```go
func IsNegative(update nup.Update[int]) bool {
	value, ok := update.Value()
	return ok && value < 0
}
```

There is a small change to the behavior of `nup.MarshalJSON` for types
containing `nup` types nested deeper than one level. Previously, our
`MarshalJSON` only performed special marshalling logic for top-level struct
fields. Now it recurses into structs. Consider the following snippet:

```go
type Child struct {
	Update nup.Update[int]
}
type Parent struct {
	Child Child
}
value := Parent{
	Child: Child{
		Update: nup.Noop[int]()
	}
}
```

Previously (ignoring v1 vs. v2 update types), `value` would have been marshalled
as `{"Child":{"Update":null}}`. Now, it will be marshalled as `{"Child":{}}`.

Note that `nup.MarshalJSON` still does not recurse into `map`s or other types
that might contain a `nup` type.
