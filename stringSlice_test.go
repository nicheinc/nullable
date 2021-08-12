package nullable

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

const testString = "e^(i*tau) = 1"

// Ensure implementation of Nullable interface.
var _ Nullable = &StringSlice{}

func TestStringSlice_UnmarshalJSON(t *testing.T) {
	testCases := []struct {
		name     string
		json     string
		expected StringSlice
	}{
		{
			name:     "EmptyJSONObject",
			json:     `{}`,
			expected: StringSlice{},
		},
		{
			name:     "NullStringSlice",
			json:     `{"stringSlice": null}`,
			expected: NewStringSlice(nil),
		},
		{
			name:     "EmptyStringSlice",
			json:     `{"stringSlice": []}`,
			expected: NewStringSlice([]string{}),
		},
		{
			name:     "NonEmptyStringSlice",
			json:     fmt.Sprintf(`{"stringSlice": ["%s"]}`, testString),
			expected: NewStringSlice([]string{testString}),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			var dst struct {
				StringSlice StringSlice `json:"stringSlice"`
			}
			if err := json.Unmarshal([]byte(testCase.json), &dst); err != nil {
				t.Errorf("Error unmarshaling JSON: %s", err)
			}
			if !reflect.DeepEqual(dst.StringSlice, testCase.expected) {
				t.Errorf("Expected: %v. Actual: %v", testCase.expected, dst.StringSlice)
			}
		})
	}
}

func TestStringSlice_SetValue(t *testing.T) {
	testCases := []struct {
		name     string
		value    []string
		expected StringSlice
	}{
		{
			name:     "NullStringSlice",
			value:    nil,
			expected: NewStringSlice(nil),
		},
		{
			name:     "EmptyStringSlice",
			value:    []string{},
			expected: NewStringSlice([]string{}),
		},
		{
			name:     "NonEmptyStringSlice",
			value:    []string{testString},
			expected: NewStringSlice([]string{testString}),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actual := StringSlice{}
			actual.SetValue(testCase.value)
			if !reflect.DeepEqual(actual, testCase.expected) {
				t.Errorf("Expected: %v, Actual: %v", testCase.expected, actual)
			}
		})
	}
}

func TestStringSlice_Removed(t *testing.T) {
	testCases := []struct {
		name     string
		strSlice StringSlice
		expected bool
	}{
		{
			name:     "NotSet",
			strSlice: StringSlice{},
			expected: false,
		},
		{
			name:     "NullStringSlice",
			strSlice: NewStringSlice(nil),
			expected: true,
		},
		{
			name:     "EmptyStringSlice",
			strSlice: NewStringSlice([]string{}),
			expected: false,
		},
		{
			name:     "NonEmptyStringSlice",
			strSlice: NewStringSlice([]string{testString}),
			expected: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actual := testCase.strSlice.Removed()
			if !reflect.DeepEqual(actual, testCase.expected) {
				t.Errorf("Expected: %v, Actual: %v", testCase.expected, actual)
			}
		})
	}
}

func TestStringSlice_IsEmpty(t *testing.T) {
	testCases := []struct {
		name     string
		strSlice StringSlice
		expected bool
	}{
		{
			name:     "NotSet",
			strSlice: StringSlice{},
			expected: false,
		},
		{
			name:     "NullStringSlice",
			strSlice: NewStringSlice(nil),
			expected: false,
		},
		{
			name:     "EmptyStringSlice",
			strSlice: NewStringSlice([]string{}),
			expected: true,
		},
		{
			name:     "NonEmptyStringSlice",
			strSlice: NewStringSlice([]string{testString}),
			expected: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actual := testCase.strSlice.IsEmpty()
			if !reflect.DeepEqual(actual, testCase.expected) {
				t.Errorf("Expected: %v, Actual: %v", testCase.expected, actual)
			}
		})
	}
}

func TestStringSlice_Value(t *testing.T) {
	s := StringSlice{}
	if s.Value() != nil {
		t.Errorf("Expected: nil, Actual: %v", s.Value())
	}
	expected := []string{"value"}
	s.SetValue(expected)
	if !reflect.DeepEqual(s.Value(), expected) {
		t.Errorf("Expected: %v, Actual: %v", expected, s.Value())
	}
}

func TestStringSlice_InterfaceValue(t *testing.T) {
	var s StringSlice
	if !reflect.ValueOf(s.InterfaceValue()).IsNil() {
		t.Errorf("Expected: nil, Actual: %v", s.InterfaceValue())
	}
	expected := []string{"value"}
	s.SetValue(expected)
	if !reflect.DeepEqual(s.InterfaceValue(), expected) {
		t.Errorf("Expected: %v, Actual: %v", expected, s.InterfaceValue())
	}
}
