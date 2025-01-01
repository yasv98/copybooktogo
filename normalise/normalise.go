// Package normalise provides functionality to normalise COBOL copybooks to the standard IBM reference format.
package normalise

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"copybooktogo/util/generic"
)

const (
	indicatorArea = 7
	areaBEnd      = 72
	dataBlockLen  = areaBEnd - (indicatorArea - 1)
)

// Format formats a copybook in memory to the standard IBM reference format
// without sequence numbers and extra characters.
// https://www.ibm.com/docs/en/cobol-zos/6.4?topic=structure-reference-format
//
// NOTE: This assumes the copybook is either already in the standard IBM reference format,
// or its sequence numbers have been removed, and it starts from the indicator area.
func Format(copybook []byte) ([]byte, error) {
	if len(copybook) == 0 {
		return []byte{}, nil
	}

	lines := strings.Split(string(copybook), "\n")
	indentation, err := getDataBlockIndentation(lines)
	if err != nil {
		return nil, err
	}

	lines = generic.Map(func(line string) string {
		return normaliseLine(line, indentation)
	}, lines)
	return []byte(strings.Join(lines, "\n")), nil
}

// 01 level can be preceded by a sequence number, spaces or valid indicator area inputs.
var recordDescriptionEntryRegex = regexp.MustCompile(`^(?P<Indentation>(?P<SequenceNumberArea>\s*.{6})?)(?P<IndicatorArea>[/Dd\s])(?P<OptionalIndentation>\s*)(?:01\s)`)

func getDataBlockIndentation(lines []string) (int, error) {
	for _, line := range lines {
		if groups, matched := findMatchGroups(recordDescriptionEntryRegex, line); matched {
			return len(groups["Indentation"]), nil
		}
	}
	return 0, fmt.Errorf("first level 01 not found")
}

func normaliseLine(line string, indentation int) string {
	line = strings.TrimRightFunc(line, unicode.IsSpace)
	start, end := indentation, dataBlockLen+indentation
	if len(line) < end {
		line += strings.Repeat(" ", end-len(line))
	}
	dataBlock := line[start:end]
	return strings.Repeat(" ", indicatorArea-1) + dataBlock
}

func findMatchGroups(re *regexp.Regexp, s string) (map[string]string, bool) {
	getNamedMatches := func(re *regexp.Regexp, matches []string) map[string]string {
		result := make(map[string]string)
		for i, name := range re.SubexpNames() {
			if i < len(matches) {
				result[name] = matches[i]
			}
		}
		return result
	}
	matches := re.FindStringSubmatch(s)
	return getNamedMatches(re, matches), len(matches) > 0
}
