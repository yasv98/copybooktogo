// Package xrequire implements additional common test helpers on top of testify
// assert.
package xrequire

import (
	"github.com/stretchr/testify/require"
)

// Single asserts that the list has exactly one element and returns it.
func Single[T any](t require.TestingT, list []T) T {
	require.Len(t, list, 1)
	return list[0]
}
