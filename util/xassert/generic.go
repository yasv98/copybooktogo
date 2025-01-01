// Package xassert implements additional common test helpers on top of testify
// assert.
package xassert

import (
	"context"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"copybooktogo/util/generic"
)

// EqualAll asserts that all elements are equal to the expected value.
func EqualAll[T any](t assert.TestingT, expected T, actuals []T) bool {
	return EqualAllFunc(t, expected, actuals, func(elem T) T { return elem })
}

// EqualAllFunc asserts that all elements are equal to the expected value after
// applying the given function to each element.
func EqualAllFunc[T, K any](t assert.TestingT, expected K, actuals []T, f func(elem T) K) bool {
	got := generic.Map(f, actuals)
	for _, actual := range got {
		if !assert.Equal(t, expected, actual, "all elements should match %v, got: %v", expected, got) {
			return false
		}
	}
	return true
}

// MockContextArg returns a mock context.Context argument to be used in testify mock argument assertions.
func MockContextArg() any {
	return mock.MatchedBy(func(context.Context) bool {
		return true
	})
}
