package generic

import (
	"testing"
)

func TestMap(t *testing.T) {
	t.Run("SquareInts", func(t *testing.T) {
		// Map over ints, returning the square of each int.
		// (Take ints, return ints.)
		x := Map(func(a int) int {
			return a * a
		}, []int{2, 3, 4})

		expected := []int{4, 9, 16}
		for i, v := range x {
			if v != expected[i] {
				t.Errorf("Index %d: expected %d, got %d", i, expected[i], v)
			}
		}
	})

	t.Run("LengthOfStrings", func(t *testing.T) {
		// Map over strings, returning the length of each string.
		// (Take ints, return strings.)
		x := Map(func(a string) int {
			return len([]rune(a))
		}, []string{"one", "two", "three"})

		expected := []int{3, 3, 5}
		for i, v := range x {
			if v != expected[i] {
				t.Errorf("Index %d: expected %d, got %d", i, expected[i], v)
			}
		}
	})
}
