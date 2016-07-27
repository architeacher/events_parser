package processing

import "sort"

type Operator struct {
}

func NewOperator() *Operator {
	return &Operator{}
}


func (self *Operator) Operate(data map[string]int) map[string]int {

	if 0 == len(data) {
		return data
	}

	//data = self.removeDuplicates(data)
	//data = self.sort(data)
	return data
}

func (*Operator) removeDuplicates(data []int) []int {

	encountered := map[int]bool{}
	result := []int{}

	for _, value := range data {
		if !encountered[value]{
			encountered[value] = true
			result = append(result, value)
		}
	}

	return result
}

func (*Operator) sort(data []int) []int {
	sort.Ints(data)
	return data
}
