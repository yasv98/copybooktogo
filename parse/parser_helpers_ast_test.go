package parse

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/yasv98/copybooktogo/util/xassert"
)

func Test_createAndAddRecordToAST(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		t.Run("Level01Record_AddedToASTAndWorkingParentsStack", func(t *testing.T) {
			builder := &astBuilder{}

			err := createAndAddRecordToAST(builder, 1, "LEVEL01-RECORD", []any{})

			require.NoError(t, err)
			expectedRecord := &Record{Level: 1, Identifier: "LEVEL01-RECORD"}
			require.Len(t, builder.ast, 1)
			require.Len(t, builder.workingParentsStack, 1)
			xassert.EqualAll(t, expectedRecord, []*Record{builder.ast[0], builder.workingParentsStack[0]})
		})

		t.Run("Level05Record_AddedToParentAndWorkingParentsStack", func(t *testing.T) {
			rootRecord := &Record{Level: 1, Identifier: "LEVEL01-RECORD"}
			builder := &astBuilder{ast: []*Record{rootRecord}, workingParentsStack: []*Record{rootRecord}}

			err := createAndAddRecordToAST(builder, 5, "LEVEL05-RECORD", []any{})

			require.NoError(t, err)
			require.Len(t, builder.workingParentsStack, 2)
			require.Len(t, builder.ast[0].Children, 1)
			expectedRecord := &Record{Level: 5, Identifier: "LEVEL05-RECORD"}
			xassert.EqualAll(t, expectedRecord, []*Record{builder.ast[0].Children[0], builder.workingParentsStack[1]})

			assert.Len(t, builder.ast, 1) // assert it is not added to ast
		})

		t.Run("Level01Record_AddedToASTAndReinitializeWorkingParentsStack", func(t *testing.T) {
			rootRecord := &Record{
				Level:      1,
				Identifier: "LEVEL01-RECORD-1",
				Children: []*Record{
					{
						Level:      5,
						Identifier: "LEVEL05-RECORD",
					},
				},
			}
			builder := &astBuilder{ast: []*Record{rootRecord}, workingParentsStack: []*Record{rootRecord, rootRecord.Children[0]}}

			err := createAndAddRecordToAST(builder, 1, "LEVEL01-RECORD-2", []any{})

			require.NoError(t, err)
			require.Len(t, builder.ast, 2)
			require.Len(t, builder.workingParentsStack, 1)
			expectedRecord := &Record{Level: 1, Identifier: "LEVEL01-RECORD-2"}
			xassert.EqualAll(t, expectedRecord, []*Record{builder.ast[1], builder.workingParentsStack[0]})
		})

		t.Run("Level05RecordWithPic_AddedToParent", func(t *testing.T) {
			rootRecord := &Record{Level: 1, Identifier: "LEVEL01-RECORD"}
			builder := &astBuilder{ast: []*Record{rootRecord}, workingParentsStack: []*Record{rootRecord}}

			err := createAndAddRecordToAST(builder, 5, "LEVEL05-RECORD", []any{Picture{PicType: Alpha, PicCount: 1}}) // a Record with a Picture clause is a leaf node

			require.NoError(t, err)
			require.Len(t, builder.ast[0].Children, 1)
			expectedRecord := &Record{Level: 5, Identifier: "LEVEL05-RECORD", Pic: Picture{PicType: Alpha, PicCount: 1}}
			assert.Equal(t, expectedRecord, builder.ast[0].Children[0])

			assert.Len(t, builder.workingParentsStack, 1) // assert it is not added to the working parent stack
			assert.Len(t, builder.ast, 1)                 // assert it is not added to ast
		})

		t.Run("Level10RecordWithPic_AddedToCorrectParent", func(t *testing.T) {
			rootRecord := &Record{
				Level:      1,
				Identifier: "LEVEL01-RECORD",
				Children: []*Record{
					{
						Level:      5,
						Identifier: "LEVEL05-RECORD",
						Children: []*Record{
							{
								Level:      10,
								Identifier: "LEVEL10-RECORD-1",
								Children: []*Record{
									{
										Level:      15,
										Identifier: "LEVEL15-RECORD",
									},
								},
							},
						},
					},
				},
			}
			builder := &astBuilder{ast: []*Record{rootRecord}, workingParentsStack: []*Record{rootRecord, rootRecord.Children[0], rootRecord.Children[0].Children[0]}}

			err := createAndAddRecordToAST(builder, 10, "LEVEL10-RECORD-2", []any{Picture{PicType: Alpha, PicCount: 1}})

			require.NoError(t, err)
			groupRecord := rootRecord.Children[0]
			require.Len(t, groupRecord.Children, 2) // Should be added to end of LEVEL05-RECORD's Children
			expectedRecord := &Record{Level: 10, Identifier: "LEVEL10-RECORD-2", Pic: Picture{PicType: Alpha, PicCount: 1}}
			assert.Equal(t, groupRecord.Children[1], expectedRecord)
		})

		t.Run("Level15Record_AddedToBottomOfTree", func(t *testing.T) {
			rootRecord := &Record{
				Level:      1,
				Identifier: "LEVEL01-RECORD",
				Children: []*Record{
					{
						Level:      5,
						Identifier: "LEVEL05-RECORD",
						Children: []*Record{
							{
								Level:      10,
								Identifier: "LEVEL10-RECORD",
							},
						},
					},
				},
			}
			builder := &astBuilder{ast: []*Record{rootRecord}, workingParentsStack: []*Record{rootRecord, rootRecord.Children[0], rootRecord.Children[0].Children[0]}}

			err := createAndAddRecordToAST(builder, 15, "LEVEL15-RECORD", []any{Picture{PicType: Alpha, PicCount: 1}})

			require.NoError(t, err)
			subGroupRecord := rootRecord.Children[0].Children[0]
			require.Len(t, subGroupRecord.Children, 1) // Should be added LEVEL10-RECORD's Children
			expectedRecord := &Record{Level: 15, Identifier: "LEVEL15-RECORD", Pic: Picture{PicType: Alpha, PicCount: 1}}
			assert.Equal(t, subGroupRecord.Children[0], expectedRecord)
		})
	})

	t.Run("Fail", func(t *testing.T) {
		t.Run("InvalidASTBuilder", func(t *testing.T) {
			err := createAndAddRecordToAST("not an astBuilder", 1, "LEVEL01-RECORD", []any{})
			assert.Error(t, err)
		})

		t.Run("NotRootRecordAddedToEmptyStack", func(t *testing.T) {
			err := createAndAddRecordToAST(&astBuilder{}, 5, "LEVEL05-RECORD", []any{})
			assert.Error(t, err)
		})

		t.Run("InvalidRecordLevel", func(t *testing.T) {
			err := createAndAddRecordToAST(&astBuilder{}, -1, "INVALID-RECORD", []any{})
			assert.Error(t, err)
		})
	})
}

func Test_isLeafNode(t *testing.T) {
	t.Run("LeafNodeRecord_ReturnsTrue", func(t *testing.T) {
		leafNodeRecord := &Record{Pic: Picture{PicType: Alpha, PicCount: 1}}
		assert.True(t, isLeafNode(leafNodeRecord))
	})

	t.Run("NonLeafNodeRecord_ReturnsFalse", func(t *testing.T) {
		nonLeafNodeRecord := &Record{}
		assert.False(t, isLeafNode(nonLeafNodeRecord))
	})
}

func Test_getAST(t *testing.T) {
	t.Run("Success_ReturnsAST", func(t *testing.T) {
		rootRecord1 := &Record{Level: 1, Identifier: "LEVEL01-RECORD-1"}
		rootRecord2 := &Record{Level: 1, Identifier: "LEVEL01-RECORD-2"}
		builder := &astBuilder{ast: []*Record{rootRecord1, rootRecord2}}

		got, err := getAST(builder)

		require.NoError(t, err)
		require.Len(t, got, 2)
		assert.Equal(t, builder.ast, got)
	})

	t.Run("Fail_InvalidASTBuilder", func(t *testing.T) {
		_, err := getAST("not an astBuilder")
		assert.Error(t, err)
	})
}
