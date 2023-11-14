package pago

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
)

type Pago[T any] struct {
	dirty      bool
	keyed      map[int]T
	keyManager *KeyManager
	sorters    map[string]*Sorter[T]
}

func (p *Pago[T]) sortAll() *Pago[T] {
	keys := make([]int, 0)
	for k := range p.keyed {
		keys = append(keys, k)
	}
	for _, sorter := range p.sorters {
		sorter.sorted = make([]int, len(keys))
		copy(sorter.sorted, keys)
		sort.Sort(sorter)
	}
	p.dirty = false
	return p
}

// Override the same key.
func (p *Pago[T]) AddSorter(key string, sorter *Sorter[T]) *Pago[T] {
	sorter.pago = p
	for k := range p.keyed {
		sorter.sorted = append(sorter.sorted, k)
	}
	p.sorters[key] = sorter
	sort.Sort(sorter)
	return p
}

func (p *Pago[T]) RemoveSorter(key string) *Pago[T] {
	delete(p.sorters, key)
	return p
}

func (p *Pago[T]) RemoveFirstBy(selectors ...func(t T) bool) *Pago[T] {
next:
	for k, v := range p.keyed {
		for _, selected := range selectors {
			if !selected(v) {
				continue next
			}
		}
		delete(p.keyed, k)
		p.keyManager.Restore(k)
		p.dirty = true
		break
	}
	return p
}

func (p *Pago[T]) RemoveAllBy(selectors ...func(t T) bool) *Pago[T] {
	removes := make([]int, 0)
next:
	for k, v := range p.keyed {
		for _, selected := range selectors {
			if !selected(v) {
				continue next
			}
		}
		removes = append(removes, k)
	}
	if len(removes) > 0 {
		p.dirty = true
	}
	for _, k := range removes {
		delete(p.keyed, k)
		p.keyManager.Restore(k)
	}
	return p
}

func (p *Pago[T]) Adds(items ...T) *Pago[T] {
	if len(items) > 0 {
		p.dirty = true
	}
	for _, item := range items {
		// TODO: KeyManager分配Key的一致性需要测试，可能出现意外的覆盖
		v := p.keyManager.Allocate()
		if _, ok := p.keyed[v]; ok {
			panic("Repeated key: " + strconv.Itoa(v))
		}
		p.keyed[v] = item
	}
	return p
}

func (p *Pago[T]) Sorted(key string, selectors ...func(t T) bool) ([]T, error) {
	// FIXME: optimize
	sorter, ok := p.sorters[key]
	if !ok {
		return nil, fmt.Errorf("No sorter named %s", key)
	}
	if p.dirty {
		p.sortAll()
	}
	results := make([]T, 0)
next:
	for _, k := range sorter.sorted {
		v := sorter.pago.keyed[k]
		for _, selected := range selectors {
			if !selected(v) {
				continue next
			}
		}
		results = append(results, v)
	}
	return results, nil
}

// Page index starts from 1, but not 0.
//
// If index is lower than 1, the first page will be returned.
// If index is greater than the total page count, the last page will be returned.
func (p *Pago[T]) Paged(key string, size, index int, selectors ...func(t T) bool) ([]T, error) {
	if size < 1 {
		return nil, errors.New("Page size can not be lower than 1")
	}
	results, err := p.Sorted(key, selectors...)
	if err != nil {
		return nil, err
	}
	if index < 1 {
		index = 1
	}
	start := size * (index - 1)
	if start > len(results) {
		start = len(results) / size * size
	}
	return results[start:min[int](start+size, len(results))], nil
}

// Return the total page count by the specified page size and selectors.
// If page size is less than 1, then -1 will be returned.
func (p *Pago[T]) Count(size int, selectors ...func(t T) bool) int {
	if size < 1 {
		return -1
	}
	sum := 0
next:
	for _, v := range p.keyed {
		for _, selected := range selectors {
			if !selected(v) {
				continue next
			}
		}
		sum++
	}
	cnt := sum / size
	if sum%size > 0 {
		return cnt + 1
	}
	return cnt
}

func NewPago[T any](items ...T) *Pago[T] {
	p := &Pago[T]{
		keyed:      make(map[int]T),
		keyManager: NewKeyManager(1),
		sorters:    make(map[string]*Sorter[T]),
	}
	p.Adds(items...)
	return p
}
