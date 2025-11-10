package utils

import (
	"testing"
)

func TestSet_ToArray(t *testing.T) {
	var set1 Set[int] = map[int]bool{}
	var set2 Set[int] = map[int]bool{1: true, 2: true, 4: true}

	expected1 := []int{}
	expected2 := []int{1, 2, 4}

	CompareSlices(t, expected1, set1.ToArray())
	CompareSlices(t, expected2, set2.ToArray())
}
