package utils

import (
	"testing"
)

func CompareSlices[P comparable](t *testing.T, expected []P, actual []P) {
	if len(expected) != len(actual) {
		t.Errorf("Lenghts mismatch: %d != %d\n  --> %v != %v", len(expected), len(actual), expected, actual)
		return
	}

	for index, exp := range expected {
		if exp != actual[index] {
			t.Errorf("Failure at index %d: %v != %v", index, exp, actual[index])
		}
	}
}
