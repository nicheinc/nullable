package nullable

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/nicheinc/go-common/v12/test"
)

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
				Set:   false,
				Value: nil,
			},
		},
		{
			name: "NullString",
			json: `{"string": null}`,
			expected: String{
				Set:   true,
				Value: nil,
			},
		},
		{
			name: "EmptyString",
			json: `{"string": ""}`,
			expected: String{
				Set:   true,
				Value: test.StrToPtr(""),
			},
		},
		{
			name: "SpaceString",
			json: `{"string": " "}`,
			expected: String{
				Set:   true,
				Value: test.StrToPtr(" "),
			},
		},
		{
			name: "ValueString",
			json: `{"string": "value"}`,
			expected: String{
				Set:   true,
				Value: test.StrToPtr("value"),
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
				Set:   true,
				Value: nil,
			},
		},
		{
			name:  "EmptyString",
			value: test.StrToPtr(""),
			expected: String{
				Set:   true,
				Value: test.StrToPtr(""),
			},
		},
		{
			name:  "SpaceString",
			value: test.StrToPtr(" "),
			expected: String{
				Set:   true,
				Value: test.StrToPtr(" "),
			},
		},
		{
			name:  "ValueString",
			value: test.StrToPtr("value"),
			expected: String{
				Set:   true,
				Value: test.StrToPtr("value"),
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
				Set:   false,
				Value: nil,
			},
			expected: false,
		},
		{
			name: "NullString",
			str: String{
				Set:   true,
				Value: nil,
			},
			expected: true,
		},
		{
			name: "EmptyString",
			str: String{
				Set:   true,
				Value: test.StrToPtr(""),
			},
			expected: false,
		},
		{
			name: "SpaceString",
			str: String{
				Set:   true,
				Value: test.StrToPtr(" "),
			},
			expected: false,
		},
		{
			name: "ValueString",
			str: String{
				Set:   true,
				Value: test.StrToPtr("value"),
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
				Set:   false,
				Value: nil,
			},
			expected: false,
		},
		{
			name: "NullString",
			str: String{
				Set:   true,
				Value: nil,
			},
			expected: false,
		},
		{
			name: "EmptyString",
			str: String{
				Set:   true,
				Value: test.StrToPtr(""),
			},
			expected: true,
		},
		{
			name: "SpaceString",
			str: String{
				Set:   true,
				Value: test.StrToPtr(" "),
			},
			expected: true,
		},
		{
			name: "TabString",
			str: String{
				Set:   true,
				Value: test.StrToPtr("\t"),
			},
			expected: true,
		},
		{
			name: "NewlineString",
			str: String{
				Set:   true,
				Value: test.StrToPtr("\n"),
			},
			expected: true,
		},
		{
			name: "ValueString",
			str: String{
				Set:   true,
				Value: test.StrToPtr("value"),
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
