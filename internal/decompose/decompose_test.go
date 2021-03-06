package decompose

import (
	"reflect"
	"testing"

	"gotest.tools/assert"
)

type DecomposeTestCase struct {
	original int
	expected []string
	length   int
	base     int
}

func TestDecompose(t *testing.T) {
	testCases := []DecomposeTestCase{
		{
			original: 123456789,
			expected: []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"},
			base:     10,
			length:   9,
		},
		{
			original: 4321,
			expected: []string{"1", "0", "0", "0", "0", "1", "1", "1", "0", "0", "0", "0", "1"},
			base:     2,
			length:   13,
		},
		{
			original: 0,
			expected: []string{"0", "0", "0", "0"},
			base:     8,
			length:   4,
		},
		{
			original: 2,
			expected: []string{"0", "2"},
			base:     10,
			length:   2,
		},
		{
			original: 1,
			expected: []string{"1"},
			base:     2,
			length:   1,
		},
	}

	for _, testCase := range testCases {
		actual := Value(testCase.original, testCase.base, testCase.length)
		assert.Assert(t, reflect.DeepEqual(testCase.expected, actual))
	}
}
