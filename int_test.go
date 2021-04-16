package nullable

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/nicheinc/go-common/v12/test"
)

const testInt = 42

// Ensure implementation of Nullable interface.
var _ Nullable = &Int{}

func TestInt_UnmarshalJSON(t *testing.T) {
	testCases := []struct {
		name     string
		json     string
		expected Int
	}{
		{
			name: "EmptyJSONObject",
			json: `{}`,
			expected: Int{
				set:   false,
				value: nil,
			},
		},
		{
			name: "NullInt",
			json: `{"int": null}`,
			expected: Int{
				set:   true,
				value: nil,
			},
		},
		{
			name: "ValueInt",
			json: fmt.Sprintf(`{"int": %v}`, testInt),
			expected: Int{
				set:   true,
				value: test.IntToPtr(testInt),
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			var dst struct {
				Int Int `json:"int"`
			}
			if err := json.Unmarshal([]byte(testCase.json), &dst); err != nil {
				t.Errorf("Error unmarshaling JSON: %s", err)
			}
			if !reflect.DeepEqual(dst.Int, testCase.expected) {
				t.Errorf("Expected: %v. Actual: %v", testCase.expected, dst.Int)
			}
		})
	}
}

func TestInt_SetValue(t *testing.T) {
	testCases := []struct {
		name     string
		value    *int
		expected Int
	}{
		{
			name:  "NullInt",
			value: nil,
			expected: Int{
				set:   true,
				value: nil,
			},
		},
		{
			name:  "ValueInt",
			value: test.IntToPtr(testInt),
			expected: Int{
				set:   true,
				value: test.IntToPtr(testInt),
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actual := Int{}
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

func TestInt_Removed(t *testing.T) {
	testCases := []struct {
		name     string
		i        Int
		expected bool
	}{
		{
			name: "NotSet",
			i: Int{
				set:   false,
				value: nil,
			},
			expected: false,
		},
		{
			name: "NullInt",
			i: Int{
				set:   true,
				value: nil,
			},
			expected: true,
		},
		{
			name: "ValueInt",
			i: Int{
				set:   true,
				value: test.IntToPtr(testInt),
			},
			expected: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actual := testCase.i.Removed()
			if !reflect.DeepEqual(actual, testCase.expected) {
				t.Errorf("Expected: %v, Actual: %v", testCase.expected, actual)
			}
		})
	}
}

func TestInt_IsZero(t *testing.T) {
	testCases := []struct {
		name     string
		i        Int
		expected bool
	}{
		{
			name: "NotSet",
			i: Int{
				set:   false,
				value: nil,
			},
			expected: false,
		},
		{
			name: "NullInt",
			i: Int{
				set:   true,
				value: nil,
			},
			expected: false,
		},
		{
			name: "ZeroInt",
			i: Int{
				set:   true,
				value: test.IntToPtr(0),
			},
			expected: true,
		},
		{
			name: "NonZeroInt",
			i: Int{
				set:   true,
				value: test.IntToPtr(1),
			},
			expected: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actual := testCase.i.IsZero()
			if !reflect.DeepEqual(actual, testCase.expected) {
				t.Errorf("Expected: %v, Actual: %v", testCase.expected, actual)
			}
		})
	}
}

func TestInt_IsNegative(t *testing.T) {
	testCases := []struct {
		name     string
		i        Int
		expected bool
	}{
		{
			name: "NotSet",
			i: Int{
				set:   false,
				value: nil,
			},
			expected: false,
		},
		{
			name: "NullInt",
			i: Int{
				set:   true,
				value: nil,
			},
			expected: false,
		},
		{
			name: "NegativeInt",
			i: Int{
				set:   true,
				value: test.IntToPtr(-1),
			},
			expected: true,
		},
		{
			name: "PositiveInt",
			i: Int{
				set:   true,
				value: test.IntToPtr(1),
			},
			expected: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actual := testCase.i.IsNegative()
			if !reflect.DeepEqual(actual, testCase.expected) {
				t.Errorf("Expected: %v, Actual: %v", testCase.expected, actual)
			}
		})
	}
}
