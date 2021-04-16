# nullable

This package provides `Nullable` field types, which enable distinguishing
between unset fields and fields explicitly set to null when marshalling to/from
JSON.

## Motivation

We define certain data updates using JSON objects, where each field in a
structure can be updated according to the following rules:

- If the field is present and non-`null`, the corresponding value is updated
- If the field is present but `null`, the corresponding field is removed
- If the field is absent, the corresponding value is left unmodified

We sometimes need to define go structs (for example the
[`EntityUpdate`](https://github.com/nicheinc/entity/blob/9c8bb0fe92e4e77e3af339c30b29fd122c190cd3/entity.go#L70-L78)
type in the `entity` service) corresponding to these updates, which need to be
marshalled to/from JSON.

If we were to use pointer fields with the `omitempty` JSON struct tag option for
these structs, then fields explicitly set to `nil` would simply be absent from
the marshalled JSON. If we were to use pointer fields _without_ `omitempty`,
then unset fields would be present and `null` in the JSON output.

The `Nullable` field types distinguish between "set to `nil`" and "not set",
allowing them to correctly and seamlessly unmarshal themselves from JSON.

## Marshalling

Unfortunately, the default marshaller is unaware of our `Nullable` types, and
providing a `MarshalJSON` implementation in the types themselves is
insufficient because it's the containing struct that determines which field
names appear in the JSON output. A custom implementation can use an ad-hoc
struct mirroring the original struct (but with an extra level of indirection),
along with a check per field that the field is set before copying it into the
output struct, as seen
[here](https://github.com/nicheinc/entity/blob/9c8bb0fe92e4e77e3af339c30b29fd122c190cd3/entity.go#L148-L167).

To avoid the need to define `MarshalJSON` for each struct containing `Nullable`
fields, this package provides the `nullable.MarshalJSON` function, which
implements a version of `json.Marshal` that respects the unset/removed status
of `Nullable` types.

`nullable.MarshalJSON` should behave exactly like
[`json.Marshal`](https://golang.org/pkg/encoding/json/#Marshal), with a few
exceptions:

- Anonymous fields are skipped
- The `string` tag option is ignored

To avoid accidentally calling the default implementation, it may be prudent to
implement a `MarshalJSON` for relevant types that simply calls
`nullable.MarshalJSON`.

## Future Improvements

Ideally we wouldn't have to maintain our own reflection-based marshalling code
to solve this problem. There is
[a proposal](https://github.com/nicheinc/entity/pull/120#discussion_r610908706)
to add support for omitting zero-valued structs with `omitempty`. If this or a
similar proposal were accepted, we could get rid of `nullable.MarshalJSON` and
implement `MarshalJSON` and `IsZero` for the `Nullable` types.

There is also a lot of code duplication across `Nullable` types. Once generics
are available, we can refactor the `Nullable` interface into a proper generic
type. Note that currently `nullable.StringSlice` has a different interface from
the other types because it stores a slice rather than a pointer, so we may
actually want two generic types: one for value-like types and one for
pointer-like types.
