package utils

import (
	"slices"
	"testing"
)

func TestSet_ToArray(t *testing.T) {
	var set1 Set[int] = map[int]bool{}
	var set2 Set[int] = map[int]bool{1: true, 2: true, 4: true}

	expected1 := []int{}
	expected2 := []int{1, 2, 4}

	val1 := set1.ToArray()
	val2 := set2.ToArray()
	slices.Sort(val1)
	slices.Sort(val2)

	CompareSlices(t, expected1, val1)
	CompareSlices(t, expected2, val2)
}

func TestUnique(t *testing.T) {
	array1 := []string{"a", "b", "a", "c"}
	array2 := []string{"c", "a", "b", "a"}

	var set Set[string] = map[string]bool{"a": true, "b": true, "c": true}

	gen1 := Unique(array1)
	gen2 := Unique(array2)

	if !gen2.Equals(gen1) {
		t.Errorf("Should match: %v == %v", gen1.ToArray(), gen2.ToArray())
	}
	if !gen1.Equals(gen2) {
		t.Errorf("Should match: %v == %v", gen2.ToArray(), gen1.ToArray())
	}

	if !set.Equals(gen1) {
		t.Errorf("Should match: %v == %v", set.ToArray(), gen1.ToArray())
	}
	if !set.Equals(gen2) {
		t.Errorf("Shouldrmatch: %v == %v", set.ToArray(), gen2.ToArray())
	}
}

func TestSet_Equals(t *testing.T) {
	var set1 Set[string] = map[string]bool{"a": true, "b": true, "c": true}
	var set2 Set[string] = map[string]bool{"a": true, "b": true}
	var set3 Set[string] = map[string]bool{"c": true, "a": true, "b": true}

	if !set1.Equals(set3) {
		t.Errorf("These should be equal: %v == %v", set1, set3)
	}
	if !set3.Equals(set1) {
		t.Errorf("These should be equal: %v == %v", set3, set1)
	}

	if set1.Equals(set2) {
		t.Errorf("These should not be equal: %v == %v", set1, set2)
	}
	if set2.Equals(set1) {
		t.Errorf("These should not be equal: %v == %v", set2, set1)
	}

	if set3.Equals(set2) {
		t.Errorf("These should not be equal: %v == %v", set3, set2)
	}
	if set2.Equals(set3) {
		t.Errorf("These should not be equal: %v == %v", set2, set3)
	}
}

func TestSet_Union(t *testing.T) {
}

func TestSet_Contains(t *testing.T) {
}
