package pago

import "testing"

func TestKeyManager(t *testing.T) {
	km := NewKeyManager(1)
	_, b, c := km.Allocate(), km.Allocate(), km.Allocate()
	if km.top != 4 {
		t.Fatalf("KeyManager top should be 4, but actually %d", km.top)
	}
	km.Restore(b)
	if len(km.temps) != 1 {
		t.Fatalf("KeyManager temps should have only 1 item, but actually %d", len(km.temps))
	}
	km.Restore(c)
	if km.top != 2 {
		t.Fatalf("KeyManager top should be 2, but actually %d", km.top)
	}
	_ = km.Allocate()
	if len(km.temps) != 0 {
		t.Fatalf("KeyManager temps should have no item, but actually %d", len(km.temps))
	}
	if km.top != 3 {
		t.Fatalf("KeyManager top should be 3, but actually %d", km.top)
	}
}
