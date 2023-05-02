package models

import (
	"errors"
	"fmt"
)

type Toilet struct {
	count       int
	ArrayUrinal []int
}

func NewToilet() *Toilet {
	return &Toilet{
		count:       0,
		ArrayUrinal: make([]int, 0, 0),
	}
}

func (t *Toilet) AddUrinal(count int) {
	t.count = count
	t.ArrayUrinal = make([]int, count, count)
}

func (t *Toilet) TakePlace(index int) error {
	if t.ArrayUrinal[index] == t.count {
		return errors.New(fmt.Sprintf("This index '%d' busy", index))
	}
	t.ArrayUrinal[index] = t.count
	for i := range t.ArrayUrinal {
		if i < index {
			newValue := t.count - (index - i)
			if newValue > t.ArrayUrinal[i] {
				t.ArrayUrinal[i] = newValue
			}
		} else {
			newValue := t.count + (index - i)
			if newValue > t.ArrayUrinal[i] {
				t.ArrayUrinal[i] = newValue
			}
		}
	}
	return nil
}

func (t *Toilet) GetOptimalPlace() (int, error) {
	minIndexMap := make(map[int][]int, t.count)
	minWt := t.count
	for i, v := range t.ArrayUrinal {
		if minWt >= v {
			minWt = v
			_, ok := minIndexMap[minWt]
			if ok {
				minIndexMap[minWt] = append(minIndexMap[minWt], i)
			} else {
				minIndexMap[minWt] = make([]int, 0, t.count)
				minIndexMap[minWt] = append(minIndexMap[minWt], i)
			}
		}
	}
	if minWt == t.count {
		return 0, errors.New("Doesn't have optimal place")
	}
	if len(minIndexMap[minWt]) == 1 {
		return minIndexMap[minWt][0], nil
	}
	minIndexValue := 0
	indexResult := 0
	for _, v := range minIndexMap[minWt] {
		sum := t.getModuleIndexSum(v)
		if minIndexValue < sum {
			minIndexValue = sum
			indexResult = v
		}
	}
	return indexResult, nil
}

func (t *Toilet) getModuleIndexSum(index int) int {
	result := [2]int{0, 0}
	for i, v := range t.ArrayUrinal {
		if i < index {
			result[0] += v
		}
		if i > index {
			result[1] += v
		}
	}
	a := result[0] - result[1]
	if a < 0 {
		return -a
	}
	return a
}
