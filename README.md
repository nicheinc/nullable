# nullable

This package provides types representing updates to struct fields,
distinguishing between no-ops, removals, and modifications when marshalling
updates to/from JSON.

See [godoc](https://pkg.go.dev/github.com/nicheinc/nullable) for usage and
examples.

## Motivation

It's often useful to define data updates using JSON objects, where each
key-value pair represents a field and its new value - using `null` to indicate
deletion. If a certain key is not present, the corresponding field is not
modified. We want to define go structs corresponding to these updates, which
need to be marshalled to/from JSON.

If we were to use pointer fields with the `omitempty` JSON struct tag option for
these structs, then fields explicitly set to `nil` to be removed would simply be
absent from the marshalled JSON, i.e. unchanged. If we were to use pointer
fields without `omitempty`, then unset fields would be present and `null` in the
JSON output, i.e. removed.

The `Nullable` field types distinguish between "unchanged" and "removed",
allowing them to correctly and seamlessly unmarshal themselves from JSON.
