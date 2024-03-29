package pago

import (
	"fmt"
	"testing"
)

type testElem struct {
	key int
}

type positiveComparator struct{}

func (c *positiveComparator) Less(i, j *testElem) bool {
	return i.key < j.key
}

type negativeComparator struct{}

func (c *negativeComparator) Less(i, j *testElem) bool {
	return i.key > j.key
}

func checkOrdered(expect []int, actual []*testElem) error {
	if len(expect) != len(actual) {
		return fmt.Errorf("Element count not equals to the expected")
	}
	for i := 0; i < len(actual); i++ {
		if expect[i] != actual[i].key {
			return fmt.Errorf("Order not equals to the expected")
		}
	}
	return nil
}

func TestSort(t *testing.T) {
	sorter := NewSorter[*testElem](&positiveComparator{})
	pago := NewPago(
		&testElem{key: 2},
		&testElem{key: 3},
		&testElem{key: 1},
		&testElem{key: 1},
		&testElem{key: -1},
		&testElem{key: 100},
	)
	sorted, err := pago.AddSorter("test", sorter).Sorted("test")
	if err != nil {
		t.Fatal(err)
	}
	expectedOrder := []int{-1, 1, 1, 2, 3, 100}
	if err := checkOrdered(expectedOrder, sorted); err != nil {
		t.Fatal(err)
	}
}

func TestPage(t *testing.T) {
	sorter := NewSorter[*testElem](&positiveComparator{})
	pago := NewPago(
		&testElem{key: 1},
		&testElem{key: 6},
		&testElem{key: 2},
		&testElem{key: 4},
		&testElem{key: 5},
		&testElem{key: 3},
	)
	pago.AddSorter("test", sorter)

	if pages := pago.Count(2); pages != 3 {
		t.Errorf("Total page count(size=3) should be 2, but actully %d", pages)
	}
	if pages := pago.Count(5); pages != 2 {
		t.Errorf("Total page count(size=5) should be 2, but actully %d", pages)
	}

	type pageTest struct {
		size, index   int
		expectedOrder []int
	}
	pageTests := []pageTest{
		{size: 3, index: 1, expectedOrder: []int{1, 2, 3}},
		{size: 4, index: 2, expectedOrder: []int{5, 6}},
		{size: 4, index: 3, expectedOrder: []int{5, 6}},
		{size: 4, index: 6, expectedOrder: []int{5, 6}},
		{size: 100, index: 3, expectedOrder: []int{1, 2, 3, 4, 5, 6}},
	}
	for _, pt := range pageTests {
		paged, err := pago.Paged("test", pt.size, pt.index)
		if err != nil {
			t.Fatal(err)
		}
		if err := checkOrdered(pt.expectedOrder, paged); err != nil {
			t.Fatal(err)
		}
	}
}

func TestRemove(t *testing.T) {
	sorter := NewSorter[*testElem](&positiveComparator{})
	pago := NewPago(
		&testElem{key: 4},
		&testElem{key: 2},
		&testElem{key: 3},
		&testElem{key: 3},
		&testElem{key: 1},
	)
	pago.AddSorter("test", sorter)

	pago.RemoveFirstBy(func(t *testElem) bool { return t.key == 2 })
	index := 1
	paged, err := pago.Paged("test", 3, index)
	if err != nil {
		t.Fatal(err)
	}
	if err := checkOrdered([]int{1, 3, 3}, paged); err != nil {
		t.Fatal(err)
	}

	index = 2
	pago.RemoveAllBy(func(t *testElem) bool { return t.key == 3 })
	paged, err = pago.Paged("test", 4, index)
	if err != nil {
		t.Fatal(err)
	}
	if err := checkOrdered([]int{1, 4}, paged); err != nil {
		t.Fatal(err)
	}

	index = 3
	pago.RemoveAllBy(func(t *testElem) bool { return true })
	paged, err = pago.Paged("test", 4, index)
	if err != nil {
		t.Fatal(err)
	}

	if err := checkOrdered([]int{}, paged); err != nil {
		t.Fatal(err)
	}
}

func TestMultiSorter(t *testing.T) {
	positiveSorter := NewSorter[*testElem](&positiveComparator{})
	negativeSorter := NewSorter[*testElem](&negativeComparator{})
	pago := NewPago(
		&testElem{key: 4},
		&testElem{key: 2},
		&testElem{key: 3},
		&testElem{key: 5},
		&testElem{key: 1},
	).AddSorter("positive", positiveSorter).AddSorter("negative", negativeSorter)

	positiveSorted, err := pago.Sorted("positive")
	if err != nil {
		t.Error(err)
	}
	expectedOrderPositive := []int{1, 2, 3, 4, 5}
	if err := checkOrdered(expectedOrderPositive, positiveSorted); err != nil {
		t.Fatal(err)
	}

	negativeSorted, err := pago.Sorted("negative")
	if err != nil {
		t.Error(err)
	}
	expectedOrderNegative := []int{5, 4, 3, 2, 1}
	if err := checkOrdered(expectedOrderNegative, negativeSorted); err != nil {
		t.Fatal(err)
	}
}
