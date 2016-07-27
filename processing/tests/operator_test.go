package processing

import (
	"testing"
	"reflect"
	"splash/processing"
)

type TestCase struct {
	data []int
	expected[]int
}

func TestProcessing(t *testing.T) {
	ShouldProcessCorrectly(t)
}

func ShouldProcessCorrectly(t *testing.T) {

	tests := map[string]*TestCase{
		"Should return an empty array.": &TestCase{
			data: []int{},
			expected: []int{},
		},
		"Should return an array without duplicats.": &TestCase{
			data: []int{1, 2, 1, 2, 3},
			expected: []int{1, 2, 3},
		},
		"Should return a sorted array": &TestCase{
			data: []int{5, 1, 2, 3},
			expected: []int{1, 2, 3, 5},
		},
		"Should return a sorted array without duplicates": &TestCase{
			data: []int{5, 1, 2, 3, 5 , 1, 7, 1},
			expected: []int{1, 2, 3, 5, 7},
		},
	}

	operator := processing.NewOperator()

	for caseName, testCase := range tests {

		got := operator.Operate(testCase.data)

		if(!reflect.DeepEqual(testCase.expected, got)) {
			t.Error(
				"For", caseName,
				"expected", testCase.expected,
				"got", got,
			)
		}
	}
}