package tools

func Map[T any, R any](s []T, f func(t T, i int) R) []R {
	r := make([]R, len(s))

	for i, v := range s {
		r[i] = f(v, i)
	}

	return r
}
