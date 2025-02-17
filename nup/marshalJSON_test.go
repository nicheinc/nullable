package nup

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/nicheinc/expect"
)

func TestMarshalJSON_OneWay(t *testing.T) {
	testCases := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{
			name:     "Nil",
			input:    nil,
			expected: "null",
		},
		{
			name:     "Primitive",
			input:    5,
			expected: "5",
		},
		{
			name: "PointerToStruct",
			input: &struct {
				Value  int
				Update Update[int]
			}{
				Value:  5,
				Update: Noop[int](),
			},
			expected: `{"Value":5}`,
		},
		{
			name: "Unexported",
			input: struct {
				unexported int
			}{
				unexported: 0,
			},
			expected: "{}",
		},
		{
			name: "Omitted",
			input: struct {
				Omitted int `json:"-"`
			}{
				Omitted: 0,
			},
			expected: `{}`,
		},
		{
			name: "Omitempy/Empty",
			input: struct {
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
			name: "Omitempy/Nonempty",
			input: struct {
				Slice      []int    `json:",omitempty"`
				B          bool     `json:",omitempty"`
				I          int      `json:",omitempty"`
				U          uint     `json:",omitempty"`
				F          float32  `json:",omitempty"`
				Ptr        *int     `json:",omitempty"`
				NeverEmpty struct{} `json:",omitempty"`
			}{
				Slice:      []int{1},
				B:          true,
				I:          1,
				U:          1,
				F:          1,
				Ptr:        func(v int) *int { return &v }(1),
				NeverEmpty: struct{}{},
			},
			expected: `{"Slice":[1],"B":true,"I":1,"U":1,"F":1,"Ptr":1,"NeverEmpty":{}}`,
		},
		{
			name: "NoTag",
			input: struct {
				NoTag bool
			}{
				NoTag: true,
			},
			expected: `{"NoTag":true}`,
		},
		{
			name: "EmptyName",
			input: struct {
				EmptyName int `json:","`
			}{
				EmptyName: 1,
			},
			expected: `{"EmptyName":1}`,
		},
		{
			name: "Update/Noop",
			input: struct {
				Field Update[int]
			}{
				Field: Noop[int](),
			},
			expected: `{}`,
		},
		{
			name: "Update/Remove",
			input: struct {
				Field Update[int]
			}{
				Field: Remove[int](),
			},
			expected: `{"Field":null}`,
		},
		{
			name: "Update/Set",
			input: struct {
				Field Update[int]
			}{
				Field: Set(1),
			},
			expected: `{"Field":1}`,
		},
		{
			name: "SliceUpdate/Noop",
			input: struct {
				Field SliceUpdate[int]
			}{
				Field: SliceNoop[int](),
			},
			expected: `{}`,
		},
		{
			name: "SliceUpdate/Remove",
			input: struct {
				Field SliceUpdate[int]
			}{
				Field: SliceRemove[int](),
			},
			expected: `{"Field":null}`,
		},
		{
			name: "SliceUpdate/RemoveOrSet/Nil",
			input: struct {
				Field SliceUpdate[int]
			}{
				Field: SliceRemoveOrSet([]int(nil)),
			},
			expected: `{"Field":null}`,
		},
		{
			name: "SliceUpdate/RemoveOrSet/Empty",
			input: struct {
				Field SliceUpdate[int]
			}{
				Field: SliceRemoveOrSet([]int{}),
			},
			expected: `{"Field":[]}`,
		},
		{
			name: "SliceUpdate/RemoveOrSet/Nonempty",
			input: struct {
				Field SliceUpdate[int]
			}{
				Field: SliceRemoveOrSet([]int{1}),
			},
			expected: `{"Field":[1]}`,
		},
		{
			name: "MultipleFields",
			input: struct {
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
			input: struct {
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
			expect.ErrorNil(t, err)
			actual := string(data)
			expect.Equal(t, actual, testCase.expected)
		})
	}
}

func TestMarshalJSON_RoundTrip(t *testing.T) {
	roundtrip(t, "Primitive", 0)
	roundtrip(t, "EmptyStruct", struct{}{})
	roundtrip(t, "Omitted", struct {
		Omitted int `json:"-"`
	}{
		Omitted: 0,
	})
	roundtrip(t, "Omitempy/Empty", struct {
		Omitempty int `json:"omitempty"`
	}{
		Omitempty: 0,
	})
	roundtrip(t, "Omitempy/Nonempty", struct {
		Omitempty int `json:"omitempty"`
	}{
		Omitempty: 1,
	})
	roundtrip(t, "NoTag", struct {
		NoTag int
	}{
		NoTag: 1,
	})
	roundtrip(t, "EmptyName", struct {
		EmptyName int `json:","`
	}{
		EmptyName: 1,
	})
	roundtrip(t, "Update/Noop", struct {
		Field Update[int]
	}{
		Field: Noop[int](),
	})
	roundtrip(t, "Update/Remove", struct {
		Field Update[int]
	}{
		Field: Remove[int](),
	})
	roundtrip(t, "Update/Set", struct {
		Field Update[int]
	}{
		Field: Set(1),
	})
	roundtrip(t, "SliceUpdate/Noop", struct {
		Field SliceUpdate[int]
	}{
		Field: SliceNoop[int](),
	})
	roundtrip(t, "SliceUpdate/Remove", struct {
		Field SliceUpdate[int]
	}{
		Field: SliceRemove[int](),
	})
	roundtrip(t, "SliceUpdate/RemoveOrSet/Nil", struct {
		Field SliceUpdate[int]
	}{
		Field: SliceRemoveOrSet[int](nil),
	})
	roundtrip(t, "SliceUpdate/RemoveOrSet/Empty", struct {
		Field SliceUpdate[int]
	}{
		Field: SliceRemoveOrSet([]int{}),
	})
	roundtrip(t, "SliceUpdate/RemoveOrSet/Nonempty", struct {
		Field SliceUpdate[int]
	}{
		Field: SliceRemoveOrSet([]int{1}),
	})
	roundtrip(t, "MultipleFields", struct {
		First  int
		Second int
		Third  int
	}{
		First:  1,
		Second: 2,
		Third:  3,
	})
}

// roundtrip checks that unmarshal(marshal(input)) == input.
func roundtrip[T any](t *testing.T, testName string, input T) {
	t.Helper()
	t.Run(testName, func(t *testing.T) {
		t.Helper()
		// Marshal input.
		data, marshalErr := MarshalJSON(input)
		expect.ErrorNil(t, marshalErr)
		// Unmarshal resulting JSON to output.
		var output T
		unmarshalErr := json.Unmarshal(data, &output)
		expect.ErrorNil(t, unmarshalErr)
		// Marshalling then unmarshalling should result in the same value.
		expect.Equal(t, input, output)
	})
}

func TestMarshalJSON_FieldErrors(t *testing.T) {
	testCases := []struct {
		name  string
		input interface{}
	}{
		{
			name: "NonUpdateField",
			input: struct {
				Field badField
			}{
				Field: badField{},
			},
		},
		{
			name: "UpdateField",
			input: struct {
				Field Update[badField]
			}{
				Field: Set(badField{}),
			},
		},
		{
			name: "SliceUpdateField",
			input: struct {
				Field SliceUpdate[badField]
			}{
				Field: SliceRemoveOrSet([]badField{{}}),
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			data, err := MarshalJSON(testCase.input)
			expect.ErrorNonNil(t, err)
			expect.Equal(t, data, nil)
		})
	}
}

type badField struct{}

func (f badField) MarshalJSON() ([]byte, error) {
	return nil, errors.New("error marshalling field")
}
