package nullable

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
			expected: NewNoopSlice[int](),
		},
		{
			name:     "NullUpdate",
			json:     `{"update": null}`,
			expected: NewRemoveSlice[int](),
		},
		{
			name:     "EmptyUpdate",
			json:     `{"update": []}`,
			expected: NewSetSlice([]int{}),
		},
		{
			name:     "NonemptyUpdate",
			json:     fmt.Sprintf(`{"update": %v}`, testSlice1),
			expected: NewSetSlice(testSlice1),
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
	}{
		{
			name:           "Noop",
			update:         NewNoopSlice[int](),
			expectedOp:     Noop,
			expectedIsNoop: true,
		},
		{
			name:             "Remove",
			update:           NewRemoveSlice[int](),
			expectedOp:       Remove,
			expectedIsRemove: true,
		},
		{
			name:          "Set",
			update:        NewSetSlice(testSlice1),
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

func TestSliceUpdate_Value(t *testing.T) {
	testCases := []struct {
		name          string
		update        SliceUpdate[int]
		expectedValue []int
		expectedOK    bool
	}{
		{
			name:          "Noop",
			update:        NewNoopSlice[int](),
			expectedValue: nil,
			expectedOK:    false,
		},
		{
			name:          "Remove",
			update:        NewRemoveSlice[int](),
			expectedValue: nil,
			expectedOK:    false,
		},
		{
			name:          "Set",
			update:        NewSetSlice(testSlice1),
			expectedValue: testSlice1,
			expectedOK:    true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			value, ok := testCase.update.Value()
			if !reflect.DeepEqual(value, testCase.expectedValue) {
				t.Errorf("Expected value: %v. Actual: %v", testCase.expectedValue, value)
			}
			if ok != testCase.expectedOK {
				t.Errorf("Expected ok: %v. Actual: %v", testCase.expectedOK, ok)
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
			u:        NewNoopSlice[int](),
			expected: testSlice1,
		},
		{
			name:     "Remove",
			u:        NewRemoveSlice[int](),
			expected: nil,
		},
		{
			name:     "Set",
			u:        NewSetSlice(testSlice2),
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
			u:        NewNoopSlice[int](),
			value:    testSlice1,
			expected: NewNoopSlice[int](),
		},
		{
			name:     "Remove/NonZeroValue",
			u:        NewRemoveSlice[int](),
			value:    testSlice1,
			expected: NewRemoveSlice[int](),
		},
		{
			name:     "Remove/ZeroValue",
			u:        NewRemoveSlice[int](),
			value:    nil,
			expected: NewNoopSlice[int](),
		},
		{
			name:     "Set/Equal",
			u:        NewSetSlice(testSlice1),
			value:    testSlice1,
			expected: NewNoopSlice[int](),
		},
		{
			name:     "Set/NotEqual",
			u:        NewSetSlice(testSlice2),
			value:    testSlice1,
			expected: NewSetSlice(testSlice2),
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
			u:        NewNoopSlice[int](),
			expected: false,
		},
		{
			name:     "Remove",
			u:        NewSetSlice[int](nil),
			expected: false,
		},
		{
			name:     "Set/NotEqual/DifferentSizes",
			u:        NewSetSlice(testSlice2),
			expected: false,
		},
		{
			name:     "Set/NotEqual/SameSize",
			u:        NewSetSlice([]int{0}),
			expected: false,
		},
		{
			name:     "Set/Equal",
			u:        NewSetSlice(testSlice1),
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
			u:        NewNoopSlice[int](),
			expected: false,
		},
		{
			name:     "Remove",
			u:        NewSetSlice[int](nil),
			expected: false,
		},
		{
			name:     "Set/NotSatisfied",
			u:        NewSetSlice(testSlice1),
			expected: false,
		},
		{
			name:     "Set/Satisfied",
			u:        NewSetSlice(testSlice2),
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
			u:        SliceUpdate[string]{},
			expected: "<no-op>",
		},
		{
			name:     "Remove",
			u:        NewRemoveSlice[string](),
			expected: "<remove>",
		},
		{
			name:     "Set/StringerSlice",
			u:        NewSetSlice([]stringer{{}}),
			expected: "[stringer]",
		},
		{
			name:     "Set/IntSlice",
			u:        NewSetSlice([]int{42}),
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
