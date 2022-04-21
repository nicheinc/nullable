package nullable

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"
)

// Ensure implementation of Nullable interface.
var _ Nullable = &Time{}

func TestTime_UnmarshalJSON(t *testing.T) {
	testCases := []struct {
		name     string
		json     string
		expected Time
	}{
		{
			name:     "EmptyJSONObject",
			json:     `{}`,
			expected: Time{},
		},
		{
			name:     "NullTime",
			json:     `{"time": null}`,
			expected: NewTimePtr(nil),
		},
		{
			name:     "ValueTime",
			json:     `{"time": "2022-04-21T09:57:01Z"}`,
			expected: NewTime(time.Date(2022, 4, 21, 9, 57, 1, 0, time.UTC)),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			var dst struct {
				Time Time `json:"time"`
			}
			if err := json.Unmarshal([]byte(testCase.json), &dst); err != nil {
				t.Errorf("Error unmarshaling JSON: %s", err)
			}
			if !reflect.DeepEqual(dst.Time, testCase.expected) {
				t.Errorf("Expected: %v. Actual: %v", testCase.expected, dst.Time)
			}
		})
	}
}

func TestTime_SetValue(t *testing.T) {
	var testTime = time.Now()
	testCases := []struct {
		name     string
		value    *time.Time
		expected Time
	}{
		{
			name:     "NullTime",
			value:    nil,
			expected: NewTimePtr(nil),
		},
		{
			name:     "ValueTime",
			value:    &testTime,
			expected: NewTime(testTime),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actual := Time{}
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

func TestTime_Removed(t *testing.T) {
	testCases := []struct {
		name     string
		t        Time
		expected bool
	}{
		{
			name:     "NotSet",
			t:        Time{},
			expected: false,
		},
		{
			name:     "NullTime",
			t:        NewTimePtr(nil),
			expected: true,
		},
		{
			name:     "ValueTime",
			t:        NewTime(time.Now()),
			expected: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actual := testCase.t.Removed()
			if !reflect.DeepEqual(actual, testCase.expected) {
				t.Errorf("Expected: %v, Actual: %v", testCase.expected, actual)
			}
		})
	}
}

func TestTime_IsZero(t *testing.T) {
	testCases := []struct {
		name     string
		t        Time
		expected bool
	}{
		{
			name:     "NotSet",
			t:        Time{},
			expected: false,
		},
		{
			name:     "NullTime",
			t:        NewTimePtr(nil),
			expected: false,
		},
		{
			name:     "ZeroTime",
			t:        NewTime(time.Time{}),
			expected: true,
		},
		{
			name:     "NonZeroTime",
			t:        NewTime(time.Now()),
			expected: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actual := testCase.t.IsZero()
			if !reflect.DeepEqual(actual, testCase.expected) {
				t.Errorf("Expected: %v, Actual: %v", testCase.expected, actual)
			}
		})
	}
}

func TestTime_Value(t *testing.T) {
	var testTime Time
	if testTime.Value() != nil {
		t.Errorf("Expected: nil, Actual: %v", testTime.Value())
	}
	expected := time.Now()
	testTime.SetValue(expected)
	if *testTime.Value() != expected {
		t.Errorf("Expected: %v, Actual: %v", expected, *testTime.Value())
	}
}

func TestTime_Equals(t *testing.T) {
	var testTime = time.Date(2022, 04, 21, 10, 37, 51, 0, time.UTC)
	testCases := []struct {
		name     string
		t        Time
		expected bool
	}{
		{
			name:     "Unset",
			t:        Time{},
			expected: false,
		},
		{
			name:     "Removed",
			t:        NewTimePtr(nil),
			expected: false,
		},
		{
			name:     "Set/NotEqualValue",
			t:        NewTime(time.Now()),
			expected: false,
		},
		{
			name:     "Set/Equal",
			t:        NewTime(testTime),
			expected: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if actual := testCase.t.Equals(testTime); actual != testCase.expected {
				t.Errorf("Expected: %v. Actual: %v", testCase.expected, actual)
			}
		})
	}
}

func TestTime_Apply(t *testing.T) {
	var (
		value1 = time.Date(2022, 04, 21, 10, 37, 51, 0, time.UTC)
		value2 = time.Now()
	)
	testCases := []struct {
		name     string
		t        Time
		expected time.Time
	}{
		{
			name:     "Unset",
			t:        Time{},
			expected: value1,
		},
		{
			name:     "Removed",
			t:        NewTimePtr(nil),
			expected: time.Time{},
		},
		{
			name:     "Set",
			t:        NewTime(value2),
			expected: value2,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if actual := testCase.t.Apply(value1); actual != testCase.expected {
				t.Errorf("Expected: %v. Actual: %v", testCase.expected, actual)
			}
		})
	}
}

func TestTime_ApplyPtr(t *testing.T) {
	var (
		value1 = time.Date(2022, 04, 21, 10, 37, 51, 0, time.UTC)
		value2 = time.Now()
	)
	testCases := []struct {
		name     string
		t        Time
		expected *time.Time
	}{
		{
			name:     "Unset",
			t:        Time{},
			expected: &value1,
		},
		{
			name:     "Removed",
			t:        NewTimePtr(nil),
			expected: nil,
		},
		{
			name:     "Set",
			t:        NewTime(value2),
			expected: &value2,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if actual := testCase.t.ApplyPtr(&value1); !reflect.DeepEqual(actual, testCase.expected) {
				t.Errorf("Expected: %v. Actual: %v", testCase.expected, actual)
			}
		})
	}
}

func TestTime_Diff(t *testing.T) {
	var (
		value1 = time.Date(2022, 04, 21, 10, 37, 51, 0, time.UTC)
		value2 = time.Now()
	)
	testCases := []struct {
		name     string
		t        Time
		value    time.Time
		expected Time
	}{
		{
			name:     "Unset",
			t:        Time{},
			value:    value1,
			expected: Time{},
		},
		{
			name:     "Removed/NonZeroValue",
			t:        NewTimePtr(nil),
			value:    value1,
			expected: NewTimePtr(nil),
		},
		{
			name:     "Removed/ZeroValue",
			t:        NewTimePtr(nil),
			value:    time.Time{},
			expected: Time{},
		},
		{
			name:     "Set/Equal",
			t:        NewTime(value1),
			value:    value1,
			expected: Time{},
		},
		{
			name:     "Set/NotEqual",
			t:        NewTime(value2),
			value:    value1,
			expected: NewTime(value2),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if actual := testCase.t.Diff(testCase.value); !reflect.DeepEqual(actual, testCase.expected) {
				t.Errorf("Expected: %v. Actual: %v", testCase.expected, actual)
			}
		})
	}
}

func TestTime_InterfaceValue(t *testing.T) {
	var testTime Time
	if !reflect.ValueOf(testTime.InterfaceValue()).IsNil() {
		t.Errorf("Expected: nil, Actual: %v", testTime.InterfaceValue())
	}
	expected := time.Now()
	testTime.SetValue(expected)
	if !reflect.DeepEqual(testTime.InterfaceValue(), &expected) {
		t.Errorf("Expected: %v, Actual: %v", expected, testTime.InterfaceValue())
	}
}
