package nup

import "testing"

func TestOperation_String(t *testing.T) {
	testCases := []struct {
		name     string
		op       Operation
		expected string
	}{
		{
			name:     "Noop",
			op:       OpNoop,
			expected: "no-op",
		},
		{
			name:     "Remove",
			op:       OpRemove,
			expected: "remove",
		},
		{
			name:     "Set",
			op:       OpSet,
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
