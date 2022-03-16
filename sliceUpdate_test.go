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
			expected: SliceUpdate[int]{},
		},
		{
			name:     "NullSlice",
			json:     `{"update": null}`,
			expected: NewSliceUpdate[int](nil),
		},
		{
			name:     "EmptySlice",
			json:     `{"update": null}`,
			expected: NewSliceUpdate[int](nil),
		},
		{
			name:     "NonemptySlice",
			json:     fmt.Sprintf(`{"update": %v}`, testSlice1),
			expected: NewSliceUpdate(testSlice1),
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

func TestSliceUpdate_SetValue(t *testing.T) {
	testCases := []struct {
		name     string
		value    []int
		expected SliceUpdate[int]
	}{
		{
			name:     "NullSlice",
			value:    nil,
			expected: NewSliceUpdate[int](nil),
		},
		{
			name:     "EmptySlice",
			value:    []int{},
			expected: NewSliceUpdate([]int{}),
		},
		{
			name:     "NonemptySlice",
			value:    testSlice1,
			expected: NewSliceUpdate(testSlice1),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actual := SliceUpdate[int]{}
			actual.SetValue(testCase.value)
			if !reflect.DeepEqual(actual, testCase.expected) {
				t.Errorf("Expected: %v. Actual: %v", testCase.expected, actual)
			}
		})
	}
}

func TestSliceUpdate_Removed(t *testing.T) {
	testCases := []struct {
		name     string
		u        SliceUpdate[int]
		expected bool
	}{
		{
			name:     "NotSet",
			u:        SliceUpdate[int]{},
			expected: false,
		},
		{
			name:     "NullInt",
			u:        NewSliceUpdate[int](nil),
			expected: true,
		},
		{
			name:     "ValueInt",
			u:        NewSliceUpdate(testSlice1),
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

func TestSliceUpdate_Value(t *testing.T) {
	var u SliceUpdate[int]
	if u.Value() != nil {
		t.Errorf("Expected: nil. Actual: %v", u.Value())
	}
	u.SetValue(testSlice1)
	if !reflect.DeepEqual(u.Value(), testSlice1) {
		t.Errorf("Expected: %v. Actual: %v", testSlice1, u.Value())
	}
}

func TestSliceUpdate_Apply(t *testing.T) {
	testCases := []struct {
		name     string
		u        SliceUpdate[int]
		expected []int
	}{
		{
			name:     "Unset",
			u:        SliceUpdate[int]{},
			expected: testSlice1,
		},
		{
			name:     "Removed",
			u:        NewSliceUpdate[int](nil),
			expected: nil,
		},
		{
			name:     "Set",
			u:        NewSliceUpdate([]int{testValue, testValue}),
			expected: []int{testValue, testValue},
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
			name:     "Unset",
			u:        SliceUpdate[int]{},
			value:    testSlice1,
			expected: SliceUpdate[int]{},
		},
		{
			name:     "Removed/NonZeroValue",
			u:        NewSliceUpdate[int](nil),
			value:    testSlice1,
			expected: NewSliceUpdate[int](nil),
		},
		{
			name:     "Removed/ZeroValue",
			u:        NewSliceUpdate[int](nil),
			value:    nil,
			expected: SliceUpdate[int]{},
		},
		{
			name:     "Set/Equal",
			u:        NewSliceUpdate(testSlice1),
			value:    testSlice1,
			expected: SliceUpdate[int]{},
		},
		{
			name:     "Set/NotEqual",
			u:        NewSliceUpdate(testSlice2),
			value:    testSlice1,
			expected: NewSliceUpdate(testSlice2),
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
			name:     "Unset",
			u:        SliceUpdate[int]{},
			expected: false,
		},
		{
			name:     "Removed",
			u:        NewSliceUpdate[int](nil),
			expected: false,
		},
		{
			name:     "Set/NotEqual/DifferentSizes",
			u:        NewSliceUpdate(testSlice2),
			expected: false,
		},
		{
			name:     "Set/NotEqual/SameSize",
			u:        NewSliceUpdate([]int{0}),
			expected: false,
		},
		{
			name:     "Set/Equal",
			u:        NewSliceUpdate(testSlice1),
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
			name:     "Unset",
			u:        SliceUpdate[int]{},
			expected: false,
		},
		{
			name:     "Removed",
			u:        NewSliceUpdate[int](nil),
			expected: false,
		},
		{
			name:     "Set/NotSatisfied",
			u:        NewSliceUpdate(testSlice1),
			expected: false,
		},
		{
			name:     "Set/Satisfied",
			u:        NewSliceUpdate(testSlice2),
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
			name:     "Unset",
			u:        SliceUpdate[string]{},
			expected: "<unset>",
		},
		{
			name:     "Removed",
			u:        NewSliceUpdate[string](nil),
			expected: "<removed>",
		},
		{
			name:     "Set/StringerSlice",
			u:        NewSliceUpdate([]stringer{{}}),
			expected: "[stringer]",
		},
		{
			name:     "Set/IntSlice",
			u:        NewSliceUpdate([]int{42}),
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
