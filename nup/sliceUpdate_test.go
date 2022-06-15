package nup

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

var (
	testSlice1 = []int{1}
	testSlice2 = []int{1, 2}
)

// Ensure implementation of the updateMarshaller interface.
var _ updateMarshaller = &SliceUpdate[int]{}

func TestSliceUpdate_UnmarshalJSON(t *testing.T) {
	testCases := []struct {
		name     string
		json     string
		expected SliceUpdate[int]
	}{
		{
			name:     "EmptyJSONObject",
			json:     `{}`,
			expected: SliceNoop[int](),
		},
		{
			name:     "NullUpdate",
			json:     `{"update": null}`,
			expected: SliceRemove[int](),
		},
		{
			name:     "EmptyUpdate",
			json:     `{"update": []}`,
			expected: SliceRemoveOrSet([]int{}),
		},
		{
			name:     "NonemptyUpdate",
			json:     fmt.Sprintf(`{"update": %v}`, testSlice1),
			expected: SliceRemoveOrSet(testSlice1),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			var dst struct {
				Update SliceUpdate[int] `json:"update"`
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

func TestSliceRemoveOrSet(t *testing.T) {
	testCases := []struct {
		name     string
		value    []int
		expected SliceUpdate[int]
	}{
		{
			name:     "Nil",
			value:    nil,
			expected: SliceRemove[int](),
		},
		{
			name:     "EmptyNonNil",
			value:    []int{},
			expected: SliceRemoveOrSet([]int{}),
		},
		{
			name:     "Nonempty",
			value:    testSlice1,
			expected: SliceRemoveOrSet(testSlice1),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actual := SliceRemoveOrSet(testCase.value)
			if !reflect.DeepEqual(actual, testCase.expected) {
				t.Errorf("Expected: %v. Actual: %v", testCase.expected, actual)
			}
		})
	}
}

func TestSliceUpdate_ValueOperation(t *testing.T) {
	testCases := []struct {
		name          string
		update        SliceUpdate[int]
		expectedValue []int
		expectedOp    Operation
	}{
		{
			name:          "Noop",
			update:        SliceNoop[int](),
			expectedValue: nil,
			expectedOp:    OpNoop,
		},
		{
			name:          "Remove",
			update:        SliceRemove[int](),
			expectedValue: nil,
			expectedOp:    OpRemove,
		},
		{
			name:          "Set",
			update:        SliceRemoveOrSet(testSlice1),
			expectedValue: testSlice1,
			expectedOp:    OpSet,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actualValue, actualOp := testCase.update.ValueOperation()
			if !reflect.DeepEqual(actualValue, testCase.expectedValue) {
				t.Errorf("Expected value: %v. Actual: %v", testCase.expectedValue, actualValue)
			}
			if actualOp != testCase.expectedOp {
				t.Errorf("Expected operation: %v. Actual: %v", testCase.expectedOp, actualOp)
			}
		})
	}
}

func TestSliceUpdate_OperationAccessors(t *testing.T) {
	testCases := []struct {
		name             string
		update           SliceUpdate[int]
		expectedOp       Operation
		expectedIsNoop   bool
		expectedIsRemove bool
		expectedIsSet    bool
		expectedIsChange bool
	}{
		{
			name:           "Noop",
			update:         SliceNoop[int](),
			expectedOp:     OpNoop,
			expectedIsNoop: true,
		},
		{
			name:             "Remove",
			update:           SliceRemove[int](),
			expectedOp:       OpRemove,
			expectedIsRemove: true,
			expectedIsChange: true,
		},
		{
			name:             "Set",
			update:           SliceRemoveOrSet(testSlice1),
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

func TestSliceUpdate_Value(t *testing.T) {
	testCases := []struct {
		name          string
		update        SliceUpdate[int]
		expectedValue []int
		expectedIsSet bool
	}{
		{
			name:          "Noop",
			update:        SliceNoop[int](),
			expectedValue: nil,
			expectedIsSet: false,
		},
		{
			name:          "Remove",
			update:        SliceRemove[int](),
			expectedValue: nil,
			expectedIsSet: false,
		},
		{
			name:          "Set",
			update:        SliceRemoveOrSet(testSlice1),
			expectedValue: testSlice1,
			expectedIsSet: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			value, isSet := testCase.update.Value()
			if !reflect.DeepEqual(value, testCase.expectedValue) {
				t.Errorf("Expected value: %v. Actual: %v", testCase.expectedValue, value)
			}
			if isSet != testCase.expectedIsSet {
				t.Errorf("Expected isSet: %v. Actual: %v", testCase.expectedIsSet, isSet)
			}
		})
	}
}

func TestSliceUpdate_ValueOrNil(t *testing.T) {
	testCases := []struct {
		name     string
		update   SliceUpdate[int]
		expected []int
	}{
		{
			name:     "Noop",
			update:   SliceNoop[int](),
			expected: nil,
		},
		{
			name:     "Remove",
			update:   SliceRemove[int](),
			expected: nil,
		},
		{
			name:     "Set",
			update:   SliceRemoveOrSet(testSlice1),
			expected: testSlice1,
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

func TestSliceUpdate_Apply(t *testing.T) {
	testCases := []struct {
		name     string
		u        SliceUpdate[int]
		expected []int
	}{
		{
			name:     "Noop",
			u:        SliceNoop[int](),
			expected: testSlice1,
		},
		{
			name:     "Remove",
			u:        SliceRemove[int](),
			expected: nil,
		},
		{
			name:     "Set",
			u:        SliceRemoveOrSet(testSlice2),
			expected: testSlice2,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if actual := testCase.u.Apply(testSlice1); !reflect.DeepEqual(actual, testCase.expected) {
				t.Errorf("Expected: %v. Actual: %v", testCase.expected, actual)
			}
		})
	}
}

func TestSliceUpdate_Diff(t *testing.T) {
	testCases := []struct {
		name     string
		u        SliceUpdate[int]
		value    []int
		expected SliceUpdate[int]
	}{
		{
			name:     "Noop",
			u:        SliceNoop[int](),
			value:    testSlice1,
			expected: SliceNoop[int](),
		},
		{
			name:     "Remove/NonZeroValue",
			u:        SliceRemove[int](),
			value:    testSlice1,
			expected: SliceRemove[int](),
		},
		{
			name:     "Remove/ZeroValue",
			u:        SliceRemove[int](),
			value:    nil,
			expected: SliceNoop[int](),
		},
		{
			name:     "Set/Equal",
			u:        SliceRemoveOrSet(testSlice1),
			value:    testSlice1,
			expected: SliceNoop[int](),
		},
		{
			name:     "Set/NotEqual",
			u:        SliceRemoveOrSet(testSlice2),
			value:    testSlice1,
			expected: SliceRemoveOrSet(testSlice2),
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

func TestSliceUpdate_IsSetTo(t *testing.T) {
	testCases := []struct {
		name     string
		u        SliceUpdate[int]
		expected bool
	}{
		{
			name:     "Noop",
			u:        SliceNoop[int](),
			expected: false,
		},
		{
			name:     "Remove",
			u:        SliceRemove[int](),
			expected: false,
		},
		{
			name:     "Set/NotEqual/DifferentSizes",
			u:        SliceRemoveOrSet(testSlice2),
			expected: false,
		},
		{
			name:     "Set/NotEqual/SameSize",
			u:        SliceRemoveOrSet([]int{0}),
			expected: false,
		},
		{
			name:     "Set/Equal",
			u:        SliceRemoveOrSet(testSlice1),
			expected: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if actual := testCase.u.IsSetTo(testSlice1); actual != testCase.expected {
				t.Errorf("Expected: %v. Actual: %v", testCase.expected, actual)
			}
		})
	}
}

func TestSliceUpdate_IsSetSuchThat(t *testing.T) {
	isCouple := func(v []int) bool {
		return len(v) == 2
	}
	testCases := []struct {
		name     string
		u        SliceUpdate[int]
		expected bool
	}{
		{
			name:     "Noop",
			u:        SliceNoop[int](),
			expected: false,
		},
		{
			name:     "Remove",
			u:        SliceRemove[int](),
			expected: false,
		},
		{
			name:     "Set/NotSatisfied",
			u:        SliceRemoveOrSet(testSlice1),
			expected: false,
		},
		{
			name:     "Set/Satisfied",
			u:        SliceRemoveOrSet(testSlice2),
			expected: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if actual := testCase.u.IsSetSuchThat(isCouple); actual != testCase.expected {
				t.Errorf("Expected: %v. Actual: %v", testCase.expected, actual)
			}
		})
	}
}

func TestSliceUpdate_String(t *testing.T) {
	testCases := []struct {
		name     string
		u        fmt.Stringer
		expected string
	}{
		{
			name:     "Noop",
			u:        SliceNoop[int](),
			expected: "<no-op>",
		},
		{
			name:     "Remove",
			u:        SliceRemove[int](),
			expected: "<remove>",
		},
		{
			name:     "RemoveOrSet/Nil",
			u:        SliceRemoveOrSet([]stringer(nil)),
			expected: "<remove>",
		},
		{
			name:     "RemoveOrSet/Empty",
			u:        SliceRemoveOrSet([]stringer{}),
			expected: "[]",
		},
		{
			name:     "RemoveOrSet/StringerSlice",
			u:        SliceRemoveOrSet([]stringer{{}}),
			expected: "[stringer]",
		},
		{
			name:     "RemoveOrSet/IntSlice",
			u:        SliceRemoveOrSet([]int{42}),
			expected: "[42]",
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

func TestSliceUpdate_Equal(t *testing.T) {
	testCases := []struct {
		name     string
		first    SliceUpdate[int]
		second   SliceUpdate[int]
		expected bool
	}{
		{
			name:     "Equal/Noop",
			first:    SliceNoop[int](),
			second:   SliceNoop[int](),
			expected: true,
		},
		{
			name:     "Equal/Remove",
			first:    SliceRemove[int](),
			second:   SliceRemove[int](),
			expected: true,
		},
		{
			name:     "Equal/Set",
			first:    SliceRemoveOrSet(testSlice1),
			second:   SliceRemoveOrSet(testSlice1),
			expected: true,
		},
		{
			name:     "NotEqual/Noop/Remove",
			first:    SliceNoop[int](),
			second:   SliceRemove[int](),
			expected: false,
		},
		{
			name:     "NotEqual/Remove/Set",
			first:    SliceRemove[int](),
			second:   SliceRemoveOrSet(testSlice1),
			expected: false,
		},
		{
			name:     "NotEqual/Set/Noop",
			first:    SliceRemoveOrSet(testSlice1),
			second:   SliceNoop[int](),
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
