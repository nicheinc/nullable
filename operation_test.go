package nullable

import "testing"

func TestOperation_String(t *testing.T) {
	testCases := []struct {
		name     string
		op       Operation
		expected string
	}{
		{
			name:     "Noop",
			op:       Noop,
			expected: "no-op",
		},
		{
			name:     "Remove",
			op:       Remove,
			expected: "remove",
		},
		{
			name:     "Set",
			op:       Set,
			expected: "set",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actual := testCase.op.String()
			if actual != testCase.expected {
				t.Errorf("Expected: %v. Actual: %v", testCase.expected, actual)
			}
		})
	}
}
