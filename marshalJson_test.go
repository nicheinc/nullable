package nullable

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/nicheinc/go-common/v12/test"
)

type testStruct struct {
	Omitted          int `json:"-"`
	NoTag            int
	EmptyName        int         `json:","`
	NonNullableI     int         `json:"nonNullableI,omitempty"`
	NullableStr      String      `json:"nullableStr,omitempty"`
	NullableI        Int         `json:"nullableI,omitempty"`
	NullableB        Bool        `json:"nullableB,omitempty"`
	NullableStrSlice StringSlice `json:"nullableStrSlice,omitempty"`
}

// roundTrip returns unmarshal(marshal(v)).
func roundTrip(t *testing.T, input, output *testStruct) {
	// Marshal input.
	data, err := MarshalJSON(input)
	if err != nil {
		t.Errorf("Error while marshalling: %v", err)
	}
	// Unmarshal resulting JSON to output.
	if err := json.Unmarshal(data, &output); err != nil {
		t.Errorf("Error while unmarshalling: %v", err)
	}
}

func TestMarshalJSON_RoundTrip(t *testing.T) {
	testCases := []struct {
		name  string
		input testStruct
	}{
		{
			name: "EmptyJSONObject",
			input: testStruct{
				Omitted:      0,
				NoTag:        1,
				EmptyName:    2,
				NonNullableI: 3,
				NullableStr: String{
					set: true,
				},
				NullableI: Int{
					set:   true,
					value: test.IntToPtr(4),
				},
				NullableB: Bool{
					set:   true,
					value: test.BoolToPtr(true),
				},
				NullableStrSlice: StringSlice{
					set:   true,
					value: []string{"Hello, world!"},
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			var output testStruct
			roundTrip(t, &testCase.input, &output)
			// Marshalling then unmarshalling should result in the same value.
			if !reflect.DeepEqual(testCase.input, output) {
				t.Errorf("Expected: %v, Actual: %v", testCase.input, output)
			}
		})
	}
}
