package infra

func MapSortV[K comparable, V any](m map[K]V, fn func(left, right V) bool) []V {
	vs := make([]V, 0, len(m))
	for _, v := range m {
		vs = append(vs, v)
	}
	SliceSort(vs, fn)
	return vs
}
