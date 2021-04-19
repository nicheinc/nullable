package nullable

import (
	"encoding/json"
	"reflect"
	"testing"
)

type kitchenSink struct {
	unexported       int
	Omitted          int `json:"-"`
	NoTag            int
	EmptyName        int         `json:","`
	NonNullableI     int         `json:"nonNullableI,omitempty"`
	NullableStr      String      `json:"nullableStr,omitempty"`
	NullableI        Int         `json:"nullableI,omitempty"`
	NullableB        Bool        `json:"nullableB,omitempty"`
	NullableStrSlice StringSlice `json:"nullableStrSlice,omitempty"`
}

func TestMarshalJSON_OneWay(t *testing.T) {
	testCases := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{
			name: "Omitted",
			input: &struct {
				Omitted int `json:"-"`
			}{
				Omitted: 0,
			},
			expected: `{}`,
		},
		{
			name: "Omitempy/Empty",
			input: &struct {
				Slice []int   `json:",omitempty"`
				B     bool    `json:",omitempty"`
				I     int     `json:",omitempty"`
				U     uint    `json:",omitempty"`
				F     float32 `json:",omitempty"`
				Ptr   *int    `json:",omitempty"`
			}{
				Slice: []int{},
				B:     false,
				I:     0,
				U:     0,
				F:     0,
				Ptr:   nil,
			},
			expected: `{}`,
		},
		{
			name: "NoTag",
			input: &struct {
				NoTag bool
			}{
				NoTag: true,
			},
			expected: `{"NoTag":true}`,
		},
		{
			name: "EmptyName",
			input: &struct {
				EmptyName int `json:","`
			}{
				EmptyName: 1,
			},
			expected: `{"EmptyName":1}`,
		},
		{
			name: "NullableInt/Unset",
			input: &struct {
				Field Int
			}{
				Field: Int{},
			},
			expected: `{}`,
		},
		{
			name: "NullableInt/Removed",
			input: &struct {
				Field Int
			}{
				Field: NewIntPtr(nil),
			},
			expected: `{"Field":null}`,
		},
		{
			name: "NullableInt/Set",
			input: &struct {
				Field Int
			}{
				Field: NewInt(1),
			},
			expected: `{"Field":1}`,
		},
		{
			name: "MultipleFields",
			input: &struct {
				First  int
				Second int
				Third  int
			}{
				First:  1,
				Second: 2,
				Third:  3,
			},
			expected: `{"First":1,"Second":2,"Third":3}`,
		},
		{
			name: "AnonymousField",
			input: &struct {
				int
			}{
				int: 1,
			},
			expected: `{}`,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			data, err := MarshalJSON(testCase.input)
			if err != nil {
				t.Errorf("Error while marshalling: %v", err)
			}
			actual := string(data)
			if actual != testCase.expected {
				t.Errorf("Expected: %v, Actual: %v", testCase.expected, actual)
			}
		})
	}
}

// roundTrip returns unmarshal(marshal(input)).
func roundTrip(t *testing.T, input interface{}) interface{} {
	// Declare output with same type as input.
	output := reflect.New(reflect.TypeOf(input).Elem()).Interface()
	// Marshal input.
	data, err := MarshalJSON(input)
	if err != nil {
		t.Errorf("Error while marshalling: %v", err)
	}
	// Unmarshal resulting JSON to output.
	if err := json.Unmarshal(data, &output); err != nil {
		t.Errorf("Error while unmarshalling: %v", err)
	}
	return output
}

func TestMarshalJSON_RoundTrip(t *testing.T) {
	testCases := []struct {
		name  string
		input interface{}
	}{
		{
			name:  "EmptyStruct",
			input: &struct{}{},
		},
		{
			name: "Omitted",
			input: &struct {
				Omitted int `json:"-"`
			}{
				Omitted: 0,
			},
		},
		{
			name: "Omitempy/Empty",
			input: &struct {
				Omitempty int `json:"omitempty"`
			}{
				Omitempty: 0,
			},
		},
		{
			name: "Omitempy/Nonempty",
			input: &struct {
				Omitempty int `json:"omitempty"`
			}{
				Omitempty: 1,
			},
		},
		{
			name: "NoTag",
			input: &struct {
				NoTag int
			}{
				NoTag: 1,
			},
		},
		{
			name: "EmptyName",
			input: &struct {
				EmptyName int `json:","`
			}{
				EmptyName: 1,
			},
		},
		{
			name: "NullableInt/Unset",
			input: &struct {
				Field Int
			}{
				Field: Int{},
			},
		},
		{
			name: "NullableInt/Removed",
			input: &struct {
				Field Int
			}{
				Field: NewIntPtr(nil),
			},
		},
		{
			name: "NullableInt/Set",
			input: &struct {
				Field Int
			}{
				Field: NewInt(1),
			},
		},
		{
			name: "NullableBool/Unset",
			input: &struct {
				Field Bool
			}{
				Field: Bool{},
			},
		},
		{
			name: "NullableBool/Removed",
			input: &struct {
				Field Bool
			}{
				Field: NewBoolPtr(nil),
			},
		},
		{
			name: "NullableBool/Set",
			input: &struct {
				Field Bool
			}{
				Field: NewBool(true),
			},
		},
		{
			name: "NullableString/Unset",
			input: &struct {
				Field String
			}{
				Field: String{},
			},
		},
		{
			name: "NullableString/Removed",
			input: &struct {
				Field String
			}{
				Field: NewStringPtr(nil),
			},
		},
		{
			name: "NullableString/Set",
			input: &struct {
				Field String
			}{
				Field: NewString(""),
			},
		},
		{
			name: "NullableStringSlice/Unset",
			input: &struct {
				Field StringSlice
			}{
				Field: StringSlice{},
			},
		},
		{
			name: "NullableStringSlice/Removed",
			input: &struct {
				Field StringSlice
			}{
				Field: NewStringSlice(nil),
			},
		},
		{
			name: "NullableStringSlice/Set",
			input: &struct {
				Field StringSlice
			}{
				Field: NewStringSlice([]string{}),
			},
		},
		{
			name: "MultipleFields",
			input: &struct {
				First  int
				Second int
				Third  int
			}{
				First:  1,
				Second: 2,
				Third:  3,
			},
		},
		{
			name: "KitchenSink",
			input: &kitchenSink{
				unexported:       0,
				Omitted:          0,
				NoTag:            1,
				EmptyName:        2,
				NonNullableI:     3,
				NullableStr:      NewStringPtr(nil),
				NullableI:        NewInt(4),
				NullableB:        NewBool(true),
				NullableStrSlice: NewStringSlice([]string{"Hello, world!"}),
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			output := roundTrip(t, testCase.input)
			// Marshalling then unmarshalling should result in the same value.
			if !reflect.DeepEqual(testCase.input, output) {
				t.Errorf("Expected: %v, Actual: %v", testCase.input, output)
			}
		})
	}
}
