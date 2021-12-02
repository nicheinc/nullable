/*
Package nullable provides types representing updates to struct fields,
distinguishing between no-ops, removals, and modifications when marshalling
updates to/from JSON.

Motivation

It's often useful to define data updates using JSON objects, where each
key-value pair represents a field and its new value - using null to indicate
deletion. If a certain key is not present, the corresponding field is not
modified. We want to define go structs corresponding to these updates, which
need to be marshalled to/from JSON.

If we were to use pointer fields with the omitempty JSON struct tag option for
these structs, then fields explicitly set to nil to be removed would simply be
absent from the marshalled JSON, i.e. unchanged. If we were to use pointer
fields without omitempty, then unset fields would be present and null in the
JSON output, i.e. removed.

The Nullable field types distinguish between "unchanged" and "removed", allowing
them to correctly and seamlessly unmarshal themselves from JSON.

Marshalling

Unfortunately, the default JSON marshaller is unaware of Nullable types, and
providing a MarshalJSON implementation in the types themselves is insufficient
because it's the containing struct that determines which field names appear in
the JSON output.

A custom implementation can use an ad-hoc struct mirroring the original struct
(but with an extra level of indirection), along with a check per field that the
field is set before copying it into the output struct, but implementing this
method for every update type is laborious and error-prone. To avoid this
boilerplate, this package provides the nullable.MarshalJSON function, which
implements a version of json.Marshal that respects the unset/removed status of
Nullable types.

Aside from Nullable fields, nullable.MarshalJSON should behave exactly like
json.Marshal (https://golang.org/pkg/encoding/json/#Marshal), with the following
exceptions:

• Anonymous fields are skipped

• The string tag option is ignored

Note that the omitempty option does not affect Nullable types. The default JSON
marshaller never omits struct values, but nullable.MarshalJSON takes the use of
a Nullable type per se as an indication to omit the field if it's unset, even if
omitempty is absent.

To avoid accidentally calling the default implementation, it's prudent to
implement a MarshalJSON for each relevant type that simply calls
nullable.MarshalJSON.

There is a proposal (https://github.com/golang/go/issues/11939) to add support
for omitting zero-valued structs with omitempty. If this or a similar proposal
were accepted, nullable.MarshalJSON could be eliminated in favor of implementing
MarshalJSON and IsZero for the Nullable types.
*/
package nullable
