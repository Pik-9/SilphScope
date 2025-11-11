package utils

import "fmt"

type Set[T comparable] map[T]bool

func (set Set[T]) ToArray() []T {
	ret := make([]T, 0, len(set))
	for k := range set {
		ret = append(ret, k)
	}
	return ret
}

func Unique[T comparable](array []T) Set[T] {
	ret := make(Set[T])
	for _, val := range array {
		ret[val] = true
	}
	return ret
}

func (set *Set[T]) Add(element T) {
	(*set)[element] = true
}

func (set Set[T]) Union(other Set[T]) Set[T] {
	ret := set
	for key := range other {
		ret[key] = true
	}

	return ret
}

func (set Set[T]) Contains(value T) bool {
	_, ret := set[value]
	return ret
}

func (set Set[T]) Equals(other Set[T]) bool {
	if len(set) != len(other) {
		return false
	}

	for val := range set {
		_, contains := other[val]
		if !contains {
			return false
		}
	}

	return true
}

func (set Set[T]) String() string {
	inner := "{"
	index := 0
	for key := range set {
		if index == 0 {
			inner = fmt.Sprintf("%s%v", inner, key)
		} else {
			inner = fmt.Sprintf("%s, %v", inner, key)
		}
		index += 1
	}
	return inner + "}"
}
