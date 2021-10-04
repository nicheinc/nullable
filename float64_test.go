package nullable

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

const testFloat64 = 4.2

// Ensure implementation of Nullable interface.
var _ Nullable = &Float64{}

func TestFloat64_UnmarshalJSON(t *testing.T) {
	testCases := []struct {
		name     string
		json     string
		expected Float64
	}{
		{
			name:     "EmptyJSONObject",
			json:     `{}`,
			expected: Float64{},
		},
		{
			name:     "NullFloat64",
			json:     `{"float64": null}`,
			expected: NewFloat64Ptr(nil),
		},
		{
			name:     "ValueFloat64",
			json:     fmt.Sprintf(`{"float64": %v}`, testFloat64),
			expected: NewFloat64(testFloat64),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			var dst struct {
				Float64 Float64 `json:"float64"`
			}
			if err := json.Unmarshal([]byte(testCase.json), &dst); err != nil {
				t.Errorf("Error unmarshaling JSON: %s", err)
			}
			if !reflect.DeepEqual(dst.Float64, testCase.expected) {
				t.Errorf("Expected: %v. Actual: %v", testCase.expected, dst.Float64)
			}
		})
	}
}

func TestFloat64_SetValue(t *testing.T) {
	testCases := []struct {
		name     string
		value    *float64
		expected Float64
	}{
		{
			name:     "NullFloat64",
			value:    nil,
			expected: NewFloat64Ptr(nil),
		},
		{
			name:     "ValueFloat64",
			value:    func(v float64) *float64 { return &v }(testFloat64),
			expected: NewFloat64(testFloat64),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actual := Float64{}
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

func TestFloat64_Removed(t *testing.T) {
	testCases := []struct {
		name     string
		f        Float64
		expected bool
	}{
		{
			name:     "NotSet",
			f:        Float64{},
			expected: false,
		},
		{
			name:     "NullFloat64",
			f:        NewFloat64Ptr(nil),
			expected: true,
		},
		{
			name:     "ValueFloat64",
			f:        NewFloat64(testFloat64),
			expected: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actual := testCase.f.Removed()
			if !reflect.DeepEqual(actual, testCase.expected) {
				t.Errorf("Expected: %v, Actual: %v", testCase.expected, actual)
			}
		})
	}
}

func TestFloat64_IsZero(t *testing.T) {
	testCases := []struct {
		name     string
		f        Float64
		expected bool
	}{
		{
			name:     "NotSet",
			f:        Float64{},
			expected: false,
		},
		{
			name:     "NullFloat64",
			f:        NewFloat64Ptr(nil),
			expected: false,
		},
		{
			name:     "ZeroFloat64",
			f:        NewFloat64(0.0),
			expected: true,
		},
		{
			name:     "NonZeroFloat64",
			f:        NewFloat64(1.0),
			expected: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actual := testCase.f.IsZero()
			if !reflect.DeepEqual(actual, testCase.expected) {
				t.Errorf("Expected: %v, Actual: %v", testCase.expected, actual)
			}
		})
	}
}

func TestFloat64_IsNegative(t *testing.T) {
	testCases := []struct {
		name     string
		f        Float64
		expected bool
	}{
		{
			name:     "NotSet",
			f:        Float64{},
			expected: false,
		},
		{
			name:     "NullFloat64",
			f:        NewFloat64Ptr(nil),
			expected: false,
		},
		{
			name:     "NegativeFloat64",
			f:        NewFloat64(-1.0),
			expected: true,
		},
		{
			name:     "PositiveFloat64",
			f:        NewFloat64(1.0),
			expected: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actual := testCase.f.IsNegative()
			if !reflect.DeepEqual(actual, testCase.expected) {
				t.Errorf("Expected: %v, Actual: %v", testCase.expected, actual)
			}
		})
	}
}

func TestFloat64_Value(t *testing.T) {
	var f Float64
	if f.Value() != nil {
		t.Errorf("Expected: nil, Actual: %v", f.Value())
	}
	expected := 1.5
	f.SetValue(expected)
	if *f.Value() != expected {
		t.Errorf("Expected: %v, Actual: %v", expected, *f.Value())
	}
}

func TestFloat64_Equals(t *testing.T) {
	value := 1.5
	testCases := []struct {
		name     string
		f        Float64
		expected bool
	}{
		{
			name:     "Unset",
			f:        Float64{},
			expected: false,
		},
		{
			name:     "Removed",
			f:        NewFloat64Ptr(nil),
			expected: false,
		},
		{
			name:     "Set/NotEqual",
			f:        NewFloat64(value + 1),
			expected: false,
		},
		{
			name:     "Set/Equal",
			f:        NewFloat64(value),
			expected: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if actual := testCase.f.Equals(value); actual != testCase.expected {
				t.Errorf("Expected: %v. Actual: %v", testCase.expected, actual)
			}
		})
	}
}

func TestFloat64_Apply(t *testing.T) {
	value := 1.5
	testCases := []struct {
		name     string
		f        Float64
		expected float64
	}{
		{
			name:     "Unset",
			f:        Float64{},
			expected: value,
		},
		{
			name:     "Removed",
			f:        NewFloat64Ptr(nil),
			expected: 0,
		},
		{
			name:     "Set",
			f:        NewFloat64(value + 1),
			expected: value + 1,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if actual := testCase.f.Apply(value); actual != testCase.expected {
				t.Errorf("Expected: %v. Actual: %v", testCase.expected, actual)
			}
		})
	}
}

func TestFloat64_ApplyPtr(t *testing.T) {
	var (
		value1 = 1.5
		value2 = value1 + 1
	)
	testCases := []struct {
		name     string
		f        Float64
		expected *float64
	}{
		{
			name:     "Unset",
			f:        Float64{},
			expected: &value1,
		},
		{
			name:     "Removed",
			f:        NewFloat64Ptr(nil),
			expected: nil,
		},
		{
			name:     "Set",
			f:        NewFloat64(value2),
			expected: &value2,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if actual := testCase.f.ApplyPtr(&value1); !reflect.DeepEqual(actual, testCase.expected) {
				t.Errorf("Expected: %v. Actual: %v", testCase.expected, actual)
			}
		})
	}
}

func TestFloat64_InterfaceValue(t *testing.T) {
	var f Float64
	if !reflect.ValueOf(f.InterfaceValue()).IsNil() {
		t.Errorf("Expected: nil, Actual: %v", f.InterfaceValue())
	}
	expected := 1.5
	f.SetValue(expected)
	if !reflect.DeepEqual(f.InterfaceValue(), &expected) {
		t.Errorf("Expected: %v, Actual: %v", expected, f.InterfaceValue())
	}
}
