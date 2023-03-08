package gopagination

import (
	"fmt"
	"testing"
)

type testElem struct {
	key int
}

type testComparator struct{}

func (c *testComparator) Less(i, j *testElem) bool {
	return i.key < j.key
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
	sorter := NewSorter[*testElem](&testComparator{})
	paginator := NewPaginator(
		&testElem{key: 2},
		&testElem{key: 3},
		&testElem{key: 1},
		&testElem{key: 1},
		&testElem{key: -1},
		&testElem{key: 100},
	)
	sorted, err := paginator.AddSorter("test", sorter).Sorted("test")
	if err != nil {
		t.Fatal(err)
	}
	expectedOrder := []int{-1, 1, 1, 2, 3, 100}
	if err := checkOrdered(expectedOrder, sorted); err != nil {
		t.Fatal(err)
	}
}

func TestPage(t *testing.T) {
	sorter := NewSorter[*testElem](&testComparator{})
	paginator := NewPaginator(
		&testElem{key: 1},
		&testElem{key: 6},
		&testElem{key: 2},
		&testElem{key: 4},
		&testElem{key: 5},
		&testElem{key: 3},
	)
	paginator.AddSorter("test", sorter)

	type pageTest struct {
		size, index   int
		expectedOrder []int
	}
	pageTests := []pageTest{
		{size: 3, index: 1, expectedOrder: []int{1, 2, 3}},
		{size: 4, index: 2, expectedOrder: []int{5, 6}},
		{size: 4, index: 3, expectedOrder: []int{5, 6}},
		{size: 100, index: 3, expectedOrder: []int{1, 2, 3, 4, 5, 6}},
	}
	for _, pt := range pageTests {
		paged, err := paginator.Paged("test", pt.size, pt.index)
		if err != nil {
			t.Fatal(err)
		}
		if err := checkOrdered(pt.expectedOrder, paged); err != nil {
			t.Fatal(err)
		}
	}
}

func TestRemove(t *testing.T) {
	sorter := NewSorter[*testElem](&testComparator{})
	paginator := NewPaginator(
		&testElem{key: 4},
		&testElem{key: 2},
		&testElem{key: 3},
		&testElem{key: 3},
		&testElem{key: 1},
	)
	paginator.AddSorter("test", sorter)

	paginator.RemoveFirst(func(t *testElem) bool { return t.key == 2 })
	paged, err := paginator.Paged("test", 3, 1)
	if err != nil {
		t.Fatal(err)
	}
	if err := checkOrdered([]int{1, 3, 3}, paged); err != nil {
		t.Fatal(err)
	}

	paginator.RemoveAll(func(t *testElem) bool { return t.key == 3 })
	paged, err = paginator.Paged("test", 4, 2)
	if err != nil {
		t.Fatal(err)
	}
	if err := checkOrdered([]int{1, 4}, paged); err != nil {
		t.Fatal(err)
	}

	paginator.RemoveAll(func(t *testElem) bool { return true })
	paged, err = paginator.Paged("test", 4, 3)
	if err != nil {
		t.Fatal(err)
	}
	if err := checkOrdered([]int{}, paged); err != nil {
		t.Fatal(err)
	}
}
