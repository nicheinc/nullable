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
			expected: NewIntPtr(nil),
		},
		{
			name:     "ValueInt",
			json:     fmt.Sprintf(`{"int": %v}`, testInt),
			expected: NewInt(testInt),
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
			expected: NewIntPtr(nil),
		},
		{
			name:     "ValueInt",
			value:    func(v int) *int { return &v }(testInt),
			expected: NewInt(testInt),
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
			i:        NewIntPtr(nil),
			expected: true,
		},
		{
			name:     "ValueInt",
			i:        NewInt(testInt),
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
			i:        NewIntPtr(nil),
			expected: false,
		},
		{
			name:     "ZeroInt",
			i:        NewInt(0),
			expected: true,
		},
		{
			name:     "NonZeroInt",
			i:        NewInt(1),
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
			i:        NewIntPtr(nil),
			expected: false,
		},
		{
			name:     "NegativeInt",
			i:        NewInt(-1),
			expected: true,
		},
		{
			name:     "PositiveInt",
			i:        NewInt(1),
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
	var i Int
	if i.Value() != nil {
		t.Errorf("Expected: nil, Actual: %v", i.Value())
	}
	expected := 1
	i.SetValue(expected)
	if *i.Value() != expected {
		t.Errorf("Expected: %v, Actual: %v", expected, *i.Value())
	}
}

func TestInt_Equals(t *testing.T) {
	value := 1
	testCases := []struct {
		name     string
		i        Int
		expected bool
	}{
		{
			name:     "Unset",
			i:        Int{},
			expected: false,
		},
		{
			name:     "Removed",
			i:        NewIntPtr(nil),
			expected: false,
		},
		{
			name:     "Set/NotEqualValue",
			i:        NewInt(value + 1),
			expected: false,
		},
		{
			name:     "Set/Equal",
			i:        NewInt(value),
			expected: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if actual := testCase.i.Equals(value); actual != testCase.expected {
				t.Errorf("Expected: %v. Actual: %v", testCase.expected, actual)
			}
		})
	}
}

func TestInt_Apply(t *testing.T) {
	value := 1
	testCases := []struct {
		name     string
		i        Int
		expected int
	}{
		{
			name:     "Unset",
			i:        Int{},
			expected: value,
		},
		{
			name:     "Removed",
			i:        NewIntPtr(nil),
			expected: 0,
		},
		{
			name:     "Set",
			i:        NewInt(value + 1),
			expected: value + 1,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if actual := testCase.i.Apply(value); actual != testCase.expected {
				t.Errorf("Expected: %v. Actual: %v", testCase.expected, actual)
			}
		})
	}
}

func TestInt_ApplyPtr(t *testing.T) {
	var (
		value1 = 1
		value2 = value1 + 1
	)
	testCases := []struct {
		name     string
		i        Int
		expected *int
	}{
		{
			name:     "Unset",
			i:        Int{},
			expected: &value1,
		},
		{
			name:     "Removed",
			i:        NewIntPtr(nil),
			expected: nil,
		},
		{
			name:     "Set",
			i:        NewInt(value2),
			expected: &value2,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if actual := testCase.i.ApplyPtr(&value1); !reflect.DeepEqual(actual, testCase.expected) {
				t.Errorf("Expected: %v. Actual: %v", testCase.expected, actual)
			}
		})
	}
}

func TestInt_Diff(t *testing.T) {
	var (
		value1 = 1
		value2 = 2
	)
	testCases := []struct {
		name     string
		i        Int
		value    int
		expected Int
	}{
		{
			name:     "Unset",
			i:        Int{},
			value:    value1,
			expected: Int{},
		},
		{
			name:     "Removed/NonZeroValue",
			i:        NewIntPtr(nil),
			value:    value1,
			expected: NewIntPtr(nil),
		},
		{
			name:     "Removed/ZeroValue",
			i:        NewIntPtr(nil),
			value:    0.0,
			expected: Int{},
		},
		{
			name:     "Set/Equal",
			i:        NewInt(value1),
			value:    value1,
			expected: Int{},
		},
		{
			name:     "Set/NotEqual",
			i:        NewInt(value2),
			value:    value1,
			expected: NewInt(value2),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if actual := testCase.i.Diff(testCase.value); !reflect.DeepEqual(actual, testCase.expected) {
				t.Errorf("Expected: %v. Actual: %v", testCase.expected, actual)
			}
		})
	}
}

func TestInt_InterfaceValue(t *testing.T) {
	var i Int
	if !reflect.ValueOf(i.InterfaceValue()).IsNil() {
		t.Errorf("Expected: nil, Actual: %v", i.InterfaceValue())
	}
	expected := 1
	i.SetValue(expected)
	if !reflect.DeepEqual(i.InterfaceValue(), &expected) {
		t.Errorf("Expected: %v, Actual: %v", expected, i.InterfaceValue())
	}
}
