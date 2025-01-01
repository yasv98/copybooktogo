package parse

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testData struct {
	level       int
	identifier  string
	redefines   string
	occursCount int
	pic         Picture
}

func newTestData() testData {
	return testData{
		level:       1,
		identifier:  "Record-1",
		redefines:   "Record-X",
		occursCount: 5,
		pic:         Picture{PicString: "X(10)", PicType: Alpha, PicCount: 10},
	}
}

func Test_createRecord(t *testing.T) {
	td := newTestData()

	t.Run("Success", func(t *testing.T) {
		result, err := createRecord(td.level, td.identifier, []any{td.redefines, td.pic, td.occursCount})
		require.NoError(t, err)
		assert.Equal(t, Record{
			Level:       td.level,
			Identifier:  td.identifier,
			Redefines:   td.redefines,
			Pic:         td.pic,
			OccursCount: td.occursCount,
		}, result)
	})

	t.Run("Fail_IncorrectLevelType", func(t *testing.T) {
		result, err := createRecord("invalid type", td.identifier, []any{})
		assert.Error(t, err)
		assert.Empty(t, result)
	})

	t.Run("Fail_IncorrectIdentifierType", func(t *testing.T) {
		result, err := createRecord(td.level, -1, []any{})
		assert.Error(t, err)
		assert.Empty(t, result)
	})

	t.Run("Fail_IncorrectClausesType", func(t *testing.T) {
		result, err := createRecord(td.level, td.identifier, "invalid type")
		assert.Error(t, err)
		assert.Empty(t, result)
	})
}

func Test_processClause(t *testing.T) {
	td := newTestData()

	tests := map[string]struct {
		record   Record
		clause   any
		expected Record
		wantErr  bool
	}{
		"Success_RedefinesClause": {
			record:   Record{},
			clause:   td.redefines,
			expected: Record{Redefines: td.redefines},
			wantErr:  false,
		},
		"Success_PictureClause": {
			record:   Record{},
			clause:   td.pic,
			expected: Record{Pic: td.pic},
			wantErr:  false,
		},
		"Success_OccursClause": {
			record:   Record{},
			clause:   td.occursCount,
			expected: Record{OccursCount: td.occursCount},
			wantErr:  false,
		},
		"Fail_RedefinesAlreadySet": {
			record:   Record{Redefines: td.redefines},
			clause:   td.redefines,
			expected: Record{Redefines: td.redefines},
			wantErr:  true,
		},
		"Fail_PictureAlreadySet": {
			record:   Record{Pic: td.pic},
			clause:   td.pic,
			expected: Record{Pic: td.pic},
			wantErr:  true,
		},
		"Fail_OccursAlreadySet": {
			record:   Record{OccursCount: td.occursCount},
			clause:   td.occursCount,
			expected: Record{OccursCount: td.occursCount},
			wantErr:  true,
		},
		"Fail_UnexpectedClauseType": {
			record:   Record{},
			clause:   12.5,
			expected: Record{},
			wantErr:  true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			err := tt.record.processClause(tt.clause)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.expected, tt.record)
		})
	}
}

func Test_getClauseDetails(t *testing.T) {
	td := newTestData()

	t.Run("Success", func(t *testing.T) {
		t.Run("RedefinesClause", func(t *testing.T) {
			result, err := getRedefinesClauseDetails(td.redefines)
			require.NoError(t, err)
			assert.Equal(t, td.redefines, result)
		})
		t.Run("PictureClause", func(t *testing.T) {
			result, err := getPictureClauseDetails("X(10)")
			require.NoError(t, err)
			assert.Equal(t, td.pic, result)
		})
		t.Run("OccursClause", func(t *testing.T) {
			result, err := getOccursClauseDetails(td.occursCount)
			require.NoError(t, err)
			assert.Equal(t, td.occursCount, result)
		})
	})

	t.Run("Fail_IncorrectType", func(t *testing.T) {
		t.Run("RedefinesClause", func(t *testing.T) {
			result, err := getRedefinesClauseDetails(-1)
			assert.Error(t, err)
			assert.Empty(t, result)
		})
		t.Run("PictureClause", func(t *testing.T) {
			result, err := getPictureClauseDetails(-1)
			assert.Error(t, err)
			assert.Empty(t, result)
		})
		t.Run("OccursClause", func(t *testing.T) {
			result, err := getOccursClauseDetails("invalid type")
			assert.Error(t, err)
			assert.Empty(t, result)
		})
	})
}

func Test_parseIntFromBytes(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		result, err := parseIntFromBytes([]byte("123"))
		assert.NoError(t, err)
		assert.Equal(t, 123, result)
	})

	t.Run("Fail_InvalidInput", func(t *testing.T) {
		result, err := parseIntFromBytes("invalid type")
		assert.Error(t, err)
		assert.Zero(t, result)
	})
}
