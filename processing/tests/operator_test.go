package processing

import (
	"reflect"
	"splash/processing"
	testingUtil "splash/testing"
	"testing"
)

func TestProcessing(t *testing.T) {
	ShouldProcessCorrectly(t)
}

func ShouldProcessCorrectly(t *testing.T) {

	testCases := []*testingUtil.TestCase{
	//&testingUtil.TestCase{
	//	Id: "Should return an empty array.",
	//	Input: []int{},
	//	Expected: []int{},
	//},
	//&testingUtil.TestCase{
	//	Id: "Should return a sorted array.",
	//	Input: []int{5, 1, 2, 3},
	//	Expected: []int{1, 2, 3, 5},
	//},
	//&testingUtil.TestCase{
	//	Id: "Should return a sorted array without duplicates.",
	//	Input: []int{5, 1, 2, 3, 5 , 1, 7, 1},
	//	Expected: []int{1, 2, 3, 5, 7},
	//},
	}

	operator := processing.NewOperator()

	for _, testCase := range testCases {

		got := operator.Operate(testCase.Input.(map[string]int))

		if !reflect.DeepEqual(testCase.Expected, got) {
			t.Error(
				"For", testCase.Id,
				"expected", testCase.Expected,
				"got", got,
			)
		}
	}
}
