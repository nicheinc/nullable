package nullable

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

// Ensure implementation of Nullable interface.
var _ Nullable = &Bool{}

func TestBool_UnmarshalJSON(t *testing.T) {
	testCases := []struct {
		name     string
		json     string
		expected Bool
	}{
		{
			name:     "EmptyJSONObject",
			json:     `{}`,
			expected: Bool{},
		},
		{
			name:     "NullBool",
			json:     `{"int": null}`,
			expected: NewBoolPtr(nil),
		},
		{
			name:     "ValueBool",
			json:     fmt.Sprintf(`{"int": %v}`, true),
			expected: NewBool(true),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			var dst struct {
				Bool Bool `json:"int"`
			}
			if err := json.Unmarshal([]byte(testCase.json), &dst); err != nil {
				t.Errorf("Error unmarshaling JSON: %s", err)
			}
			if !reflect.DeepEqual(dst.Bool, testCase.expected) {
				t.Errorf("Expected: %v. Actual: %v", testCase.expected, dst.Bool)
			}
		})
	}
}

func TestBool_SetValue(t *testing.T) {
	testCases := []struct {
		name     string
		value    *bool
		expected Bool
	}{
		{
			name:     "NullBool",
			value:    nil,
			expected: NewBoolPtr(nil),
		},
		{
			name:     "ValueBool",
			value:    func(v bool) *bool { return &v }(true),
			expected: NewBool(true),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actual := Bool{}
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

func TestBool_Removed(t *testing.T) {
	testCases := []struct {
		name     string
		b        Bool
		expected bool
	}{
		{
			name:     "NotSet",
			b:        Bool{},
			expected: false,
		},
		{
			name:     "NullBool",
			b:        NewBoolPtr(nil),
			expected: true,
		},
		{
			name:     "ValueBool",
			b:        NewBool(true),
			expected: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actual := testCase.b.Removed()
			if !reflect.DeepEqual(actual, testCase.expected) {
				t.Errorf("Expected: %v, Actual: %v", testCase.expected, actual)
			}
		})
	}
}

func TestBool_Value(t *testing.T) {
	var b Bool
	if b.Value() != nil {
		t.Errorf("Expected: nil, Actual: %v", b.Value())
	}
	expected := true
	b.SetValue(expected)
	if *b.Value() != expected {
		t.Errorf("Expected: %v, Actual: %v", expected, *b.Value())
	}
}

func TestBool_InterfaceValue(t *testing.T) {
	var b Bool
	if !reflect.ValueOf(b.InterfaceValue()).IsNil() {
		t.Errorf("Expected: nil, Actual: %v", b.InterfaceValue())
	}
	expected := true
	b.SetValue(expected)
	if !reflect.DeepEqual(b.InterfaceValue(), &expected) {
		t.Errorf("Expected: %v, Actual: %v", expected, b.InterfaceValue())
	}
}

func TestBool_Scan(t *testing.T) {
	testCases := []struct {
		name     string
		dest     Bool
		src      interface{}
		expected Bool
		err      error
	}{
		{
			name: "ScanTypeError",
			dest: Bool{},
			src:  0,
			err: &ScanTypeError{
				Src:  0,
				Dest: &Bool{},
			},
		},
		{
			name:     "Nil",
			dest:     NewBool(true),
			src:      nil,
			expected: Bool{},
		},
		{
			name:     "Bool",
			dest:     Bool{},
			src:      true,
			expected: NewBool(true),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			err := testCase.dest.Scan(testCase.src)
			if !reflect.DeepEqual(err, testCase.err) {
				t.Errorf("Expected error: %v, Actual: %v", testCase.err, err)
			}
			if !reflect.DeepEqual(testCase.dest, testCase.expected) {
				t.Errorf("Expected: %v, Actual: %v", testCase.expected, testCase.dest)
			}
		})
	}
}
