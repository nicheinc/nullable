package nup

import (
	"encoding/json"
	"errors"
	"reflect"
	"testing"
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
			name: "SliceUpdate/Set/Nil",
			input: struct {
				Field SliceUpdate[int]
			}{
				Field: SliceSet([]int(nil)),
			},
			expected: `{"Field":null}`,
		},
		{
			name: "SliceUpdate/Set/Empty",
			input: struct {
				Field SliceUpdate[int]
			}{
				Field: SliceSet([]int{}),
			},
			expected: `{"Field":[]}`,
		},
		{
			name: "SliceUpdate/Set/Nonempty",
			input: struct {
				Field SliceUpdate[int]
			}{
				Field: SliceSet([]int{1}),
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
			if err != nil {
				t.Errorf("Error while marshalling: %v", err)
			}
			actual := string(data)
			if actual != testCase.expected {
				t.Errorf("Expected: %v. Actual: %v", testCase.expected, actual)
			}
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
	// NOTE: SliceSet[T](nil) does not survive the roundtrip; it will be
	// unmarshalled as SliceRemove[T](). That's because SliceUpdate[T]'s
	// UnmarshalJSON always treats null as "remove", not as "set to the nil
	// slice", and that's okay because removing a slice field is semantically
	// equivalent to setting it to nil.
	roundtrip(t, "SliceUpdate/Set/Empty", struct {
		Field SliceUpdate[int]
	}{
		Field: SliceSet([]int{}),
	})
	roundtrip(t, "SliceUpdate/Set/Nonempty", struct {
		Field SliceUpdate[int]
	}{
		Field: SliceSet([]int{1}),
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
		data, err := MarshalJSON(input)
		if err != nil {
			t.Errorf("Error while marshalling: %v", err)
		}
		// Unmarshal resulting JSON to output.
		var output T
		if err := json.Unmarshal(data, &output); err != nil {
			t.Errorf("Error while unmarshalling: %v", err)
		}
		// Marshalling then unmarshalling should result in the same value.
		if !reflect.DeepEqual(input, output) {
			t.Errorf("Expected: %v. Actual: %v", input, output)
		}
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
				Field: SliceSet([]badField{{}}),
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			data, err := MarshalJSON(testCase.input)
			if data != nil {
				t.Errorf("Expected nil data. Actual: %v", string(data))
			}
			if err == nil {
				t.Error("Expected a non-nil error")
			}
		})
	}
}

type badField struct{}

func (f badField) MarshalJSON() ([]byte, error) {
	return nil, errors.New("error marshalling field")
}
