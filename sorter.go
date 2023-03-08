package pago

type Comparator[T any] interface {
	Less(i, j T) bool
}

type Sorter[T any] struct {
	pago       *Pago[T]
	comparator Comparator[T]
	sorted     []int
}

func (s *Sorter[T]) Len() int {
	return len(s.sorted)
}

func (s *Sorter[T]) Less(i, j int) bool {
	return s.comparator.Less(s.pago.keyed[s.sorted[i]], s.pago.keyed[s.sorted[j]])
}

func (s *Sorter[T]) Swap(i, j int) {
	s.sorted[i], s.sorted[j] = s.sorted[j], s.sorted[i]
}

func NewSorter[T any](comparator Comparator[T]) *Sorter[T] {
	return &Sorter[T]{comparator: comparator, sorted: make([]int, 0)}
}
