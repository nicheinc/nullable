package nullable

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/nicheinc/go-common/v12/test"
)

// Ensure implementation of Nullable interface.
var _ Nullable = &String{}

func TestString_UnmarshalJSON(t *testing.T) {
	testCases := []struct {
		name     string
		json     string
		expected String
	}{
		{
			name: "EmptyJSONObject",
			json: `{}`,
			expected: String{
				set:   false,
				value: nil,
			},
		},
		{
			name: "NullString",
			json: `{"string": null}`,
			expected: String{
				set:   true,
				value: nil,
			},
		},
		{
			name: "EmptyString",
			json: `{"string": ""}`,
			expected: String{
				set:   true,
				value: test.StrToPtr(""),
			},
		},
		{
			name: "SpaceString",
			json: `{"string": " "}`,
			expected: String{
				set:   true,
				value: test.StrToPtr(" "),
			},
		},
		{
			name: "ValueString",
			json: `{"string": "value"}`,
			expected: String{
				set:   true,
				value: test.StrToPtr("value"),
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			var dst struct {
				String String `json:"string"`
			}
			if err := json.Unmarshal([]byte(testCase.json), &dst); err != nil {
				t.Errorf("Error unmarshaling JSON: %s", err)
			}
			if !reflect.DeepEqual(dst.String, testCase.expected) {
				t.Errorf("Expected: %v. Actual: %v", testCase.expected, dst.String)
			}
		})
	}
}

func TestString_SetValue(t *testing.T) {
	testCases := []struct {
		name     string
		value    *string
		expected String
	}{
		{
			name:  "NullString",
			value: nil,
			expected: String{
				set:   true,
				value: nil,
			},
		},
		{
			name:  "EmptyString",
			value: test.StrToPtr(""),
			expected: String{
				set:   true,
				value: test.StrToPtr(""),
			},
		},
		{
			name:  "SpaceString",
			value: test.StrToPtr(" "),
			expected: String{
				set:   true,
				value: test.StrToPtr(" "),
			},
		},
		{
			name:  "ValueString",
			value: test.StrToPtr("value"),
			expected: String{
				set:   true,
				value: test.StrToPtr("value"),
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actual := String{}
			if testCase.value != nil {
				actual.SetValue(*testCase.value)
			} else {
				actual.SetPtr(nil)
			}
			if !reflect.DeepEqual(actual, testCase.expected) {
				t.Errorf("Expected: %v, Actual: %v", testCase.expected, actual)
			}
		})
	}
}

func TestString_Removed(t *testing.T) {
	testCases := []struct {
		name     string
		str      String
		expected bool
	}{
		{
			name: "NotSet",
			str: String{
				set:   false,
				value: nil,
			},
			expected: false,
		},
		{
			name: "NullString",
			str: String{
				set:   true,
				value: nil,
			},
			expected: true,
		},
		{
			name: "EmptyString",
			str: String{
				set:   true,
				value: test.StrToPtr(""),
			},
			expected: false,
		},
		{
			name: "SpaceString",
			str: String{
				set:   true,
				value: test.StrToPtr(" "),
			},
			expected: false,
		},
		{
			name: "ValueString",
			str: String{
				set:   true,
				value: test.StrToPtr("value"),
			},
			expected: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actual := testCase.str.Removed()
			if !reflect.DeepEqual(actual, testCase.expected) {
				t.Errorf("Expected: %v, Actual: %v", testCase.expected, actual)
			}
		})
	}
}

func TestString_IsEmpty(t *testing.T) {
	testCases := []struct {
		name     string
		str      String
		expected bool
	}{
		{
			name: "NotSet",
			str: String{
				set:   false,
				value: nil,
			},
			expected: false,
		},
		{
			name: "NullString",
			str: String{
				set:   true,
				value: nil,
			},
			expected: false,
		},
		{
			name: "EmptyString",
			str: String{
				set:   true,
				value: test.StrToPtr(""),
			},
			expected: true,
		},
		{
			name: "SpaceString",
			str: String{
				set:   true,
				value: test.StrToPtr(" "),
			},
			expected: true,
		},
		{
			name: "TabString",
			str: String{
				set:   true,
				value: test.StrToPtr("\t"),
			},
			expected: true,
		},
		{
			name: "NewlineString",
			str: String{
				set:   true,
				value: test.StrToPtr("\n"),
			},
			expected: true,
		},
		{
			name: "ValueString",
			str: String{
				set:   true,
				value: test.StrToPtr("value"),
			},
			expected: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actual := testCase.str.IsEmpty()
			if !reflect.DeepEqual(actual, testCase.expected) {
				t.Errorf("Expected: %v, Actual: %v", testCase.expected, actual)
			}
		})
	}
}
