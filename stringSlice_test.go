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
	expected := []string{testString}
	s.SetValue(expected)
	if !reflect.DeepEqual(s.Value(), expected) {
		t.Errorf("Expected: %v, Actual: %v", expected, s.Value())
	}
}

func TestStringSlice_Equals(t *testing.T) {
	testCases := []struct {
		name     string
		s        StringSlice
		value    []string
		expected bool
	}{
		{
			name:     "Unset",
			s:        StringSlice{},
			value:    nil,
			expected: false,
		},
		{
			name:     "Removed",
			s:        NewStringSlice(nil),
			value:    nil,
			expected: false,
		},
		{
			name:     "Set/NotEqual/DifferentSizes",
			s:        NewStringSlice([]string{testString, testString}),
			value:    []string{testString},
			expected: false,
		},
		{
			name:     "Set/NotEqual/DifferentElements",
			s:        NewStringSlice([]string{"value1"}),
			value:    []string{"value2"},
			expected: false,
		},
		{
			name:     "Set/Equal/Empty/Nil",
			s:        NewStringSlice([]string{}),
			value:    nil,
			expected: true,
		},
		{
			name:     "Set/Equal/Empty/Empty",
			s:        NewStringSlice([]string{}),
			value:    []string{},
			expected: true,
		},
		{
			name:     "Set/Equal/Nonempty/Nonempty",
			s:        NewStringSlice([]string{testString}),
			value:    []string{testString},
			expected: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if actual := testCase.s.Equals(testCase.value); actual != testCase.expected {
				t.Errorf("Expected: %v. Actual: %v", testCase.expected, actual)
			}
		})
	}
}

func TestStringSlice_Apply(t *testing.T) {
	value := []string{testString}
	testCases := []struct {
		name     string
		s        StringSlice
		expected []string
	}{
		{
			name:     "Unset",
			s:        StringSlice{},
			expected: value,
		},
		{
			name:     "Removed",
			s:        NewStringSlice(nil),
			expected: nil,
		},
		{
			name:     "Set",
			s:        NewStringSlice([]string{testString, testString}),
			expected: []string{testString, testString},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if actual := testCase.s.Apply(value); !reflect.DeepEqual(actual, testCase.expected) {
				t.Errorf("Expected: %v. Actual: %v", testCase.expected, actual)
			}
		})
	}
}

func TestStringSlice_Diff(t *testing.T) {
	var (
		value1 = []string{"value1"}
		value2 = []string{"value2"}
	)
	testCases := []struct {
		name     string
		s        StringSlice
		value    []string
		expected StringSlice
	}{
		{
			name:     "Unset",
			s:        StringSlice{},
			value:    value1,
			expected: StringSlice{},
		},
		{
			name:     "Removed/NonZeroValue",
			s:        NewStringSlice(nil),
			value:    value1,
			expected: NewStringSlice(nil),
		},
		{
			name:     "Removed/ZeroValue",
			s:        NewStringSlice(nil),
			value:    nil,
			expected: StringSlice{},
		},
		{
			name:     "Set/Equal",
			s:        NewStringSlice(value1),
			value:    value1,
			expected: StringSlice{},
		},
		{
			name:     "Set/NotEqual",
			s:        NewStringSlice(value2),
			value:    value1,
			expected: NewStringSlice(value2),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if actual := testCase.s.Diff(testCase.value); !reflect.DeepEqual(actual, testCase.expected) {
				t.Errorf("Expected: %v. Actual: %v", testCase.expected, actual)
			}
		})
	}
}

func TestStringSlice_InterfaceValue(t *testing.T) {
	var s StringSlice
	if !reflect.ValueOf(s.InterfaceValue()).IsNil() {
		t.Errorf("Expected: nil, Actual: %v", s.InterfaceValue())
	}
	expected := []string{testString}
	s.SetValue(expected)
	if !reflect.DeepEqual(s.InterfaceValue(), expected) {
		t.Errorf("Expected: %v, Actual: %v", expected, s.InterfaceValue())
	}
}