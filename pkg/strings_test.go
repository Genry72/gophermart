package pkg

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestToSliceInterface(t *testing.T) {

	input := []int{1, 2, 3, 4, 5}
	expected := []interface{}{1, 2, 3, 4, 5}
	result := ToSliceInterface(input)
	assert.Equal(t, expected, result)

	input2 := []string{"1", "2", "3", "4", "5"}
	expected = []interface{}{"1", "2", "3", "4", "5"}
	result = ToSliceInterface(input2)
	assert.Equal(t, expected, result)

	// Positive test case with empty input
	input = []int{}
	expected = []interface{}{}
	result = ToSliceInterface(input)
	assert.Equal(t, expected, result)

}
