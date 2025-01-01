package parse

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"copybooktogo/util/xrequire"
)

func Test_ParseNormalisedCopybook(t *testing.T) {
	tests := map[string]struct {
		input       []byte
		expected    []*Record
		assertError assert.ErrorAssertionFunc
	}{
		"ValidCopyBook_ReturnsParsedAST": {
			input: []byte(`      ******************************************************************
      *
       01  RECORD-1.                                                    
           03  FILLER              PIC X(31).                           
           03  RECORD-2                                                 
                                   PIC X(01).                           
               88  RECORD-3                        VALUE 'S'.           
               88  RECORD-4                         VALUE 'P'.           
           03  RECORD-5.                                                
               05  RECORD-6                        OCCURS 10 TIMES      
                                                   INDEXED BY           
                                                   X-:XXXX:-DBT.        
                   07  RECORD-7                                         
                                   PIC S9(07)      COMP-3.              
           03  FILLER              PIC X(190).                          
       01  RECORD-8.                                                    
           05  RECORD-9                           PIC  X(02).           
           05  RECORD-10       REDEFINES                                
               RECORD-9        PIC  X(02).                              
`),
			expected: []*Record{
				{
					Level:      1,
					Identifier: "RECORD-1",
					Children: []*Record{
						{
							Level:      3,
							Identifier: "FILLER",
							Pic:        Picture{PicString: "X(31)", PicType: Alpha, PicCount: 31},
						},
						{
							Level:      3,
							Identifier: "RECORD-2",
							Pic:        Picture{PicString: "X(01)", PicType: Alpha, PicCount: 1},
						},
						{
							Level:      3,
							Identifier: "RECORD-5",
							Children: []*Record{
								{
									Level:       5,
									Identifier:  "RECORD-6",
									OccursCount: 10,
									Children: []*Record{
										{
											Level:      7,
											Identifier: "RECORD-7",
											Pic:        Picture{PicString: "S9(07)", PicType: Signed, PicCount: 8},
										},
									},
								},
							},
						},
						{
							Level:      3,
							Identifier: "FILLER",
							Pic:        Picture{PicString: "X(190)", PicType: Alpha, PicCount: 190},
						},
					},
				},
				{
					Level:      1,
					Identifier: "RECORD-8",
					Children: []*Record{
						{
							Level:      5,
							Identifier: "RECORD-9",
							Pic:        Picture{PicString: "X(02)", PicType: Alpha, PicCount: 2},
						},
						{
							Level:      5,
							Identifier: "RECORD-10",
							Pic:        Picture{PicString: "X(02)", PicType: Alpha, PicCount: 2},
							Redefines:  "RECORD-9",
						},
					},
				},
			},
			assertError: assert.NoError,
		},
		"InvalidCopybookWithNo01Record_ReturnsError": {
			input: []byte(`       05  RECORD-1.                                                    
           10  FILLER              PIC X(31).                           
           10  RECORD-2                                                 
`),
			assertError: assert.Error,
		},
		"InvalidCopybookWithInvalidRecordLevel_ReturnsError": {
			input: []byte(`       00  RECORD-1.                                                    
           10  FILLER              PIC X(31).                           
           10  RECORD-2                                                 
`),
			assertError: assert.Error,
		},
	}

	for name, test := range tests {
		tt := test
		t.Run(name, func(t *testing.T) {
			got, err := BuildAST(tt.input)
			tt.assertError(t, err)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func Test_ParseIndividualRecord(t *testing.T) {
	t.Run("Level01Record", func(t *testing.T) {
		input := []byte("       01  RECORD-1.                                                    ")
		expected := []*Record{{Level: 1, Identifier: "RECORD-1"}}
		got, err := BuildAST(input)
		require.NoError(t, err)
		assert.Equal(t, expected, got)
	})

	t.Run("Level01RecordWithRedefines", func(t *testing.T) {
		input := []byte(`       01  RECORD-2                                                     
                                   REDEFINES RECORD-1.                  `)
		expected := []*Record{{Level: 1, Identifier: "RECORD-2", Redefines: "RECORD-1"}}
		got, err := BuildAST(input)
		require.NoError(t, err)
		assert.Equal(t, expected, got)
	})

	tests := map[string]struct {
		input    []byte
		expected []*Record
	}{
		"PIC single line": {
			input: []byte("           05  RECORD                             PIC  X(505).          "),
			expected: []*Record{
				{
					Level:      5,
					Identifier: "RECORD",
					Pic:        Picture{PicString: "X(505)", PicType: Alpha, PicCount: 505},
				},
			},
		},
		"PIC single line with dot as decimal and trailing sign": {
			input: []byte("           10  RECORD                          PIC  9(03).9(4)-.        "),
			expected: []*Record{
				{
					Level:      10,
					Identifier: "RECORD",
					Pic:        Picture{PicString: "9(03).9(4)-", PicType: Decimal, PicCount: 9},
				},
			},
		},
		"PIC with COMP": {
			input: []byte(`               05  RECORD          PIC S9(09)      COMP-3.              
`),
			expected: []*Record{
				{
					Level:      5,
					Identifier: "RECORD",
					Pic:        Picture{PicString: "S9(09)", PicType: Signed, PicCount: 10},
				},
			},
		},
		"PIC multi-line": {
			input: []byte(`                 07 RECORD                                              
                                       PIC X(10).                       
`),
			expected: []*Record{
				{
					Level:      7,
					Identifier: "RECORD",
					Pic:        Picture{PicString: "X(10)", PicType: Alpha, PicCount: 10},
				},
			},
		},
		"PIC with JUSTIFIED RIGHT": {
			input: []byte(`         10  RECORD                           PIC X(15)                
                                                 JUSTIFIED RIGHT.                                 
`),
			expected: []*Record{
				{
					Level:      10,
					Identifier: "RECORD",
					Pic:        Picture{PicString: "X(15)", PicType: Alpha, PicCount: 15},
				},
			},
		},
		"REDEFINES single line": {
			input: []byte("           05  RECORD               REDEFINES RECORD.           "),
			expected: []*Record{
				{
					Level:      5,
					Identifier: "RECORD",
					Redefines:  "RECORD",
				},
			},
		},
		"REDEFINES multi-line": {
			input: []byte(`           03  RECORD                                                   
                                   REDEFINES RECORD-2.                  
`),
			expected: []*Record{
				{
					Level:      3,
					Identifier: "RECORD",
					Redefines:  "RECORD-2",
				},
			},
		},
		"REDEFINES with PIC": {
			input: []byte(`                 07 RECORD                                              
                                   REDEFINES RECORD-2                   
                                       PIC S9(13) COMP-3.               
`),
			expected: []*Record{
				{
					Level:      7,
					Identifier: "RECORD",
					Redefines:  "RECORD-2",
					Pic:        Picture{PicString: "S9(13)", PicType: Signed, PicCount: 14},
				},
			},
		},
		"OCCURS single line": {
			input: []byte("       15  RECORD               OCCURS 4.                          "),
			expected: []*Record{
				{
					Level:       15,
					Identifier:  "RECORD",
					OccursCount: 4,
				},
			},
		},
		"OCCURS single line with PIC": {
			input: []byte("       15  RECORD               PIC X(40) OCCURS 4.                "),
			expected: []*Record{
				{
					Level:       15,
					Identifier:  "RECORD",
					Pic:         Picture{PicString: "X(40)", PicType: Alpha, PicCount: 40},
					OccursCount: 4,
				},
			},
		},
		"OCCURS with times": {
			input: []byte("       15  RECORD               PIC X(40) OCCURS 4 TIMES.          "),
			expected: []*Record{
				{
					Level:       15,
					Identifier:  "RECORD",
					Pic:         Picture{PicString: "X(40)", PicType: Alpha, PicCount: 40},
					OccursCount: 4,
				},
			},
		},
		"OCCURS multi-line": {
			input: []byte(`               10  RECORD                         PIC  X(40)            
                                                  OCCURS 4 TIMES.       
`),
			expected: []*Record{
				{
					Level:       10,
					Identifier:  "RECORD",
					Pic:         Picture{PicString: "X(40)", PicType: Alpha, PicCount: 40},
					OccursCount: 4,
				},
			},
		},
		"OCCURS with INDEXED BY": {
			input: []byte(`               05  RECORD-6                        OCCURS 10 TIMES      
                                                   INDEXED BY           
                                                   X-:XXXX:-DBT.        
`),
			expected: []*Record{
				{
					Level:       5,
					Identifier:  "RECORD-6",
					OccursCount: 10,
				},
			},
		},
	}

	for name, test := range tests {
		tt := test
		t.Run(name, func(t *testing.T) {
			// Add a dummy Record to the input so individual records that are not top Level records
			// can be parsed without error.
			createTestRecord := func(individual []byte) []byte {
				return slices.Concat([]byte("       01  DUMMY-RECORD.                                                \n"), individual)
			}
			got, err := BuildAST(createTestRecord(tt.input))
			require.NoError(t, err)
			root := xrequire.Single(t, got)
			assert.Equal(t, "DUMMY-RECORD", root.Identifier)
			assert.Equal(t, tt.expected, root.Children)
		})
	}
}

func Test_ParseDataIgnores(t *testing.T) {
	tests := map[string]struct {
		input []byte
	}{
		"BlankLine": {
			input: []byte("                                                                        "),
		},
		"CommentLine": {
			input: []byte("      ******************************************************************        "),
		},
		"UnknownLine": {
			input: []byte("       U    N    K    O    W    N                                       "),
		},
	}

	for name, test := range tests {
		tt := test
		t.Run(name, func(t *testing.T) {
			got, err := BuildAST(tt.input)
			require.NoError(t, err)
			assert.Nil(t, got)
		})
	}
}
