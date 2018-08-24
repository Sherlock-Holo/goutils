package set

import "sort"

type Comparable interface {
	Less(x interface{}) bool
}

type innerSet []Comparable

func (s innerSet) Len() int {
	return len(s)
}

func (s innerSet) Less(i, j int) bool {
	return s[i].Less(s[j])
}

func (s innerSet) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

type SortedSet struct {
	inner innerSet
}

func NewSortedSet() *SortedSet {
	return &SortedSet{
		make(innerSet, 0),
	}
}

func (s *SortedSet) Search(x Comparable) (int, bool) {
	index := sort.Search(s.inner.Len(), func(i int) bool {
		return !s.inner[i].Less(x)
	})

	if index < s.inner.Len() && s.inner[index] == x {
		return index, true
	}

	return 0, false
}

func (s *SortedSet) Add(x Comparable) bool {
	if _, ok := s.Search(x); ok {
		return false
	}

	s.inner = append(s.inner, x)
	sort.Sort(s.inner)

	return true
}

func (s *SortedSet) Delete(x Comparable) bool {
	if i, ok := s.Search(x); ok {
		s.inner = append(s.inner[:i], s.inner[i+1:]...)
		return true
	}

	return false
}

func (s *SortedSet) List() []Comparable {
	return s.inner
}

func (s *SortedSet) Get(index int) (Comparable, bool) {
	if index >= s.inner.Len() {
		return nil, false
	}

	return s.inner[index], true
}
