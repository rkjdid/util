package util

import (
	"sync"
	"testing"
	"time"
)

func sliceEq(a []int, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestDatalog_Add(t *testing.T) {
	dtest := NewTimeSeries(5, Duration(time.Second))

	dtest.Add(0)
	expect := []int{0}
	if !sliceEq(dtest.Data, expect) {
		t.Errorf("%v != %v", dtest.Data, expect)
	}

	for i := 1; i < 5; i++ {
		dtest.Add(i)
	}
	expect = []int{0, 1, 2, 3, 4}
	if !sliceEq(dtest.Data, expect) {
		t.Errorf("%v != %v", dtest.Data, expect)
	}

	dtest.Add(5)
	expect = []int{1, 2, 3, 4, 5}
	if !sliceEq(dtest.Data, expect) {
		t.Errorf("%v != %v", dtest.Data, expect)
	}
}

func TestDatalog_Padded(t *testing.T) {
	dtest := NewTimeSeries(5, Duration(time.Second))

	expect := make([]int, 5)
	if p := dtest.Padded(); !sliceEq(p, expect) {
		t.Errorf("%v != %v", p, expect)
	}

	dtest.Add(1)
	expect = []int{1, 1, 1, 1, 1}
	if p := dtest.Padded(); !sliceEq(p, expect) {
		t.Errorf("%v != %v", p, expect)
	}

	dtest.Add(2)
	expect[4] = 2
	if p := dtest.Padded(); !sliceEq(p, expect) {
		t.Errorf("%v != %v", p, expect)
	}
}

func TestDatalog_SetMaxLength(t *testing.T) {
	dtest := NewTimeSeries(2, Duration(time.Second))
	t0 := dtest.Start
	dtest.Add(1)
	dtest.Add(2)
	dtest.SetMaxLength(1)
	if expect := []int{2}; !sliceEq(dtest.Data, expect) {
		t.Errorf("%v != %v", dtest.Data, expect)
	}
	if dtest.Start.Sub(t0) != time.Second {
		t.Errorf("start time should've shifted by 1sec, got %s", dtest.Start.Sub(t0))
	}
}

func TestDatalog_Subscribe(t *testing.T) {
	dtest := NewTimeSeries(2, Duration(time.Second))
	i, ch := dtest.Subscribe()

	var wg sync.WaitGroup
	wg.Add(1)
	expect := []int{}
	go func() {
		for x := range ch {
			expect = append(expect, x)
		}
		wg.Done()
	}()

	dtest.Add(1)
	dtest.Add(2)
	dtest.Unsubscribe(i)
	wg.Wait()
	if !sliceEq(dtest.Data, []int{1, 2}) {
		t.Errorf("%v != %v", dtest.Data, expect)
	}
	if !sliceEq(dtest.Data, expect) {
		t.Errorf("%v != %v", dtest.Data, expect)
	}
}
