package nullable

import (
	"encoding/json"
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
			name:     "Primitive",
			input:    5,
			expected: "5",
		},
		{
			name: "StructValue",
			input: struct {
				Value int
			}{
				Value: 5,
			},
			expected: `{"Value":5}`,
		},
		{
			name: "Unexported",
			input: &struct {
				unexported int
			}{
				unexported: 0,
			},
			expected: "{}",
		},
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
			name: "Omitempy/Nonempty",
			input: &struct {
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
			name: "Update/Unset",
			input: &struct {
				Field Update[int]
			}{
				Field: Update[int]{},
			},
			expected: `{}`,
		},
		{
			name: "Update/Removed",
			input: &struct {
				Field Update[int]
			}{
				Field: NewUpdatePtr[int](nil),
			},
			expected: `{"Field":null}`,
		},
		{
			name: "Update/Set",
			input: &struct {
				Field Update[int]
			}{
				Field: NewUpdate(1),
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
				t.Errorf("Expected: %v. Actual: %v", testCase.expected, actual)
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
	if err := json.Unmarshal(data, output); err != nil {
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
			name:  "Primitive",
			input: new(int),
		},
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
			name: "Update/Unset",
			input: &struct {
				Field Update[int]
			}{
				Field: Update[int]{},
			},
		},
		{
			name: "Update/Removed",
			input: &struct {
				Field Update[int]
			}{
				Field: NewUpdatePtr[int](nil),
			},
		},
		{
			name: "Update/Set",
			input: &struct {
				Field Update[int]
			}{
				Field: NewUpdate(1),
			},
		},
		{
			name: "SliceUpdate/Unset",
			input: &struct {
				Field SliceUpdate[int]
			}{
				Field: SliceUpdate[int]{},
			},
		},
		{
			name: "SliceUpdate/Removed",
			input: &struct {
				Field SliceUpdate[int]
			}{
				Field: NewSliceUpdate[int](nil),
			},
		},
		{
			name: "SliceUpdate/Set",
			input: &struct {
				Field SliceUpdate[int]
			}{
				Field: NewSliceUpdate([]int{}),
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
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			output := roundTrip(t, testCase.input)
			// Marshalling then unmarshalling should result in the same value.
			if !reflect.DeepEqual(testCase.input, output) {
				t.Errorf("Expected: %v. Actual: %v", testCase.input, output)
			}
		})
	}
}
