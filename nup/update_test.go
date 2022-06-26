package nup

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
			expected: Noop[int](),
		},
		{
			name:     "NullUpdate",
			json:     `{"update": null}`,
			expected: Remove[int](),
		},
		{
			name:     "ValueUpdate",
			json:     fmt.Sprintf(`{"update": %v}`, testValue),
			expected: Set(testValue),
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
			if dst.Update != testCase.expected {
				t.Errorf("Expected: %v. Actual: %v", testCase.expected, dst.Update)
			}
		})
	}
}

func TestRemoveOrSet(t *testing.T) {
	var (
		zero = 0
		one  = 1
	)
	testCases := []struct {
		name     string
		ptr      *int
		expected Update[int]
	}{
		{
			name:     "Nil",
			ptr:      nil,
			expected: Remove[int](),
		},
		{
			name:     "Zero",
			ptr:      &zero,
			expected: Set(0),
		},
		{
			name:     "Nonzero",
			ptr:      &one,
			expected: Set(1),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actual := RemoveOrSet(testCase.ptr)
			if actual != testCase.expected {
				t.Errorf("Expected: %v. Actual: %v", testCase.expected, actual)
			}
		})
	}
}

func TestUpdate_ValueOperation(t *testing.T) {
	testCases := []struct {
		name          string
		update        Update[int]
		expectedValue int
		expectedOp    Operation
	}{
		{
			name:          "Noop",
			update:        Noop[int](),
			expectedValue: 0,
			expectedOp:    OpNoop,
		},
		{
			name:          "Remove",
			update:        Remove[int](),
			expectedValue: 0,
			expectedOp:    OpRemove,
		},
		{
			name:          "Set",
			update:        Set(testValue),
			expectedValue: testValue,
			expectedOp:    OpSet,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actualValue, actualOp := testCase.update.ValueOperation()
			if actualValue != testCase.expectedValue {
				t.Errorf("Expected value: %v. Actual: %v", testCase.expectedValue, actualValue)
			}
			if actualOp != testCase.expectedOp {
				t.Errorf("Expected operation: %v. Actual: %v", testCase.expectedOp, actualOp)
			}
		})
	}
}

func TestUpdate_OperationAccessors(t *testing.T) {
	testCases := []struct {
		name             string
		update           Update[int]
		expectedOp       Operation
		expectedIsNoop   bool
		expectedIsRemove bool
		expectedIsSet    bool
		expectedIsChange bool
	}{
		{
			name:           "Noop",
			update:         Noop[int](),
			expectedOp:     OpNoop,
			expectedIsNoop: true,
		},
		{
			name:             "Remove",
			update:           Remove[int](),
			expectedOp:       OpRemove,
			expectedIsRemove: true,
			expectedIsChange: true,
		},
		{
			name:             "Set",
			update:           Set(testValue),
			expectedOp:       OpSet,
			expectedIsSet:    true,
			expectedIsChange: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			var (
				op       = testCase.update.Operation()
				isNoop   = testCase.update.IsNoop()
				isRemove = testCase.update.IsRemove()
				isSet    = testCase.update.IsSet()
				isChange = testCase.update.IsChange()
			)
			if op != testCase.expectedOp {
				t.Errorf("Expected Operation(): %v. Actual: %v", testCase.expectedOp, op)
			}
			if isNoop != testCase.expectedIsNoop {
				t.Errorf("Expected IsNoop(): %v. Actual: %v", testCase.expectedIsNoop, isNoop)
			}
			if isRemove != testCase.expectedIsRemove {
				t.Errorf("Expected IsRemove(): %v. Actual: %v", testCase.expectedIsRemove, isRemove)
			}
			if isSet != testCase.expectedIsSet {
				t.Errorf("Expected IsSet(): %v. Actual: %v", testCase.expectedIsSet, isSet)
			}
			if isChange != testCase.expectedIsChange {
				t.Errorf("Expected IsChange(): %v. Actual: %v", testCase.expectedIsChange, isChange)
			}
		})
	}
}

func TestUpdate_Value(t *testing.T) {
	testCases := []struct {
		name          string
		update        Update[int]
		expectedValue int
		expectedOK    bool
	}{
		{
			name:          "Noop",
			update:        Noop[int](),
			expectedValue: 0,
			expectedOK:    false,
		},
		{
			name:          "Remove",
			update:        Remove[int](),
			expectedValue: 0,
			expectedOK:    false,
		},
		{
			name:          "Set",
			update:        Set(testValue),
			expectedValue: testValue,
			expectedOK:    true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			value, isSet := testCase.update.Value()
			if value != testCase.expectedValue {
				t.Errorf("Expected value: %v. Actual: %v", testCase.expectedValue, value)
			}
			if isSet != testCase.expectedOK {
				t.Errorf("Expected isSet: %v. Actual: %v", testCase.expectedOK, isSet)
			}
		})
	}
}

func TestUpdate_ValueOrNil(t *testing.T) {
	testCases := []struct {
		name     string
		update   Update[int]
		expected *int
	}{
		{
			name:     "Noop",
			update:   Noop[int](),
			expected: nil,
		},
		{
			name:     "Remove",
			update:   Remove[int](),
			expected: nil,
		},
		{
			name:     "Set",
			update:   Set(testValue),
			expected: &testValue,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actual := testCase.update.ValueOrNil()
			if !reflect.DeepEqual(actual, testCase.expected) {
				t.Errorf("Expected value: %v. Actual: %v", testCase.expected, actual)
			}
		})
	}
}

func TestUpdate_Apply(t *testing.T) {
	value := 1
	testCases := []struct {
		name     string
		update   Update[int]
		expected int
	}{
		{
			name:     "Noop",
			update:   Noop[int](),
			expected: value,
		},
		{
			name:     "Remove",
			update:   Remove[int](),
			expected: 0,
		},
		{
			name:     "Set",
			update:   Set(value + 1),
			expected: value + 1,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if actual := testCase.update.Apply(value); actual != testCase.expected {
				t.Errorf("Expected: %v. Actual: %v", testCase.expected, actual)
			}
		})
	}
}

func TestUpdate_ApplyPtr(t *testing.T) {
	var (
		value1 = 1
		value2 = 2
	)
	testCases := []struct {
		name     string
		update   Update[int]
		expected *int
	}{
		{
			name:     "Noop",
			update:   Noop[int](),
			expected: &value1,
		},
		{
			name:     "Remove",
			update:   Remove[int](),
			expected: nil,
		},
		{
			name:     "Set",
			update:   Set(value2),
			expected: &value2,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if actual := testCase.update.ApplyPtr(&value1); !reflect.DeepEqual(actual, testCase.expected) {
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
		update   Update[int]
		value    int
		expected Update[int]
	}{
		{
			name:     "Noop",
			update:   Noop[int](),
			value:    value1,
			expected: Noop[int](),
		},
		{
			name:     "Remove/NonZeroValue",
			update:   Remove[int](),
			value:    value1,
			expected: Remove[int](),
		},
		{
			name:     "Remove/ZeroValue",
			update:   Remove[int](),
			value:    0,
			expected: Noop[int](),
		},
		{
			name:     "Set/Equal",
			update:   Set(value1),
			value:    value1,
			expected: Noop[int](),
		},
		{
			name:     "Set/NotEqual",
			update:   Set(value2),
			value:    value1,
			expected: Set(value2),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if actual := testCase.update.Diff(testCase.value); actual != testCase.expected {
				t.Errorf("Expected: %v. Actual: %v", testCase.expected, actual)
			}
		})
	}
}

func TestUpdate_DiffPtr(t *testing.T) {
	var (
		zero   = 0
		value1 = 1
		value2 = 2
	)
	testCases := []struct {
		name     string
		update   Update[int]
		value    *int
		expected Update[int]
	}{
		{
			name:     "Noop",
			update:   Noop[int](),
			value:    &value1,
			expected: Noop[int](),
		},
		{
			name:     "Remove/ZeroValue",
			update:   Remove[int](),
			value:    &zero,
			expected: Remove[int](),
		},
		{
			name:     "Remove/NonZeroValue",
			update:   Remove[int](),
			value:    &value1,
			expected: Remove[int](),
		},
		{
			name:     "Remove/Nil",
			update:   Remove[int](),
			value:    nil,
			expected: Noop[int](),
		},
		{
			name:     "Set/Equal",
			update:   Set(value1),
			value:    &value1,
			expected: Noop[int](),
		},
		{
			name:     "Set/NotEqual",
			update:   Set(value2),
			value:    &value1,
			expected: Set(value2),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if actual := testCase.update.DiffPtr(testCase.value); actual != testCase.expected {
				t.Errorf("Expected: %v. Actual: %v", testCase.expected, actual)
			}
		})
	}
}

func TestUpdate_IsSetTo(t *testing.T) {
	value := 1
	testCases := []struct {
		name     string
		update   Update[int]
		expected bool
	}{
		{
			name:     "Noop",
			update:   Noop[int](),
			expected: false,
		},
		{
			name:     "Remove",
			update:   Remove[int](),
			expected: false,
		},
		{
			name:     "Set/NotEqual",
			update:   Set(value + 1),
			expected: false,
		},
		{
			name:     "Set/Equal",
			update:   Set(value),
			expected: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if actual := testCase.update.IsSetTo(value); actual != testCase.expected {
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
		update   Update[int]
		expected bool
	}{
		{
			name:     "Noop",
			update:   Noop[int](),
			expected: false,
		},
		{
			name:     "Remove",
			update:   Remove[int](),
			expected: false,
		},
		{
			name:     "Set/NotSatisfied",
			update:   Set(1),
			expected: false,
		},
		{
			name:     "Set/Satisfied",
			update:   Set(-1),
			expected: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if actual := testCase.update.IsSetSuchThat(isNegative); actual != testCase.expected {
				t.Errorf("Expected: %v. Actual: %v", testCase.expected, actual)
			}
		})
	}
}

func TestUpdate_String(t *testing.T) {
	testCases := []struct {
		name     string
		update   fmt.Stringer
		expected string
	}{
		{
			name:     "Noop",
			update:   Noop[int](),
			expected: "<no-op>",
		},
		{
			name:     "Remove",
			update:   Remove[int](),
			expected: "<remove>",
		},
		{
			name:     "Set/Stringer",
			update:   Set(stringer{}),
			expected: "stringer",
		},
		{
			name:     "Set/String",
			update:   Set("value"),
			expected: "value",
		},
		{
			name:     "Set/Int",
			update:   Set(42),
			expected: "42",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if actual := testCase.update.String(); actual != testCase.expected {
				t.Errorf("Expected: %v. Actual: %v", testCase.expected, actual)
			}
		})
	}
}

type stringer struct{}

func (s stringer) String() string {
	return "stringer"
}

func TestUpdate_Equal(t *testing.T) {
	testCases := []struct {
		name     string
		first    Update[int]
		second   Update[int]
		expected bool
	}{
		{
			name:     "Equal/Noop",
			first:    Noop[int](),
			second:   Noop[int](),
			expected: true,
		},
		{
			name:     "Equal/Remove",
			first:    Remove[int](),
			second:   Remove[int](),
			expected: true,
		},
		{
			name:     "Equal/Set",
			first:    Set(testValue),
			second:   Set(testValue),
			expected: true,
		},
		{
			name:     "NotEqual/Noop/Remove",
			first:    Noop[int](),
			second:   Remove[int](),
			expected: false,
		},
		{
			name:     "NotEqual/Remove/Set",
			first:    Remove[int](),
			second:   Set(testValue),
			expected: false,
		},
		{
			name:     "NotEqual/Set/Noop",
			first:    Set(testValue),
			second:   Noop[int](),
			expected: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actual := testCase.first.Equal(testCase.second)
			if actual != testCase.expected {
				t.Errorf("Expected: %v. Actual: %v", testCase.expected, actual)
			}
		})
	}
}
