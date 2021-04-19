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
			name:     "EmptyJSONObject",
			json:     `{}`,
			expected: String{},
		},
		{
			name:     "NullString",
			json:     `{"string": null}`,
			expected: MakeStringPtr(nil),
		},
		{
			name:     "EmptyString",
			json:     `{"string": ""}`,
			expected: MakeString(""),
		},
		{
			name:     "SpaceString",
			json:     `{"string": " "}`,
			expected: MakeString(" "),
		},
		{
			name:     "ValueString",
			json:     `{"string": "value"}`,
			expected: MakeString("value"),
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
			name:     "NullString",
			value:    nil,
			expected: MakeStringPtr(nil),
		},
		{
			name:     "EmptyString",
			value:    test.StrToPtr(""),
			expected: MakeString(""),
		},
		{
			name:     "SpaceString",
			value:    test.StrToPtr(" "),
			expected: MakeString(" "),
		},
		{
			name:     "ValueString",
			value:    test.StrToPtr("value"),
			expected: MakeString("value"),
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
			name:     "NotSet",
			str:      String{},
			expected: false,
		},
		{
			name:     "NullString",
			str:      MakeStringPtr(nil),
			expected: true,
		},
		{
			name:     "EmptyString",
			str:      MakeString(""),
			expected: false,
		},
		{
			name:     "SpaceString",
			str:      MakeString(" "),
			expected: false,
		},
		{
			name:     "ValueString",
			str:      MakeString("value"),
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
			name:     "NotSet",
			str:      String{},
			expected: false,
		},
		{
			name:     "NullString",
			str:      MakeStringPtr(nil),
			expected: false,
		},
		{
			name:     "EmptyString",
			str:      MakeString(""),
			expected: true,
		},
		{
			name:     "SpaceString",
			str:      MakeString(" "),
			expected: true,
		},
		{
			name:     "TabString",
			str:      MakeString("\t"),
			expected: true,
		},
		{
			name:     "NewlineString",
			str:      MakeString("\n"),
			expected: true,
		},
		{
			name:     "ValueString",
			str:      MakeString("value"),
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

func TestString_Value(t *testing.T) {
	var s String
	if s.Value() != nil {
		t.Errorf("Expected: nil, Actual: %v", s.Value())
	}
	expected := "value"
	s.SetValue(expected)
	if *s.Value() != expected {
		t.Errorf("Expected: %v, Actual: %v", expected, *s.Value())
	}
}
