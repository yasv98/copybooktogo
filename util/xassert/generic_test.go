package xassert

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestEqualAllFunc(t *testing.T) {
	t.Parallel()

	t.Run("True", func(t *testing.T) {
		t.Run("AllEqual_Single", func(t *testing.T) {
			t.Parallel()
			expected := "1"
			actuals := []int{1}
			assert.True(t, EqualAllFunc(mockT{}, expected, actuals, strconv.Itoa))
		})

		t.Run("AllEqual_Multiple", func(t *testing.T) {
			t.Parallel()
			expected := "1"
			actuals := []int{1, 1, 1, 1}
			assert.True(t, EqualAllFunc(mockT{}, expected, actuals, strconv.Itoa))
		})

		t.Run("AllEqual_Empty", func(t *testing.T) {
			t.Parallel()
			expected := "1"
			var actuals []int
			assert.True(t, EqualAllFunc(mockT{}, expected, actuals, strconv.Itoa))
		})
	})

	t.Run("False", func(t *testing.T) {
		t.Run("NotEqual_Single", func(t *testing.T) {
			t.Parallel()
			expected := "1"
			actuals := []int{0}
			assert.False(t, EqualAllFunc(mockT{}, expected, actuals, strconv.Itoa))
		})

		t.Run("SomeEqual", func(t *testing.T) {
			t.Parallel()
			expected := "1"
			actuals := []int{1, 1, 0, 1}
			assert.False(t, EqualAllFunc(mockT{}, expected, actuals, strconv.Itoa))
		})

		t.Run("NoneEqual", func(t *testing.T) {
			t.Parallel()
			expected := "1"
			actuals := []int{0, 0, 0, 0}
			assert.False(t, EqualAllFunc(mockT{}, expected, actuals, strconv.Itoa))
		})
	})
}

type mockT struct {
	*mock.Mock
}

func (m mockT) Errorf(_ string, _ ...interface{}) {}
func (m mockT) FailNow() {
	m.Mock.Called("FailNow")
}
