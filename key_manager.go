package pago

import "sort"

type KeyManager struct {
	base  int
	top   int
	temps map[int]struct{}
}

func (m *KeyManager) retop() {
	ints := make(sort.IntSlice, 0)
	for k := range m.temps {
		ints = append(ints, k)
	}
	sort.Sort(sort.Reverse(ints))
	for _, v := range ints {
		if v != m.top-1 {
			return
		}
		delete(m.temps, v)
		m.top--
	}
}

func (m *KeyManager) Restore(v int) {
	if v < m.base {
		return
	}
	m.temps[v] = struct{}{}
	m.retop()
}

func (m *KeyManager) Allocate() int {
	for k := range m.temps {
		delete(m.temps, k)
		return k
	}
	k := m.top
	m.top++
	return k
}

func NewKeyManager(base int) *KeyManager {
	if base < 0 {
		base = 0
	}
	return &KeyManager{base: base, top: base, temps: make(map[int]struct{})}
}
