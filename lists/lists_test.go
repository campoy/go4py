package lists

import (
	"reflect"
	"testing"
)

func checkRange(t *testing.T, from, to, step int, r []int) {
	if len(r) != cap(r) {
		t.Errorf("%v (%v, %v, %v) - len(r) %v - cap(r) %v", r, from, to, step, len(r), cap(r))
	}
	for i, v := range r {
		if e := from + i*step; e != v {
			t.Errorf("%v - item %v should be %v, is %v", r, i, e, v)
		}
	}
}

func TestRangeTo(t *testing.T) {
	for to := -10; to < 0; to++ {
		checkRange(t, 0, to, -1, RangeTo(to))
	}
	for to := 1; to < 10; to++ {
		checkRange(t, 0, to, 1, RangeTo(to))
	}
}

func TestRange(t *testing.T) {
	for from := 0; from < 10; from++ {
		for to := from; to < 10; to++ {
			checkRange(t, from, to, 1, Range(from, to))
		}

		for to := 10; to > from; to-- {
			checkRange(t, from, to, 1, Range(from, to))
		}
	}
}

func TestRangeStep(t *testing.T) {
	for from := 0; from < 10; from++ {
		for to := from; to < 10; to++ {
			for step := 1; step < 10; step++ {
				checkRange(t, from, to, step, RangeStep(from, to, step))
			}
		}
	}
}

func TestMapIntInt(t *testing.T) {
	cases := []struct {
		f   func(int) int
		in  []int
		out []int
	}{
		{nil, nil, nil},
		{nil, []int{}, []int{}},
		{func(x int) int { return x }, Range(0, 10), Range(0, 10)},
		{func(x int) int { return 2 * x }, Range(0, 10), RangeStep(0, 20, 2)},
	}

	for _, c := range cases {
		r := MapIntInt(c.f, c.in)
		if len(r) == 0 && len(c.out) == 0 {
			continue
		}
		if !reflect.DeepEqual(r, c.out) {
			t.Errorf("expected %v, got %v", c.out, r)
		}
	}
}

func TestFilterInt(t *testing.T) {
	cases := []struct {
		f   func(int) bool
		in  []int
		out []int
	}{
		{nil, nil, nil},
		{nil, []int{}, []int{}},
		{func(x int) bool { return true }, Range(0, 10), Range(0, 10)},
		{func(x int) bool { return false }, Range(0, 10), nil},
		{func(x int) bool { return x%2 == 0 }, Range(0, 10), RangeStep(0, 10, 2)},
		{func(x int) bool { return x%2 != 0 }, Range(0, 10), RangeStep(1, 10, 2)},
	}

	for _, c := range cases {
		r := FilterInt(c.f, c.in)
		if len(r) == 0 && len(c.out) == 0 {
			continue
		}
		if !reflect.DeepEqual(r, c.out) {
			t.Errorf("expected %v, got %v", c.out, r)
		}
	}
}

func TestMap(t *testing.T) {
	cases := []struct {
		f   interface{}
		in  interface{}
		out interface{}
	}{
		{func(x int) int { return x }, Range(0, 10), Range(0, 10)},
		{func(x int) int { return 2 * x }, Range(0, 10), RangeStep(0, 20, 2)},
		{func(x int) float32 { return 1.5 * float32(x) }, Range(0, 4), []float32{0, 1.5, 3, 4.5}},
	}

	for _, c := range cases {
		r := Map(c.f, c.in)
		if err, ok := r.(error); ok {
			t.Error(err)
			continue
		}
		if reflect.ValueOf(r).Len() == 0 && reflect.ValueOf(c.out).Len() == 0 {
			continue
		}
		if !reflect.DeepEqual(r, c.out) {
			t.Errorf("expected %v, got %v", c.out, r)
		}
	}
}

func TestFilter(t *testing.T) {
	cases := []struct {
		f   interface{}
		in  interface{}
		out interface{}
	}{
		{nil, []int{}, []int{}},
		{func(x int) bool { return true }, Range(0, 10), Range(0, 10)},
		{func(x int) bool { return false }, Range(0, 10), []int{}},
		{func(x int) bool { return x%2 == 0 }, Range(0, 10), RangeStep(0, 10, 2)},
		{func(x int) bool { return x%2 != 0 }, Range(0, 10), RangeStep(1, 10, 2)},
	}

	for _, c := range cases {
		r := Filter(c.f, c.in)
		if err, ok := r.(error); ok {
			t.Error(err)
			continue
		}
		if reflect.ValueOf(r).Len() == 0 && reflect.ValueOf(c.out).Len() == 0 {
			continue
		}
		if !reflect.DeepEqual(r, c.out) {
			t.Errorf("expected %v, got %v", c.out, r)
		}
	}
}

func BenchmarkConcrete(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FilterInt(func(x int) bool { return x%2 == 0 },
			MapIntInt(func(x int) int { return x * x },
				Range(0, 11)))
	}
}

func BenchmarkGeneric(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Filter(func(x int) bool { return x%2 == 0 },
			Map(func(x int) int { return x * x },
				Range(0, 11)))
	}
}
