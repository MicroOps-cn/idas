package sets

import (
	"fmt"
	"sort"
)

type Comparable interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64 | ~string
}

func New[T Comparable](items ...T) Set[T] {
	s := Set[T]{}
	s.Insert(items...)
	return s
}

type Set[T Comparable] map[T]int

func (s Set[T]) String() string {
	return fmt.Sprintf("%v", s.List())
}

// Insert adds items to the set.
func (s Set[T]) Insert(items ...T) Set[T] {
	for _, item := range items {
		if _, ok := s[item]; !ok {
			s[item] = s.Len()
		}
	}
	return s
}

// Delete removes all items from the set.
func (s Set[T]) Delete(items ...T) Set[T] {
	for _, item := range items {
		delete(s, item)
	}
	return s
}

// Has returns true if and only if item is contained in the set.
func (s Set[T]) Has(item T) bool {
	_, contained := s[item]
	return contained
}

// HasAll returns true if and only if all items are contained in the set.
func (s Set[T]) HasAll(items ...T) bool {
	for _, item := range items {
		if !s.Has(item) {
			return false
		}
	}
	return true
}

// HasAny returns true if any items are contained in the set.
func (s Set[T]) HasAny(items ...T) bool {
	for _, item := range items {
		if s.Has(item) {
			return true
		}
	}
	return false
}

// Difference returns a set of objects that are not in s2
// For example:
// s1 = {a1, a2, a3}
// s2 = {a1, a2, a4, a5}
// s1.Difference(s2) = {a3}
// s2.Difference(s1) = {a4, a5}
func (s Set[T]) Difference(s2 Set[T]) Set[T] {
	result := New[T]()
	for key := range s {
		if !s2.Has(key) {
			result.Insert(key)
		}
	}
	return result
}

// Union returns a new set which includes items in either s1 or s2.
// For example:
// s1 = {a1, a2}
// s2 = {a3, a4}
// s1.Union(s2) = {a1, a2, a3, a4}
// s2.Union(s1) = {a1, a2, a3, a4}
func (s Set[T]) Union(s2 Set[T]) Set[T] {
	result := New[T]()
	for key := range s {
		result.Insert(key)
	}
	for key := range s2 {
		result.Insert(key)
	}
	return result
}

// Intersection returns a new set which includes the item in BOTH s1 and s2
// For example:
// s1 = {a1, a2}
// s2 = {a2, a3}
// s1.Intersection(s2) = {a2}
func (s Set[T]) Intersection(s2 Set[T]) Set[T] {
	var walk, other Set[T]
	result := New[T]()
	if s.Len() < s2.Len() {
		walk = s
		other = s2
	} else {
		walk = s2
		other = s
	}
	for key := range walk {
		if other.Has(key) {
			result.Insert(key)
		}
	}
	return result
}

// IsSuperset returns true if and only if s1 is a superset of s2.
func (s Set[T]) IsSuperset(s2 Set[T]) bool {
	for item := range s2 {
		if !s.Has(item) {
			return false
		}
	}
	return true
}

type sortableSlice[T Comparable] []T

func (s sortableSlice[T]) Len() int           { return len(s) }
func (s sortableSlice[T]) Less(i, j int) bool { return less[T](s[i], s[j]) }
func (s sortableSlice[T]) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

// Equal returns true if and only if s1 is equal (as a set) to s2.
// Two sets are equal if their membership is identical.
// (In practice, this means same elements, order doesn't matter)
func (s Set[T]) Equal(s2 Set[T]) bool {
	return s.Len() == s2.Len() && s.IsSuperset(s2)
}

// SortedList returns the slice with contents in random order.
func (s Set[T]) SortedList() []T {
	res := make([]T, 0, s.Len())
	for key := range s {
		res = append(res, key)
	}
	sort.Sort(sortableSlice[T](res))
	return res
}

// List returns the contents as a sorted string slice.
func (s Set[T]) List() []T {
	res := make(sortableSet[T], 0, s.Len())

	for item, index := range s {
		res = append([]setItem[T](res), setItem[T]{index: index, item: item})
	}
	sort.Sort(res)
	return res.List()
}

// PopAny returns a single element from the set.
func (s Set[T]) PopAny() (T, bool) {
	for key := range s {
		s.Delete(key)
		return key, true
	}
	var zeroValue T
	return zeroValue, false
}

// Len returns the size of the set.
func (s Set[T]) Len() int {
	return len(s)
}

func less[T Comparable](lhs, rhs T) bool {
	return lhs < rhs
}

type setItem[T Comparable] struct {
	index int
	item  T
}

type sortableSet[T Comparable] []setItem[T]

func (s sortableSet[T]) Len() int {
	return len(s)
}

func (s sortableSet[T]) Less(i, j int) bool {
	return less[int](s[i].index, s[j].index)
}

func (s sortableSet[T]) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s sortableSet[T]) List() []T {
	res := make([]T, s.Len())
	for idx, item := range s {
		res[idx] = item.item
	}
	return res
}
