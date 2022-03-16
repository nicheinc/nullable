package nullable

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

var testValue = 42

// Ensure implementation of the updateMarshaller interface.
var _ updateMarshaller = &Update[int]{}

func TestUpdate_UnmarshalJSON(t *testing.T) {
	testCases := []struct {
		name     string
		json     string
		expected Update[int]
	}{
		{
			name:     "EmptyJSONObject",
			json:     `{}`,
			expected: Update[int]{},
		},
		{
			name:     "NullUpdate",
			json:     `{"update": null}`,
			expected: NewUpdatePtr[int](nil),
		},
		{
			name:     "ValueUpdate",
			json:     fmt.Sprintf(`{"update": %v}`, testValue),
			expected: NewUpdate(testValue),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			var dst struct {
				Update Update[int] `json:"update"`
			}
			if err := json.Unmarshal([]byte(testCase.json), &dst); err != nil {
				t.Errorf("Error unmarshaling JSON: %s", err)
			}
			if !reflect.DeepEqual(dst.Update, testCase.expected) {
				t.Errorf("Expected: %v. Actual: %v", testCase.expected, dst.Update)
			}
		})
	}
}

func TestUpdate_SetValue(t *testing.T) {
	testCases := []struct {
		name     string
		value    *int
		expected Update[int]
	}{
		{
			name:     "NullUpdate",
			value:    nil,
			expected: NewUpdatePtr[int](nil),
		},
		{
			name:     "ValueUpdate",
			value:    &testValue,
			expected: NewUpdate(testValue),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actual := Update[int]{}
			if testCase.value != nil {
				actual.SetValue(*testCase.value)
			} else {
				actual.SetPtr(nil)
			}
			if !reflect.DeepEqual(actual, testCase.expected) {
				t.Errorf("Expected: %v. Actual: %v", testCase.expected, actual)
			}
		})
	}
}

func TestUpdate_Removed(t *testing.T) {
	testCases := []struct {
		name     string
		u        Update[int]
		expected bool
	}{
		{
			name:     "NotSet",
			u:        Update[int]{},
			expected: false,
		},
		{
			name:     "NullUpdate",
			u:        NewUpdatePtr[int](nil),
			expected: true,
		},
		{
			name:     "ValueUpdate",
			u:        NewUpdate(testValue),
			expected: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actual := testCase.u.Removed()
			if !reflect.DeepEqual(actual, testCase.expected) {
				t.Errorf("Expected: %v. Actual: %v", testCase.expected, actual)
			}
		})
	}
}

func TestUpdate_Value(t *testing.T) {
	var u Update[int]
	if u.Value() != nil {
		t.Errorf("Expected: nil. Actual: %v", u.Value())
	}
	expected := 1
	u.SetValue(expected)
	if *u.Value() != expected {
		t.Errorf("Expected: %v. Actual: %v", expected, *u.Value())
	}
}

func TestUpdate_Apply(t *testing.T) {
	value := 1
	testCases := []struct {
		name     string
		u        Update[int]
		expected int
	}{
		{
			name:     "Unset",
			u:        Update[int]{},
			expected: value,
		},
		{
			name:     "Removed",
			u:        NewUpdatePtr[int](nil),
			expected: 0,
		},
		{
			name:     "Set",
			u:        NewUpdate(value + 1),
			expected: value + 1,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if actual := testCase.u.Apply(value); actual != testCase.expected {
				t.Errorf("Expected: %v. Actual: %v", testCase.expected, actual)
			}
		})
	}
}

func TestUpdate_ApplyPtr(t *testing.T) {
	var (
		value1 = 1
		value2 = value1 + 1
	)
	testCases := []struct {
		name     string
		u        Update[int]
		expected *int
	}{
		{
			name:     "Unset",
			u:        Update[int]{},
			expected: &value1,
		},
		{
			name:     "Removed",
			u:        NewUpdatePtr[int](nil),
			expected: nil,
		},
		{
			name:     "Set",
			u:        NewUpdate(value2),
			expected: &value2,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if actual := testCase.u.ApplyPtr(&value1); !reflect.DeepEqual(actual, testCase.expected) {
				t.Errorf("Expected: %v. Actual: %v", testCase.expected, actual)
			}
		})
	}
}

func TestUpdate_Diff(t *testing.T) {
	var (
		value1 = 1
		value2 = 2
	)
	testCases := []struct {
		name     string
		u        Update[int]
		value    int
		expected Update[int]
	}{
		{
			name:     "Unset",
			u:        Update[int]{},
			value:    value1,
			expected: Update[int]{},
		},
		{
			name:     "Removed/NonZeroValue",
			u:        NewUpdatePtr[int](nil),
			value:    value1,
			expected: NewUpdatePtr[int](nil),
		},
		{
			name:     "Removed/ZeroValue",
			u:        NewUpdatePtr[int](nil),
			value:    0.0,
			expected: Update[int]{},
		},
		{
			name:     "Set/Equal",
			u:        NewUpdate(value1),
			value:    value1,
			expected: Update[int]{},
		},
		{
			name:     "Set/NotEqual",
			u:        NewUpdate(value2),
			value:    value1,
			expected: NewUpdate(value2),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if actual := testCase.u.Diff(testCase.value); !reflect.DeepEqual(actual, testCase.expected) {
				t.Errorf("Expected: %v. Actual: %v", testCase.expected, actual)
			}
		})
	}
}

func TestUpdate_IsSetTo(t *testing.T) {
	value := 1
	testCases := []struct {
		name     string
		u        Update[int]
		expected bool
	}{
		{
			name:     "Unset",
			u:        Update[int]{},
			expected: false,
		},
		{
			name:     "Removed",
			u:        NewUpdatePtr[int](nil),
			expected: false,
		},
		{
			name:     "Set/NotEqual",
			u:        NewUpdate(value + 1),
			expected: false,
		},
		{
			name:     "Set/Equal",
			u:        NewUpdate(value),
			expected: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if actual := testCase.u.IsSetTo(value); actual != testCase.expected {
				t.Errorf("Expected: %v. Actual: %v", testCase.expected, actual)
			}
		})
	}
}

func TestUpdate_IsSetSuchThat(t *testing.T) {
	isNegative := func(v int) bool {
		return v < 0
	}
	testCases := []struct {
		name     string
		u        Update[int]
		expected bool
	}{
		{
			name:     "Unset",
			u:        Update[int]{},
			expected: false,
		},
		{
			name:     "Removed",
			u:        NewUpdatePtr[int](nil),
			expected: false,
		},
		{
			name:     "Set/NotSatisfied",
			u:        NewUpdate(1),
			expected: false,
		},
		{
			name:     "Set/Satisfied",
			u:        NewUpdate(-1),
			expected: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if actual := testCase.u.IsSetSuchThat(isNegative); actual != testCase.expected {
				t.Errorf("Expected: %v. Actual: %v", testCase.expected, actual)
			}
		})
	}
}

func TestUpdate_String(t *testing.T) {
	testCases := []struct {
		name     string
		u        fmt.Stringer
		expected string
	}{
		{
			name:     "Unset",
			u:        Update[string]{},
			expected: "<unset>",
		},
		{
			name:     "Removed",
			u:        NewUpdatePtr[string](nil),
			expected: "<removed>",
		},
		{
			name:     "Set/Stringer",
			u:        NewUpdate(stringer{}),
			expected: "stringer",
		},
		{
			name:     "Set/String",
			u:        NewUpdate("value"),
			expected: "value",
		},
		{
			name:     "Set/Int",
			u:        NewUpdate(42),
			expected: "42",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if actual := testCase.u.String(); actual != testCase.expected {
				t.Errorf("Expected: %v. Actual: %v", testCase.expected, actual)
			}
		})
	}
}

type stringer struct{}

func (s stringer) String() string {
	return "stringer"
}
