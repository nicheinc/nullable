package nup

import (
	"testing"

	"github.com/nicheinc/expect"
)

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
			expect.Equal(t, actual, testCase.expected)
		})
	}
}
