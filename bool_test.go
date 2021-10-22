package nullable

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

// Ensure implementation of Nullable interface.
var _ Nullable = &Bool{}

func TestBool_UnmarshalJSON(t *testing.T) {
	testCases := []struct {
		name     string
		json     string
		expected Bool
	}{
		{
			name:     "EmptyJSONObject",
			json:     `{}`,
			expected: Bool{},
		},
		{
			name:     "NullBool",
			json:     `{"int": null}`,
			expected: NewBoolPtr(nil),
		},
		{
			name:     "ValueBool",
			json:     fmt.Sprintf(`{"int": %v}`, true),
			expected: NewBool(true),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			var dst struct {
				Bool Bool `json:"int"`
			}
			if err := json.Unmarshal([]byte(testCase.json), &dst); err != nil {
				t.Errorf("Error unmarshaling JSON: %s", err)
			}
			if !reflect.DeepEqual(dst.Bool, testCase.expected) {
				t.Errorf("Expected: %v. Actual: %v", testCase.expected, dst.Bool)
			}
		})
	}
}

func TestBool_SetValue(t *testing.T) {
	testCases := []struct {
		name     string
		value    *bool
		expected Bool
	}{
		{
			name:     "NullBool",
			value:    nil,
			expected: NewBoolPtr(nil),
		},
		{
			name:     "ValueBool",
			value:    func(v bool) *bool { return &v }(true),
			expected: NewBool(true),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actual := Bool{}
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

func TestBool_Removed(t *testing.T) {
	testCases := []struct {
		name     string
		b        Bool
		expected bool
	}{
		{
			name:     "NotSet",
			b:        Bool{},
			expected: false,
		},
		{
			name:     "NullBool",
			b:        NewBoolPtr(nil),
			expected: true,
		},
		{
			name:     "ValueBool",
			b:        NewBool(true),
			expected: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actual := testCase.b.Removed()
			if !reflect.DeepEqual(actual, testCase.expected) {
				t.Errorf("Expected: %v, Actual: %v", testCase.expected, actual)
			}
		})
	}
}

func TestBool_Value(t *testing.T) {
	var b Bool
	if b.Value() != nil {
		t.Errorf("Expected: nil, Actual: %v", b.Value())
	}
	expected := true
	b.SetValue(expected)
	if *b.Value() != expected {
		t.Errorf("Expected: %v, Actual: %v", expected, *b.Value())
	}
}

func TestBool_Equals(t *testing.T) {
	value := true
	testCases := []struct {
		name     string
		b        Bool
		expected bool
	}{
		{
			name:     "Unset",
			b:        Bool{},
			expected: false,
		},
		{
			name:     "Removed",
			b:        NewBoolPtr(nil),
			expected: false,
		},
		{
			name:     "Set/NotEqual",
			b:        NewBool(!value),
			expected: false,
		},
		{
			name:     "Set/Equal",
			b:        NewBool(value),
			expected: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if actual := testCase.b.Equals(value); actual != testCase.expected {
				t.Errorf("Expected: %v. Actual: %v", testCase.expected, actual)
			}
		})
	}
}

func TestBool_Apply(t *testing.T) {
	value := true
	testCases := []struct {
		name     string
		b        Bool
		expected bool
	}{
		{
			name:     "Unset",
			b:        Bool{},
			expected: value,
		},
		{
			name:     "Removed",
			b:        NewBoolPtr(nil),
			expected: false,
		},
		{
			name:     "Set",
			b:        NewBool(!value),
			expected: !value,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if actual := testCase.b.Apply(value); actual != testCase.expected {
				t.Errorf("Expected: %v. Actual: %v", testCase.expected, actual)
			}
		})
	}
}

func TestBool_ApplyPtr(t *testing.T) {
	var (
		trueValue  = true
		falseValue = false
	)
	testCases := []struct {
		name     string
		b        Bool
		expected *bool
	}{
		{
			name:     "Unset",
			b:        Bool{},
			expected: &trueValue,
		},
		{
			name:     "Removed",
			b:        NewBoolPtr(nil),
			expected: nil,
		},
		{
			name:     "Set",
			b:        NewBool(falseValue),
			expected: &falseValue,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if actual := testCase.b.ApplyPtr(&trueValue); !reflect.DeepEqual(actual, testCase.expected) {
				t.Errorf("Expected: %v. Actual: %v", testCase.expected, actual)
			}
		})
	}
}

func TestBool_Diff(t *testing.T) {
	testCases := []struct {
		name     string
		b        Bool
		value    bool
		expected Bool
	}{
		{
			name:     "Unset",
			b:        Bool{},
			value:    true,
			expected: Bool{},
		},
		{
			name:     "Removed/NonZeroValue",
			b:        NewBoolPtr(nil),
			value:    true,
			expected: NewBoolPtr(nil),
		},
		{
			name:     "Removed/ZeroValue",
			b:        NewBoolPtr(nil),
			value:    false,
			expected: Bool{},
		},
		{
			name:     "Set/Equal",
			b:        NewBool(true),
			value:    true,
			expected: Bool{},
		},
		{
			name:     "Set/NotEqual",
			b:        NewBool(true),
			value:    false,
			expected: NewBool(true),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if actual := testCase.b.Diff(testCase.value); !reflect.DeepEqual(actual, testCase.expected) {
				t.Errorf("Expected: %v. Actual: %v", testCase.expected, actual)
			}
		})
	}
}

func TestBool_InterfaceValue(t *testing.T) {
	var b Bool
	if !reflect.ValueOf(b.InterfaceValue()).IsNil() {
		t.Errorf("Expected: nil, Actual: %v", b.InterfaceValue())
	}
	expected := true
	b.SetValue(expected)
	if !reflect.DeepEqual(b.InterfaceValue(), &expected) {
		t.Errorf("Expected: %v, Actual: %v", expected, b.InterfaceValue())
	}
}
