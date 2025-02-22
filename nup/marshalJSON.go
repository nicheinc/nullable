package nup

import (
	"encoding/json"
	"reflect"
	"strings"
)

// Deprecated: As of Go 1.24, json.Marshal handles fields of nup types correctly
// by default, as long as those fields are marked with the "omitzero" JSON
// struct tag.
//
// MarshalJSON is a reimplementation of json.Marshal that understands nup types.
// Any struct that contains Update or SliceUpdate fields should call this
// function instead of the default json.Marshal. For more info, see
// https://pkg.go.dev/github.com/nicheinc/nullable/#hdr-Marshalling.
func MarshalJSON(v interface{}) ([]byte, error) {
	// Marshal nil as null.
	if v == nil {
		return []byte("null"), nil
	}
	var (
		reflectedValue = reflect.ValueOf(v)
		reflectedType  = reflectedValue.Type()
	)
	// Dereference the reflected value once, if necessary, to support pointers.
	if reflectedType.Kind() == reflect.Pointer {
		reflectedValue = reflectedValue.Elem()
		reflectedType = reflectedValue.Type()
	}
	// Delegate non-struct values to the default implementation.
	if reflectedType.Kind() != reflect.Struct {
		return json.Marshal(v)
	}

	var (
		buf          = []byte{'{'}
		prependComma = false
	)
	for i := 0; i < reflectedValue.NumField(); i++ {
		var (
			field      = reflectedType.Field(i)
			fieldValue = reflectedValue.Field(i)
		)
		key := getKeyName(field, fieldValue)
		if key == nil {
			continue
		}
		appendField := func(fieldInterface interface{}) error {
			valueBuf, err := json.Marshal(fieldInterface)
			if err != nil {
				return err
			}
			// Allocate space for the quoted key, a colon, and the value, as
			// well as a leading comma if this isn't the first field marshalled.
			capacity := len(*key) + len(valueBuf) + 3
			if prependComma {
				capacity++
			}
			fieldBuf := make([]byte, 0, capacity)
			// Append the components of the field data.
			if prependComma {
				fieldBuf = append(fieldBuf, `,"`...)
			} else {
				fieldBuf = append(fieldBuf, '"')
				prependComma = true
			}
			fieldBuf = append(fieldBuf, *key...)
			fieldBuf = append(fieldBuf, `":`...)
			fieldBuf = append(fieldBuf, valueBuf...)
			// Append this field's buffer to the overall buffer.
			buf = append(buf, fieldBuf...)
			return nil
		}
		switch field := fieldValue.Interface().(type) {
		case updateMarshaller:
			// Only marshal changes (not no-ops).
			if field.IsChange() {
				if err := appendField(field.interfaceValue()); err != nil {
					return nil, err
				}
			}
		default:
			if err := appendField(field); err != nil {
				return nil, err
			}
		}
	}
	buf = append(buf, '}')
	return buf, nil
}

type updateMarshaller interface {
	// IsChange utilizes the IsChange methods on Update and SliceUpdate to
	// detect whether the update should be marshalled to JSON.
	IsChange() bool
	// interfaceValue returns the (possibly nil) updated value as an interface{}
	// to be marshalled to JSON.
	interfaceValue() interface{}
}

// getKeyName tries to extract the marshalled key name from a struct field and
// returns nil if the field should be skipped.
func getKeyName(field reflect.StructField, fieldValue reflect.Value) *string {
	// Skip anonymous fields.
	if field.Anonymous {
		return nil
	}
	// Skip unexported fields (which have a non-empty PkgPath).
	if field.PkgPath != "" {
		return nil
	}
	jsonTag := field.Tag.Get("json")
	switch jsonTag {
	case "":
		// No JSON tag; use the field name.
		return &field.Name
	case "-":
		// Skip fields marked as always omitted.
		return nil
	default:
		opts := strings.Split(jsonTag, ",")
		// Skip empty fields with the omitempty option.
		for j := 1; j < len(opts); j++ {
			if opts[j] == "omitempty" {
				if isEmptyValue(fieldValue) {
					return nil
				}
			}
		}
		// The first option is the key name.
		if opts[0] == "" {
			// Key name option is empty; use the field name instead.
			return &field.Name
		}
		return &opts[0]
	}
}

// isEmptyValue checks if a reflected value is empty, as defined by
// https://golang.org/pkg/encoding/json/#Marshal. The implementation is from
// encoding/json/encode.go. This function falls under the following license:
//
// Copyright (c) 2009 The Go Authors. All rights reserved.
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//
//   - Redistributions of source code must retain the above copyright
//
// notice, this list of conditions and the following disclaimer.
//   - Redistributions in binary form must reproduce the above
//
// copyright notice, this list of conditions and the following disclaimer
// in the documentation and/or other materials provided with the
// distribution.
//   - Neither the name of Google Inc. nor the names of its
//
// contributors may be used to endorse or promote products derived from
// this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return false
}
