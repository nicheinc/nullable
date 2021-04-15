package nullable

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

const testString = "e^(i*tau) = 1"

func TestStringSlice_UnmarshalJSON(t *testing.T) {
	testCases := []struct {
		name     string
		json     string
		expected StringSlice
	}{
		{
			name: "EmptyJSONObject",
			json: `{}`,
			expected: StringSlice{
				Set:   false,
				Value: nil,
			},
		},
		{
			name: "NullStringSlice",
			json: `{"stringSlice": null}`,
			expected: StringSlice{
				Set:   true,
				Value: nil,
			},
		},
		{
			name: "EmptyStringSlice",
			json: `{"stringSlice": []}`,
			expected: StringSlice{
				Set:   true,
				Value: []string{},
			},
		},
		{
			name: "NonEmptyStringSlice",
			json: fmt.Sprintf(`{"stringSlice": ["%s"]}`, testString),
			expected: StringSlice{
				Set:   true,
				Value: []string{testString},
			},
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
			name:  "NullStringSlice",
			value: nil,
			expected: StringSlice{
				Set:   true,
				Value: nil,
			},
		},
		{
			name:  "EmptyStringSlice",
			value: []string{},
			expected: StringSlice{
				Set:   true,
				Value: []string{},
			},
		},
		{
			name:  "NonEmptyStringSlice",
			value: []string{testString},
			expected: StringSlice{
				Set:   true,
				Value: []string{testString},
			},
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
			name: "NotSet",
			strSlice: StringSlice{
				Set:   false,
				Value: nil,
			},
			expected: false,
		},
		{
			name: "NullStringSlice",
			strSlice: StringSlice{
				Set:   true,
				Value: nil,
			},
			expected: true,
		},
		{
			name: "EmptyStringSlice",
			strSlice: StringSlice{
				Set:   true,
				Value: []string{},
			},
			expected: false,
		},
		{
			name: "NonEmptyStringSlice",
			strSlice: StringSlice{
				Set:   true,
				Value: []string{testString},
			},
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
			name: "NotSet",
			strSlice: StringSlice{
				Set:   false,
				Value: nil,
			},
			expected: false,
		},
		{
			name: "NullStringSlice",
			strSlice: StringSlice{
				Set:   true,
				Value: nil,
			},
			expected: false,
		},
		{
			name: "EmptyStringSlice",
			strSlice: StringSlice{
				Set:   true,
				Value: []string{},
			},
			expected: true,
		},
		{
			name: "NonEmptyStringSlice",
			strSlice: StringSlice{
				Set:   true,
				Value: []string{testString},
			},
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
