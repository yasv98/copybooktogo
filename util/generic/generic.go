// Package generic contains a collection of helpers for generic types and collections that don't but should belong in stdlib.
package generic

// Map applies a given function to each element of a list, returning a new list.
func Map[A, B any](f func(A) B, list []A) []B {
	res := make([]B, len(list))
	for i, v := range list {
		res[i] = f(v)
	}
	return res
}
