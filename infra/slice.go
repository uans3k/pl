package infra

import "sort"

func SliceExist[T comparable](slice []T, target T) bool {
	for _, source := range slice {
		if source == target {
			return true
		}
	}
	return false
}

type sortable[T any] struct {
	slice []T
	less  func(left, right T) bool
}

func (s *sortable[T]) Len() int {
	return len(s.slice)
}

func (s *sortable[T]) Less(i, j int) bool {
	return s.less(s.slice[i], s.slice[j])
}

func (s *sortable[T]) Swap(i, j int) {
	tmp := s.slice[i]
	s.slice[i] = s.slice[j]
	s.slice[j] = tmp
}

func SliceSort[T any](slice []T, less func(left, right T) bool) {
	sort.Sort(&sortable[T]{
		slice: slice,
		less:  less,
	})
	return
}
