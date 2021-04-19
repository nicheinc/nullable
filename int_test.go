package nullable

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
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
			name:     "EmptyJSONObject",
			json:     `{}`,
			expected: Int{},
		},
		{
			name:     "NullInt",
			json:     `{"int": null}`,
			expected: MakeIntPtr(nil),
		},
		{
			name:     "ValueInt",
			json:     fmt.Sprintf(`{"int": %v}`, testInt),
			expected: MakeInt(testInt),
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
			name:     "NullInt",
			value:    nil,
			expected: MakeIntPtr(nil),
		},
		{
			name:     "ValueInt",
			value:    func(v int) *int { return &v }(testInt),
			expected: MakeInt(testInt),
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
			name:     "NotSet",
			i:        Int{},
			expected: false,
		},
		{
			name:     "NullInt",
			i:        MakeIntPtr(nil),
			expected: true,
		},
		{
			name:     "ValueInt",
			i:        MakeInt(testInt),
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
			name:     "NotSet",
			i:        Int{},
			expected: false,
		},
		{
			name:     "NullInt",
			i:        MakeIntPtr(nil),
			expected: false,
		},
		{
			name:     "ZeroInt",
			i:        MakeInt(0),
			expected: true,
		},
		{
			name:     "NonZeroInt",
			i:        MakeInt(1),
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
			name:     "NotSet",
			i:        Int{},
			expected: false,
		},
		{
			name:     "NullInt",
			i:        MakeIntPtr(nil),
			expected: false,
		},
		{
			name:     "NegativeInt",
			i:        MakeInt(-1),
			expected: true,
		},
		{
			name:     "PositiveInt",
			i:        MakeInt(1),
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

func TestInt_Value(t *testing.T) {
	i := Int{}
	if i.Value() != nil {
		t.Errorf("Expected: nil, Actual: %v", i.Value())
	}
	expected := 1
	i.SetValue(expected)
	if *i.Value() != expected {
		t.Errorf("Expected: %v, Actual: %v", expected, *i.Value())
	}
}
