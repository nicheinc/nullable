# nup

[![Build Status](https://github.com/nicheinc/nullable/actions/workflows/ci.yml/badge.svg)](https://github.com/nicheinc/nullable/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/nicheinc/nullable)](https://goreportcard.com/report/github.com/nicheinc/nullable)
[![Godoc](https://godoc.org/github.com/nicheinc/nullable?status.svg)](https://godoc.org/github.com/nicheinc/nullable)
[![license](https://img.shields.io/github/license/nicheinc/nullable.svg?cacheSeconds=2592000)](LICENSE)

This package provides types representing updates to struct fields,
distinguishing between no-ops, removals, and modifications when marshalling
those updates to/from JSON.

See [godoc](https://pkg.go.dev/github.com/nicheinc/nullable/v2) for usage and
examples.

## Motivation

It's often useful to define data updates using JSON objects, where each
key-value pair represents an update to a field, using null to indicate deletion.
If a certain key is not present, the corresponding field is not modified. We
want to define go structs corresponding to these updates, which need to be
marshalled to/from JSON.

If we were to use pointer fields with the `omitzero` JSON struct tag option for
these structs, then fields explicitly set to `nil` to be removed would instead
simply be absent from the marshalled JSON, i.e. unchanged. If we were to use
pointer fields without `omitzero`, then `nil` fields would be present and `null`
in the JSON output, i.e. removed.

The `nup.Update` and `nup.SliceUpdate` types distinguish between no-op and
removal updates, allowing them to correctly and seamlessly unmarshal themselves
from JSON.

## Marshalling

For best results, use
[`json.Marshal`](https://pkg.go.dev/encoding/json#Marshal)'s `omitzero` struct
tag option on all struct fields of type `nup.Update` or `nup.SliceUpdate`. This
will ensure that if the field is a no-op, it's correctly omitted from the JSON
output. (If the `omitzero` tag is absent, the field will be marshalled as
`null`.)

## Installation

This package can be imported into a module-aware Go project as follows:

`go get github.com/nicheinc/nullable/v2`

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for details on contributing to the `nup`
package.
