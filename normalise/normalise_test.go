package normalise

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormaliseCopybook(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected []byte
		wantErr  bool
	}{
		{
			name:     "Empty input",
			input:    []byte{},
			expected: []byte{},
			wantErr:  false,
		},
		{
			name: "Simple copybook",
			input: []byte(` 01  RECORD-1.
     05  FIELD-A    PIC X(10).
     05  FIELD-B    PIC 9(5).`),
			expected: []byte(`       01  RECORD-1.                                                    
           05  FIELD-A    PIC X(10).                                    
           05  FIELD-B    PIC 9(5).                                     `),
			wantErr: false,
		},
		{
			name: "Format with sequence numbers and extra characters",
			input: []byte(`000100 01  RECORD-1.                                                       extra
000200     05  FIELD-A    PIC X(10).                                       extra
000300     05  FIELD-B    PIC 9(5).                                         extra`),
			expected: []byte(`       01  RECORD-1.                                                    
           05  FIELD-A    PIC X(10).                                    
           05  FIELD-B    PIC 9(5).                                     `),
			wantErr: false,
		},
		{
			name: "Format sequence numbers and indicators",
			input: []byte(`000100/01  RECORD-1.                                                       
000200*    05  FIELD-A    PIC X(10).                                       
000300d    05  FIELD-B    PIC 9(5).                                         `),
			expected: []byte(`      /01  RECORD-1.                                                    
      *    05  FIELD-A    PIC X(10).                                    
      d    05  FIELD-B    PIC 9(5).                                     `),
			wantErr: false,
		},
		{
			name: "Preformatted copybook",
			input: []byte(`       01  RECORD-1.                                                    
           05  FIELD-A    PIC X(10).                                    
           05  FIELD-B    PIC 9(5).                                     `),
			expected: []byte(`       01  RECORD-1.                                                    
           05  FIELD-A    PIC X(10).                                    
           05  FIELD-B    PIC 9(5).                                     `),
			wantErr: false,
		},
		{
			name:     "Format right space unicode normalised",
			input:    []byte("       01  RECORD-1.\r"),
			expected: []byte("       01  RECORD-1.                                                    "),
			wantErr:  false,
		},
		{
			name: "Format with no 01 level",
			input: []byte(`000100 05  RECORD-1.
000200     05  FIELD-A    PIC X(10).
000300     05  FIELD-B    PIC 9(5).`),
			expected: nil,
			wantErr:  true,
		},
		{
			name: "Format with indicator area cut off",
			input: []byte(`01  RECORD-1.
    05  FIELD-A    PIC X(10).
    05  FIELD-B    PIC 9(5).`),
			expected: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Format(tt.input)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func Test_getDataBlockIndentation(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected int
		wantErr  bool
	}{
		{
			name:     "Data block starts at column 0",
			input:    []string{" 01 RECORD."},
			expected: 0,
			wantErr:  false,
		},
		{
			name:     "Data block starts at column 7 with sequence number",
			input:    []string{"012345 01 RECORD."},
			expected: 6,
			wantErr:  false,
		},
		{
			name:     "Data block starts at column 7 with character sequence number",
			input:    []string{"012345 01 RECORD."},
			expected: 6,
			wantErr:  false,
		},
		{
			name:     "Data block starts at column 7 with first level in column 10 of Area A",
			input:    []string{"012345   01 RECORD."},
			expected: 6,
			wantErr:  false,
		},
		{
			name:     "Data block starts at column 7 with indicator",
			input:    []string{"012345d01 RECORD."},
			expected: 6,
			wantErr:  false,
		},
		{
			name:     "Line has incomplete sequence number",
			input:    []string{"2345 01 RECORD."},
			expected: 0,
			wantErr:  true,
		},
		{
			name:     "No indicator area",
			input:    []string{"01 RECORD."},
			expected: 0,
			wantErr:  true,
		},
		{
			name:     "No 01 found",
			input:    []string{"05 FIELD PIC X(10)."},
			expected: 0,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getDataBlockIndentation(tt.input)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.expected, got)
		})
	}
}
