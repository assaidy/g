package utils

import "github.com/assaidy/g"

// IfElse returns the appropriate value based on a boolean condition.
//
// This generic function is useful for inline conditional expressions in
// builder-style code where you need to choose between two values without
// breaking the chain of method calls.
//
// Example:
//
//	div := Div(KV{"class": IfElse(isActive, "active", "inactive")})
//
//	Body(
//		IfElse(isAdmin,
//			Div(g.Text("Admin content")),
//			P(g.Text("Regular user content")),
//		),
//	)
func IfElse[T any](condition bool, result, alternative T) T {
	if condition {
		return result
	}
	return alternative
}

// Conditionally returns a Node based on a boolean condition.
//
// This function returns an empty Node (not nil) when the
// condition is false, which prevents nil pointer issues when building
// DOM trees.
//
// Example:
//
//	Body(
//		If(showHeader, Header(...)),
//		Main(...),
//	)
func If(condition bool, result g.Node) g.Node {
	if condition {
		return result
	}
	return g.Empty()
}

// Repeat generates multiple Nodes by calling a function n times.
//
// The provided function is called exactly n times, and each resulting Node
// is aggregated into a single container Node. Using a function ensures each
// Node instance is unique (important for elements with mutable state).
//
// Example:
//
//	Ul(
//		Repeat(5, func() g.Node {
//			return Li(Text("List item"))
//		}),
//	)
func Repeat(n int, f func() g.Node) g.Node {
	result := g.Empty()
	for range n {
		result.Children = append(result.Children, f())
	}
	return result
}

// Map transforms a slice of items into Nodes by applying a function to each element.
//
// Each element in the input slice is transformed using the provided function, and
// all resulting Nodes are aggregated into a single container Node.
//
// Example:
//
//	items := []string{"Apple", "Banana", "Cherry"}
//	Ul(
//		Map(items, func(item string) g.Node {
//			return Li(Text(item))
//		}),
//	)
func Map[T any](input []T, f func(T) g.Node) g.Node {
	result := g.Empty()
	for _, item := range input {
		result.Children = append(result.Children, f(item))
	}
	return result
}
