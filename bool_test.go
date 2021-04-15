package nullable

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/nicheinc/go-common/v12/test"
)

func TestBool_UnmarshalJSON(t *testing.T) {
	testCases := []struct {
		name     string
		json     string
		expected Bool
	}{
		{
			name: "EmptyJSONObject",
			json: `{}`,
			expected: Bool{
				Set:   false,
				Value: nil,
			},
		},
		{
			name: "NullBool",
			json: `{"int": null}`,
			expected: Bool{
				Set:   true,
				Value: nil,
			},
		},
		{
			name: "ValueBool",
			json: fmt.Sprintf(`{"int": %v}`, true),
			expected: Bool{
				Set:   true,
				Value: test.BoolToPtr(true),
			},
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
			name:  "NullBool",
			value: nil,
			expected: Bool{
				Set:   true,
				Value: nil,
			},
		},
		{
			name:  "ValueBool",
			value: test.BoolToPtr(true),
			expected: Bool{
				Set:   true,
				Value: test.BoolToPtr(true),
			},
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
			name: "NotSet",
			b: Bool{
				Set:   false,
				Value: nil,
			},
			expected: false,
		},
		{
			name: "NullBool",
			b: Bool{
				Set:   true,
				Value: nil,
			},
			expected: true,
		},
		{
			name: "ValueBool",
			b: Bool{
				Set:   true,
				Value: test.BoolToPtr(true),
			},
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
