package nullable

import (
	"encoding/json"
	"reflect"
	"testing"
)

// Ensure implementation of Nullable interface.
var _ Nullable = &String{}

func TestString_UnmarshalJSON(t *testing.T) {
	testCases := []struct {
		name     string
		json     string
		expected String
	}{
		{
			name:     "EmptyJSONObject",
			json:     `{}`,
			expected: String{},
		},
		{
			name:     "NullString",
			json:     `{"string": null}`,
			expected: NewStringPtr(nil),
		},
		{
			name:     "EmptyString",
			json:     `{"string": ""}`,
			expected: NewString(""),
		},
		{
			name:     "SpaceString",
			json:     `{"string": " "}`,
			expected: NewString(" "),
		},
		{
			name:     "ValueString",
			json:     `{"string": "value"}`,
			expected: NewString("value"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			var dst struct {
				String String `json:"string"`
			}
			if err := json.Unmarshal([]byte(testCase.json), &dst); err != nil {
				t.Errorf("Error unmarshaling JSON: %s", err)
			}
			if !reflect.DeepEqual(dst.String, testCase.expected) {
				t.Errorf("Expected: %v. Actual: %v", testCase.expected, dst.String)
			}
		})
	}
}

func TestString_SetValue(t *testing.T) {
	testCases := []struct {
		name     string
		value    *string
		expected String
	}{
		{
			name:     "NullString",
			value:    nil,
			expected: NewStringPtr(nil),
		},
		{
			name:     "EmptyString",
			value:    func(v string) *string { return &v }(""),
			expected: NewString(""),
		},
		{
			name:     "SpaceString",
			value:    func(v string) *string { return &v }(" "),
			expected: NewString(" "),
		},
		{
			name:     "ValueString",
			value:    func(v string) *string { return &v }("value"),
			expected: NewString("value"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actual := String{}
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

func TestString_Removed(t *testing.T) {
	testCases := []struct {
		name     string
		str      String
		expected bool
	}{
		{
			name:     "NotSet",
			str:      String{},
			expected: false,
		},
		{
			name:     "NullString",
			str:      NewStringPtr(nil),
			expected: true,
		},
		{
			name:     "EmptyString",
			str:      NewString(""),
			expected: false,
		},
		{
			name:     "SpaceString",
			str:      NewString(" "),
			expected: false,
		},
		{
			name:     "ValueString",
			str:      NewString("value"),
			expected: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actual := testCase.str.Removed()
			if !reflect.DeepEqual(actual, testCase.expected) {
				t.Errorf("Expected: %v, Actual: %v", testCase.expected, actual)
			}
		})
	}
}

func TestString_IsEmpty(t *testing.T) {
	testCases := []struct {
		name     string
		str      String
		expected bool
	}{
		{
			name:     "NotSet",
			str:      String{},
			expected: false,
		},
		{
			name:     "NullString",
			str:      NewStringPtr(nil),
			expected: false,
		},
		{
			name:     "EmptyString",
			str:      NewString(""),
			expected: true,
		},
		{
			name:     "SpaceString",
			str:      NewString(" "),
			expected: true,
		},
		{
			name:     "TabString",
			str:      NewString("\t"),
			expected: true,
		},
		{
			name:     "NewlineString",
			str:      NewString("\n"),
			expected: true,
		},
		{
			name:     "ValueString",
			str:      NewString("value"),
			expected: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actual := testCase.str.IsEmpty()
			if !reflect.DeepEqual(actual, testCase.expected) {
				t.Errorf("Expected: %v, Actual: %v", testCase.expected, actual)
			}
		})
	}
}

func TestString_Value(t *testing.T) {
	var s String
	if s.Value() != nil {
		t.Errorf("Expected: nil, Actual: %v", s.Value())
	}
	expected := "value"
	s.SetValue(expected)
	if *s.Value() != expected {
		t.Errorf("Expected: %v, Actual: %v", expected, *s.Value())
	}
}

func TestString_InterfaceValue(t *testing.T) {
	var s String
	if !reflect.ValueOf(s.InterfaceValue()).IsNil() {
		t.Errorf("Expected: nil, Actual: %v", s.InterfaceValue())
	}
	expected := "value"
	s.SetValue(expected)
	if !reflect.DeepEqual(s.InterfaceValue(), &expected) {
		t.Errorf("Expected: %v, Actual: %v", expected, s.InterfaceValue())
	}
}

func TestString_String(t *testing.T) {
	var s String
	if s.String() != "<unset>" {
		t.Errorf("Expected: <unset>, Actual: %v", s.String())
	}
	s.SetPtr(nil)
	if s.String() != "<removed>" {
		t.Errorf("Expected: <removed>, Actual: %v", s.String())
	}
	expected := "value"
	s.SetValue(expected)
	if s.String() != expected {
		t.Errorf("Expected: %v, Actual: %v", expected, s.String())
	}
}

func TestString_Scan(t *testing.T) {
	testCases := []struct {
		name     string
		dest     String
		src      interface{}
		expected String
		err      error
	}{
		{
			name: "ScanTypeError",
			dest: String{},
			src:  true,
			err: &ScanTypeError{
				Src:  true,
				Dest: &String{},
			},
		},
		{
			name:     "Nil",
			dest:     String{},
			src:      nil,
			expected: NewStringPtr(nil),
		},
		{
			name:     "String",
			dest:     String{},
			src:      "value",
			expected: NewString("value"),
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
