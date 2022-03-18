package nully

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
			expected: SliceSet([]int{}),
		},
		{
			name:     "NonemptyUpdate",
			json:     fmt.Sprintf(`{"update": %v}`, testSlice1),
			expected: SliceSet(testSlice1),
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
			update:           SliceSet(testSlice1),
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
			update:        SliceSet(testSlice1),
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
			u:        SliceSet(testSlice2),
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
			u:        SliceSet(testSlice1),
			value:    testSlice1,
			expected: SliceNoop[int](),
		},
		{
			name:     "Set/NotEqual",
			u:        SliceSet(testSlice2),
			value:    testSlice1,
			expected: SliceSet(testSlice2),
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
			u:        SliceSet[int](nil),
			expected: false,
		},
		{
			name:     "Set/NotEqual/DifferentSizes",
			u:        SliceSet(testSlice2),
			expected: false,
		},
		{
			name:     "Set/NotEqual/SameSize",
			u:        SliceSet([]int{0}),
			expected: false,
		},
		{
			name:     "Set/Equal",
			u:        SliceSet(testSlice1),
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
			u:        SliceSet[int](nil),
			expected: false,
		},
		{
			name:     "Set/NotSatisfied",
			u:        SliceSet(testSlice1),
			expected: false,
		},
		{
			name:     "Set/Satisfied",
			u:        SliceSet(testSlice2),
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
			name:     "Set/StringerSlice",
			u:        SliceSet([]stringer{{}}),
			expected: "[stringer]",
		},
		{
			name:     "Set/IntSlice",
			u:        SliceSet([]int{42}),
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
