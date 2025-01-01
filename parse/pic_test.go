package parse

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_parsePICType(t *testing.T) {
	tests := map[string]struct {
		Input    string
		Expected PicType
	}{
		"Alpha":                    {"X", Alpha},
		"Alpha with count":         {"X(5)", Alpha},
		"Alphanumeric":             {"A(10)", Alpha},
		"Decimal":                  {"9V99", Decimal},
		"Decimal period":           {"9.99", Decimal},
		"Decimal period with sign": {"9(03).9(4)-", Decimal},
		"Signed integer":           {"S9(5)", Signed},
		"Unsigned integer":         {"9(9)", Unsigned},
		"Complex decimal":          {"S9(5)V99", Decimal},
		"Unknown":                  {"?", Unknown},
		"Empty string":             {"", Unknown},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			result := parsePICType(tt.Input)
			assert.Equal(t, tt.Expected, result)
		})
	}
}

func Test_parsePICCount(t *testing.T) {
	tests := map[string]struct {
		Input    string
		Expected int
	}{
		"Single character":         {"X", 1},
		"Multiple characters":      {"XXX", 3},
		"Single count":             {"X(5)", 5},
		"Multiple counts":          {"X(2)9(3)", 5},
		"Implied decimal":          {"9(15)V99", 17},
		"Scaling factor start":     {"PPP9", 1},
		"Scaling factor repeated":  {"P(3)9", 1},
		"Scaling factor end":       {"9PPP", 1},
		"Mixed format":             {"S9(5)V99", 8},
		"Complex format":           {"S9(5)V9(2)", 8},
		"Complex format with sign": {"9(03).9(4)-", 9},
		"Invalid format":           {"X(A)", -1},
		"Empty string":             {"", 0},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			result := parsePICCount(tt.Input)
			assert.Equal(t, tt.Expected, result)
		})
	}
}
