package infra

type Set[T comparable] map[T]struct{}

func NewSet[T comparable](ts ...T) Set[T] {
	s := make(Set[T], IfThenElse(len(ts) == 0, 8, len(ts)))
	for _, v := range ts {
		s[v] = struct{}{}
	}
	return s
}

func (s Set[T]) AddIfNotExist(t T) (exist bool) {
	_, exist = s[t]
	if !exist {
		s[t] = struct{}{}
	}
	return
}

func (s Set[T]) Add(t T) {
	s[t] = struct{}{}
}

func (s Set[T]) Contains(t T) bool {
	_, ok := s[t]
	return ok
}

func (s Set[T]) ContainsAll(ts []T) bool {
	for _, t := range ts {
		if _, ok := s[t]; !ok {
			return false
		}
	}
	return true
}

func (s Set[T]) Members() []T {
	members := make([]T, 0, len(s))
	for m := range s {
		members = append(members, m)
	}
	return members
}

func (s Set[T]) Equal(target Set[T]) bool {
	if len(s) != len(target) {
		return false
	}
	for k := range s {
		if _, ok := target[k]; !ok {
			return false
		}
	}
	return true
}
