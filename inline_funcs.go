package g

func If[T any](condition bool, result T) T {
	if condition {
		return result
	}
	var zero T
	return zero
}

func IfElse[T any](condition bool, result, alternative T) T {
	if condition {
		return result
	}
	return alternative
}

// this accepts a function that returns a node
// this ensure each node is different
func Repeat(n int, f func() Node) Node {
	result := make([]Node, n)
	for i := range n {
		result[i] = f()
	}
	return Empty().Add(result...)
}

func Map[T any](input []T, f func(T) Node) Node {
	result := make([]Node, len(input))
	for i, item := range input {
		result[i] = f(item)
	}
	return Empty().Add(result...)
}
