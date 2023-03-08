package gopagination

import (
	"fmt"
	"sort"
)

type Comparator[T any] interface {
	Less(i, j T) bool
}

type Sorter[T any] struct {
	paginator  *Paginator[T]
	comparator Comparator[T]
	sorted     []int
}

func (s *Sorter[T]) Len() int {
	// s.paginator.keyed好一些还是s.sorted好一些？
	return len(s.paginator.keyed)
}

func (s *Sorter[T]) Less(i, j int) bool {
	return s.comparator.Less(s.paginator.keyed[s.sorted[i]], s.paginator.keyed[s.sorted[j]])
}

func (s *Sorter[T]) Swap(i, j int) {
	s.sorted[i], s.sorted[j] = s.sorted[j], s.sorted[i]
}

func NewSorter[T any](comparator Comparator[T]) *Sorter[T] {
	return &Sorter[T]{comparator: comparator, sorted: make([]int, 0)}
}

type Paginator[T any] struct {
	keyed   map[int]T
	sorters map[string]*Sorter[T]
}

func (p *Paginator[T]) sortAll() *Paginator[T] {
	keys := make([]int, 0)
	for k := range p.keyed {
		keys = append(keys, k)
	}
	for _, sorter := range p.sorters {
		sorter.sorted = make([]int, len(keys))
		copy(sorter.sorted, keys)
		sort.Sort(sorter)
	}
	return p
}

// Override the same key.
func (p *Paginator[T]) AddSorter(key string, sorter *Sorter[T]) *Paginator[T] {
	sorter.paginator = p
	for k := range p.keyed {
		sorter.sorted = append(sorter.sorted, k)
	}
	p.sorters[key] = sorter
	sort.Sort(sorter)
	return p
}

func (p *Paginator[T]) RemoveFirst(selected func(t T) bool) *Paginator[T] {
	for k, v := range p.keyed {
		if selected(v) {
			delete(p.keyed, k)
			break
		}
	}
	p.sortAll()
	return p
}

func (p *Paginator[T]) RemoveAll(selected func(t T) bool) *Paginator[T] {
	removes := make([]int, 0)
	for k, v := range p.keyed {
		if selected(v) {
			removes = append(removes, k)
		}
	}
	for _, k := range removes {
		delete(p.keyed, k)
	}
	p.sortAll()
	return p
}

func (p *Paginator[T]) Adds(items ...T) *Paginator[T] {
	itemKeys := make([]int, 0)
	for _, item := range items {
		key := len(p.keyed)
		itemKeys = append(itemKeys, key)
		p.keyed[key] = item
	}
	p.sortAll()
	return p
}

func (p *Paginator[T]) Sorted(key string) ([]T, error) {
	results := make([]T, 0)
	// FIXME: 需要优化
	sorter, ok := p.sorters[key]
	if !ok {
		return nil, fmt.Errorf("No sorter named %s", key)
	}
	for _, k := range sorter.sorted {
		results = append(results, sorter.paginator.keyed[k])
	}
	return results, nil
}

// Page index starts from 1, but not 0.
// If index is less than 1, it will be the first page.
// If index is greater than the total page count, it will be the last page.
func (p *Paginator[T]) Paged(key string, size, index int) ([]T, error) {
	if size <= 0 {
		return nil, fmt.Errorf("Page size out of range")
	}
	start := 0
	if index > 1 {
		start = size * (index - 1)
	}
	results, err := p.Sorted(key)
	if err != nil {
		return nil, err
	}
	for start >= len(results) {
		start -= size
	}
	if start < 0 {
		start = 0
	}
	return results[start:min(start+size, len(results))], nil
}

func NewPaginator[T any](items ...T) *Paginator[T] {
	p := &Paginator[T]{
		keyed:   make(map[int]T),
		sorters: make(map[string]*Sorter[T]),
	}
	for i, t := range items {
		p.keyed[i] = t
	}
	return p
}
