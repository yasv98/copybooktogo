package parse

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// PicType defines the different types of PIC definitions.
//
//go:generate enumer -type PicType -output "pictype_enumer.generated.go" -linecomment
type PicType int

const (
	// Unknown represents an unknown PIC type.
	Unknown PicType = iota // unknown
	// Unsigned represents an unsigned number (e.g. 9(5)).
	Unsigned // unsigned
	// Signed represents a signed number (e.g. S9(5)).
	Signed // signed
	// Decimal represents a decimal number (e.g. 9(5)V99).
	Decimal // decimal
	// Alpha represents an alphanumeric string (e.g. X(5)).
	Alpha // alpha

	alphaIndicators     = "XA"
	decimalIndicators   = ".VP"
	signedIntIndicators = "S"
	intIndicators       = "9"
)

// zeroWidthIndicatorRegex matches field type indicators that do not
// contribute to width.
var zeroWidthIndicatorRegex = regexp.MustCompile(`V|P(?:\(\d+\))?`)

func PicTypeFromString(s string) (PicType, error) {
	picType, err := PicTypeString(strings.ToLower(s))
	if err != nil {
		return 0, fmt.Errorf("%q must be a valid PIC type: %v", picType, PicTypeValues())
	}
	return picType, nil
}

// parsePICType identifies an equivalent Go type from the given substring
// that contains a PIC definition.
func parsePICType(s string) PicType {
	if strings.ContainsAny(s, alphaIndicators) {
		return Alpha
	}

	if strings.ContainsAny(s, decimalIndicators) {
		return Decimal
	}

	if strings.ContainsAny(s, signedIntIndicators) {
		return Signed
	}

	if strings.ContainsAny(s, intIndicators) {
		return Unsigned
	}

	return Unknown
}

// parsePICCount identifies the fixed width, or length, of the given
// PIC definition such as: X(2), XX, 9(9), etc.
// For example:
// S9(5)V9(7): "S" = 1, "9(5)" = 5, "V" = 0, "9(7)" = 7 => 19
// S9(5).9(7): "S" = 1, "9(5)" = 5, "." = 1, "9(7)" = 7 => 20
// PPP9(5): "PPP" = 0, "9(5)" = 5 => 5
func parsePICCount(s string) int {
	// Remove indicators that do not contribute to the width.
	s = zeroWidthIndicatorRegex.ReplaceAllString(s, "")

	size := 0
	for strings.Contains(s, "(") {
		left := strings.Index(s, "(")
		right := strings.Index(s, ")")
		// capture type when using parentheses "9(99)" should be stripped to
		// "" so that it results in 99+0, not 99+len("9").
		start := left - 1
		end := right + 1
		amount, err := strconv.Atoi(s[left+1 : right])
		if err != nil {
			return -1
		}

		size += amount
		s = s[:start] + s[end:]
	}

	return size + len(s)
}
