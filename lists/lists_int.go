package lists

func FilterInt(q func(int) bool, l []int) []int {
	var v []int
	for _, x := range l {
		if q(x) {
			v = append(v, x)
		}
	}
	return v
}

func MapIntInt(f func(int) int, l []int) []int {
	v := make([]int, len(l))
	for i, x := range l {
		v[i] = f(x)
	}
	return v
}
