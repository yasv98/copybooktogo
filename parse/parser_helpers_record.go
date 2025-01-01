// Package parse defines functionality to parse COBOL copybooks.
package parse

import (
	"fmt"
	"strconv"
)

//go:generate pigeon -o parser.generated.go -optimize-parser -optimize-basic-latin copybook.peg

// Record defines a single record in a COBOL copybook.
type Record struct {
	Level       int
	Identifier  string
	Redefines   string
	Pic         Picture
	OccursCount int
	Children    []*Record
}

// Picture defines the PIC clause details for a record.
type Picture struct {
	PicString string
	PicType   PicType
	PicCount  int
}

func createRecord(level, identifier, clauses any) (Record, error) {
	levelInt, ok := level.(int)
	if !ok {
		return Record{}, fmt.Errorf("level is not an int: %v", level)
	}

	identifierString, ok := identifier.(string)
	if !ok {
		return Record{}, fmt.Errorf("identifier is not a string: %v", identifier)
	}

	newRecord := Record{
		Level:      levelInt,
		Identifier: identifierString,
	}

	clauseSlice, ok := clauses.([]any)
	if !ok {
		return Record{}, fmt.Errorf("clauses is not a []any: %v", clauses)
	}
	if len(clauseSlice) > 0 {
		// Process each clause and add to the Record's clauseDetails.
		for _, clause := range clauseSlice {
			if err := newRecord.processClause(clause); err != nil {
				return Record{}, fmt.Errorf("failed to process clause: %w", err)
			}
		}
	}

	return newRecord, nil
}

func (r *Record) processClause(clause any) error {
	switch typedClause := clause.(type) {
	case string:
		if r.Redefines != "" {
			return fmt.Errorf("redefines clause already set: %v", r.Redefines)
		}
		r.Redefines = typedClause
	case Picture:
		if r.Pic != (Picture{}) {
			return fmt.Errorf("picture clause already set: %v", r.Pic)
		}
		r.Pic = typedClause
	case int:
		if r.OccursCount != 0 {
			return fmt.Errorf("occurs clause already set: %v", r.OccursCount)
		}
		r.OccursCount = typedClause
	default:
		return fmt.Errorf("unexpected clause type: %T", clause)
	}
	return nil
}

func getRedefinesClauseDetails(identifier any) (string, error) {
	identifierString, ok := identifier.(string)
	if !ok {
		return "", fmt.Errorf("identifier is not a string: %v", identifier)
	}

	return identifierString, nil
}

func getPictureClauseDetails(pic any) (Picture, error) {
	picString, ok := pic.(string)
	if !ok {
		return Picture{}, fmt.Errorf("pic is not a string: %v", pic)
	}

	return Picture{
		PicString: picString,
		PicType:   parsePICType(picString),
		PicCount:  parsePICCount(picString),
	}, nil
}

func getOccursClauseDetails(count any) (int, error) {
	countInt, ok := count.(int)
	if !ok {
		return 0, fmt.Errorf("count is not an int: %v", count)
	}

	return countInt, nil
}

func parseIntFromBytes(value any) (int, error) {
	valueSlice, ok := value.([]byte)
	if !ok {
		return 0, fmt.Errorf("value is not a byte slice: %v", value)
	}

	num, err := strconv.Atoi(string(valueSlice))
	if err != nil {
		return 0, fmt.Errorf("failed to parse integer: %v", err)
	}

	return num, nil
}
