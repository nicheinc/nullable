package nully

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
			if !reflect.DeepEqual(dst.Update, testCase.expected) {
				t.Errorf("Expected: %v. Actual: %v", testCase.expected, dst.Update)
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
		},
		{
			name:          "Set",
			update:        Set(testValue),
			expectedOp:    OpSet,
			expectedIsSet: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			var (
				op       = testCase.update.Operation()
				isNoop   = testCase.update.IsNoop()
				isRemove = testCase.update.IsRemove()
				isSet    = testCase.update.IsSet()
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
			update:        Set(42),
			expectedValue: 42,
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
			value:    0.0,
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
			if actual := testCase.update.Diff(testCase.value); !reflect.DeepEqual(actual, testCase.expected) {
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
