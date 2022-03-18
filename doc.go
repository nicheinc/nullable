/*
Package nully provides types representing updates to struct fields,
distinguishing between no-ops, removals, and modifications when marshalling
updates to/from JSON.

Motivation

It's often useful to define data updates using JSON objects, where each
key-value pair represents an update to a field - using null to indicate
deletion. If a certain key is not present, the corresponding field is not
modified. We want to define go structs corresponding to these updates, which
need to be marshalled to/from JSON.

If we were to use pointer fields with the omitempty JSON struct tag option for
these structs, then fields explicitly set to nil to be removed would instead
simply be absent from the marshalled JSON, i.e. unchanged. If we were to use
pointer fields without omitempty, then nil fields would be present and null in
the JSON output, i.e. removed.

The Update and SliceUpdate types distinguish between no-op and removal updates,
allowing them to correctly and seamlessly unmarshal themselves from JSON.

Marshalling

Unfortunately, the default JSON marshaller is unaware of nully update types, and
providing a MarshalJSON implementation in the types themselves is insufficient
because it's the containing struct that determines which field names appear in
the JSON output.

A custom implementation can use an ad-hoc struct mirroring the original struct
(but with an extra level of indirection), along with a check per field that the
field is set before copying it into the output struct, but implementing this
method for every update type is laborious and error-prone. To avoid this
boilerplate, this package provides the nully.MarshalJSON function, which
implements a version of json.Marshal that respects the no-op/remove distinction.

Besides its treatment of Update/SliceUpdate fields, nully.MarshalJSON behaves
exactly like json.Marshal (https://golang.org/pkg/encoding/json/#Marshal), with
the following exceptions:

• Anonymous fields are skipped

• The string tag option is ignored

Note that the omitempty option does not affect nully update types. The default
JSON marshaller never omits struct values, but nully.MarshalJSON takes the use
of a nully update type per se as an indication to omit the field if it's a
no-op, even if omitempty is absent.

To avoid accidentally calling the default implementation, it's prudent to
implement for each relevant type a MarshalJSON method that simply calls
nully.MarshalJSON.

There are several outstanding golang proposals that could eliminate the need for
a custom MarshalJSON implementation in the future. One proposal
(https://github.com/golang/go/issues/11939) would allow zero-valued structs to
be treated as empty with respect to omitempty fields. Another proposal
(https://github.com/golang/go/issues/50480) would allow types to return (nil,
nil) from MarshalJSON to indicate they should be treated as empty by omitempty.
*/
package nully
