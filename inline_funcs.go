package g

// IfElse returns `result` if `condition` is true, otherwise `alternative`.
// Useful for inline conditional expressions in builder-style code.
func IfElse[T any](condition bool, result, alternative T) T {
	if condition {
		return result
	}
	return alternative
}

// If returns `result` when `condition` is true, otherwise an empty Node.
// This avoids nils when conditionally rendering DOM fragments.
func If(condition bool, result Node) Node {
	if condition {
		return result
	}
	return Empty()
}

// Repeat calls `f` exactly `n` times and aggregates the resulting Nodes.
// The passed function is used to ensure each Node instance is unique.
func Repeat(n int, f func() Node) Node {
	result := make([]Node, n)
	for i := range n {
		result[i] = f()
	}
	return Empty().Add(result...)
}

// Map converts a slice into Nodes by applying `f` to each element and
// aggregating the results into a single Node.
func Map[T any](input []T, f func(T) Node) Node {
	result := make([]Node, len(input))
	for i, item := range input {
		result[i] = f(item)
	}
	return Empty().Add(result...)
}
