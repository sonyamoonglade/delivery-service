package tgdelivery

import (
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestSixifyOrderId(t *testing.T) {

	testingTable := []TestId{
		{
			id:       1,
			expected: "000001",
		},
		{
			id:       2712,
			expected: "002712",
		},
		{
			id:       111111,
			expected: "111111",
		},
	}

	for _, testId := range testingTable {
		actual := SixifyOrderId(testId.id)
		expected := testId.expected
		assert.Equal(t, actual, expected)
	}

}

type TestId struct {
	id       int64
	expected string
}
