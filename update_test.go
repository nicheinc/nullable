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
			expected: NewNoop[int](),
		},
		{
			name:     "NullUpdate",
			json:     `{"update": null}`,
			expected: NewRemove[int](),
		},
		{
			name:     "ValueUpdate",
			json:     fmt.Sprintf(`{"update": %v}`, testValue),
			expected: NewSet(testValue),
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
			update:         NewNoop[int](),
			expectedOp:     Noop,
			expectedIsNoop: true,
		},
		{
			name:             "Remove",
			update:           NewRemove[int](),
			expectedOp:       Remove,
			expectedIsRemove: true,
		},
		{
			name:          "Set",
			update:        NewSet(testValue),
			expectedOp:    Set,
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
			update:        NewNoop[int](),
			expectedValue: 0,
			expectedOK:    false,
		},
		{
			name:          "Remove",
			update:        NewRemove[int](),
			expectedValue: 0,
			expectedOK:    false,
		},
		{
			name:          "Set",
			update:        NewSet(42),
			expectedValue: 42,
			expectedOK:    true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			value, ok := testCase.update.Value()
			if value != testCase.expectedValue {
				t.Errorf("Expected value: %v. Actual: %v", testCase.expectedValue, value)
			}
			if ok != testCase.expectedOK {
				t.Errorf("Expected ok: %v. Actual: %v", testCase.expectedOK, ok)
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
			update:   NewNoop[int](),
			expected: value,
		},
		{
			name:     "Remove",
			update:   NewRemove[int](),
			expected: 0,
		},
		{
			name:     "Set",
			update:   NewSet(value + 1),
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
			update:   NewNoop[int](),
			value:    value1,
			expected: NewNoop[int](),
		},
		{
			name:     "Remove/NonZeroValue",
			update:   NewRemove[int](),
			value:    value1,
			expected: NewRemove[int](),
		},
		{
			name:     "Remove/ZeroValue",
			update:   NewRemove[int](),
			value:    0.0,
			expected: NewNoop[int](),
		},
		{
			name:     "Set/Equal",
			update:   NewSet(value1),
			value:    value1,
			expected: NewNoop[int](),
		},
		{
			name:     "Set/NotEqual",
			update:   NewSet(value2),
			value:    value1,
			expected: NewSet(value2),
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
			update:   NewNoop[int](),
			expected: false,
		},
		{
			name:     "Remove",
			update:   NewRemove[int](),
			expected: false,
		},
		{
			name:     "Set/NotEqual",
			update:   NewSet(value + 1),
			expected: false,
		},
		{
			name:     "Set/Equal",
			update:   NewSet(value),
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
			update:   NewNoop[int](),
			expected: false,
		},
		{
			name:     "Remove",
			update:   NewRemove[int](),
			expected: false,
		},
		{
			name:     "Set/NotSatisfied",
			update:   NewSet(1),
			expected: false,
		},
		{
			name:     "Set/Satisfied",
			update:   NewSet(-1),
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
			update:   NewNoop[int](),
			expected: "<no-op>",
		},
		{
			name:     "Remove",
			update:   NewRemove[int](),
			expected: "<remove>",
		},
		{
			name:     "Set/Stringer",
			update:   NewSet(stringer{}),
			expected: "stringer",
		},
		{
			name:     "Set/String",
			update:   NewSet("value"),
			expected: "value",
		},
		{
			name:     "Set/Int",
			update:   NewSet(42),
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
