package nup

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/nicheinc/expect"
)

var testValue = 42

// Ensure implementation of the updateMarshaller interface.
var _ updateMarshaller = &Update[int]{}

func TestUpdate_MarshalJSON(t *testing.T) {
	type testCase struct {
		update   Update[int]
		expected string
	}
	run := func(name string, testCase testCase) {
		t.Helper()
		t.Run(name, func(t *testing.T) {
			t.Helper()
			actual, err := json.Marshal(testCase.update)
			expect.ErrorNil(t, err)
			expect.Equal(t, string(actual), testCase.expected)
		})
	}

	run("Noop", testCase{
		update:   Noop[int](),
		expected: "null",
	})
	run("Remove", testCase{
		update:   Remove[int](),
		expected: "null",
	})
	run("Set", testCase{
		update:   Set[int](testValue),
		expected: "42",
	})
}

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
			err := json.Unmarshal([]byte(testCase.json), &dst)
			expect.ErrorNil(t, err)
			expect.Equal(t, dst.Update, testCase.expected)
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
			expect.Equal(t, actual, testCase.expected)
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
			expect.Equal(t, actualValue, testCase.expectedValue)
			expect.Equal(t, actualOp, testCase.expectedOp)
		})
	}
}

func TestUpdate_OperationAccessors(t *testing.T) {
	testCases := []struct {
		name             string
		update           Update[int]
		expectedOp       Operation
		expectedIsNoop   bool
		expectedIsZero   bool
		expectedIsRemove bool
		expectedIsSet    bool
		expectedIsChange bool
	}{
		{
			name:           "Noop",
			update:         Noop[int](),
			expectedOp:     OpNoop,
			expectedIsNoop: true,
			expectedIsZero: true,
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
			expect.Equal(t, testCase.update.Operation(), testCase.expectedOp)
			expect.Equal(t, testCase.update.IsNoop(), testCase.expectedIsNoop)
			expect.Equal(t, testCase.update.IsZero(), testCase.expectedIsZero)
			expect.Equal(t, testCase.update.IsRemove(), testCase.expectedIsRemove)
			expect.Equal(t, testCase.update.IsSet(), testCase.expectedIsSet)
			expect.Equal(t, testCase.update.IsChange(), testCase.expectedIsChange)
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
			expect.Equal(t, value, testCase.expectedValue)
			expect.Equal(t, isSet, testCase.expectedOK)
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
			expect.Equal(t, actual, testCase.expected)
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
			actual := testCase.update.Apply(value)
			expect.Equal(t, actual, testCase.expected)
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
			actual := testCase.update.ApplyPtr(&value1)
			expect.Equal(t, actual, testCase.expected)
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
			actual := testCase.update.Diff(testCase.value)
			expect.Equal(t, actual, testCase.expected)
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
			actual := testCase.update.DiffPtr(testCase.value)
			expect.Equal(t, actual, testCase.expected)
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
			actual := testCase.update.IsSetTo(value)
			expect.Equal(t, actual, testCase.expected)
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
			actual := testCase.update.IsSetSuchThat(isNegative)
			expect.Equal(t, actual, testCase.expected)
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
			actual := testCase.update.String()
			expect.Equal(t, actual, testCase.expected)
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
			expect.Equal(t, actual, testCase.expected)
		})
	}
}
