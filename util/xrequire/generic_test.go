package xrequire

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"copybooktogo/util/xrequire/mocks"
)

func TestSingle(t *testing.T) {
	t.Run("ExactlyOne", func(t *testing.T) {
		mockT := mocks.NewTestingT(t)
		x := Single(mockT, []int{100})
		assert.Equal(t, 100, x)
	})

	t.Run("BoundaryConditions", func(t *testing.T) {
		tests := []struct {
			Name string
			List []int
		}{
			{
				Name: "LessThanOne",
				List: []int{},
			},
			{
				Name: "MoreThanOne",
				List: []int{1, 2},
			},
		}

		for _, tt := range tests {
			t.Run(tt.Name, func(t *testing.T) {
				assert.PanicsWithValue(t, "mock fail now",
					func() {
						_ = Single(mockFailRequireT(t), tt.List)
						// not called due to test ending early
						assert.FailNow(t, "unreachable")
					},
				)
			})
		}
	})
}

func mockFailRequireT(t *testing.T) *mocks.TestingT {
	mockT := mocks.NewTestingT(t)
	mockT.EXPECT().Errorf(mock.Anything, mock.Anything)
	mockT.EXPECT().FailNow().Panic("mock fail now")
	return mockT
}
