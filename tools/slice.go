package tools

func Reduce[T any, U any](s []T, f func(p, n T) U) U {
	var r U = f()

	return r
}
