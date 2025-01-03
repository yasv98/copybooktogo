package generic

import (
	"testing"

	"github.com/stretchr/testify/assert"
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

func TestMergeMaps(t *testing.T) {
	t.Run("StringKeysIntValues", func(t *testing.T) {
		map1 := map[string]int{"a": 1, "b": 2}
		map2 := map[string]int{"c": 3}

		expected := map[string]int{"a": 1, "b": 2, "c": 3}
		actual := MergeMaps(map1, map2)

		assert.Equal(t, expected, actual)
	})

	t.Run("IntKeysStringValues", func(t *testing.T) {
		map1 := map[int]string{1: "a", 2: "b"}
		map2 := map[int]string{3: "c"}

		expected := map[int]string{1: "a", 2: "b", 3: "c"}
		actual := MergeMaps(map1, map2)

		assert.Equal(t, expected, actual)
	})

	t.Run("ThreeMaps", func(t *testing.T) {
		map1 := map[string]int{"a": 1}
		map2 := map[string]int{"b": 2, "c": 3}
		map3 := map[string]int{"d": 4}

		expected := map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}
		actual := MergeMaps(map1, map2, map3)

		assert.Equal(t, expected, actual)
	})

	t.Run("Overwrite", func(t *testing.T) {
		map1 := map[string]int{"a": 1, "b": 2}
		map2 := map[string]int{"b": 3}

		expected := map[string]int{"a": 1, "b": 3}
		actual := MergeMaps(map1, map2)

		assert.Equal(t, expected, actual)
	})

	t.Run("OverwriteMultipleMapsWithSameKey", func(t *testing.T) {
		map1 := map[string]int{"a": 1, "b": 2}
		map2 := map[string]int{"b": 3}
		map3 := map[string]int{"b": 4}

		expected := map[string]int{"a": 1, "b": 4}
		actual := MergeMaps(map1, map2, map3)

		assert.Equal(t, expected, actual)
	})

	t.Run("OneEmptyMap", func(t *testing.T) {
		map1 := map[string]int{"a": 1}
		map2 := map[string]int{}

		expected := map[string]int{"a": 1}
		actual := MergeMaps(map1, map2)

		assert.Equal(t, expected, actual)
	})

	t.Run("AllEmpty", func(t *testing.T) {
		map1 := make(map[int]string)
		map2 := make(map[int]string)

		expected := make(map[int]string)
		actual := MergeMaps(map1, map2)

		assert.Equal(t, expected, actual)
	})
}
