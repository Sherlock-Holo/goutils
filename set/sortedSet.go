package set

import "sort"

type Comparable interface {
	Less(x interface{}) bool
}

type innerSet []Comparable

/*func (s innerSet) Len() int {
	return len(s)
}

func (s innerSet) Less(i, j int) bool {
	return s[i].Less(s[j])
}

func (s innerSet) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}*/

type SortedSet struct {
	inner innerSet
}

func NewSortedSet() *SortedSet {
	return &SortedSet{
		make(innerSet, 0),
	}
}

func (s *SortedSet) Search(x Comparable) (int, bool) {
	index := sort.Search(len(s.inner), func(i int) bool {
		return !s.inner[i].Less(x)
	})

	if index < len(s.inner) && s.inner[index] == x {
		return index, true
	}

	return index, false
}

func (s *SortedSet) Add(x Comparable) bool {
	/*if index, ok := s.Search(x); ok {
		return false
	}

	s.inner = append(s.inner, x)
	// sort.Sort(s.inner)
	sort.Slice(s.inner, func(i, j int) bool {
		return s.inner[i].Less(s.inner[j])
	})

	return true*/
	index, ok := s.Search(x)
	if ok {
		return false
	}

	tmp := s.inner[index-1:]

	s.inner = append(s.inner, x)
	s.inner = append(s.inner, tmp...)

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
	if index >= len(s.inner) {
		return nil, false
	}

	return s.inner[index], true
}
