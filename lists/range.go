package lists

func RangeTo(to int) []int {
	return Range(0, to)
}

func Range(from, to int) []int {
	s := 1
	if from > to {
		s = -1
	}
	return RangeStep(from, to, s)
}

func RangeStep(from, to, step int) []int {
	l := to - from
	c := l / step
	if l%step != 0 {
		c++
	}
	v := make([]int, 0, c)
	for i, x := 0, from; i < c; i, x = i+1, x+step {
		v = append(v, x)
	}
	return v
}
