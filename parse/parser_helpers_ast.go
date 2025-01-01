package parse

import (
	"fmt"
)

// parser_helpers_ast defines functionality to build an Abstract Syntax Tree (AST)
// for parsing COBOL copybooks.
//
// The package uses an astBuilder struct which consists of two main components:
//   - ast: a slice holding the top Level 1 records
//   - workingParentsStack: an operational stack used to track parent records
//
// AST Building Process:
//
// When adding a Record to the AST:
//  1. For Level 1 records:
//     - Add to the ast slice
//     - Re-initialize working parents stack with this new Level 1 Record
//  2. For levels > 1:
//     - Pop records from the working parents stack until the parent is at the top
//     - Append the current Record to its parent (the Record at the top of the working parents stack)
//     - If the Record has Children (non-leaf node), append it to working parents stack
//
// This process repeats until all records are processed, after which the ast slice will hold the parsed AST.
//
// Example:
//
// Given this COBOL copybook:
//
//	01  EXAMPLE.
//	   05  DUMMY-GROUP-1.
//	      10  DUMMY-SUB-GROUP-1.
//	         15  DUMMY-GROUP-1-OBJECT-A   PIC 9(4).
//	   05  DUMMY-GROUP-2                    PIC X(201).
//
// The parsed output would be:
//
//	Level: 1, Identifier: "EXAMPLE", Children: {
//	   Level: 5, Identifier: "DUMMY-GROUP-1", Children: {
//	      Level: 10, Identifier: "DUMMY-SUB-GROUP-1", Children: {
//	         Level: 15, Identifier: "DUMMY-GROUP-1-OBJECT-A", Pic: { PicType: "uint", PicCount: 4 }
//	      },
//	   },
//	   Level: 5, Identifier: "DUMMY-GROUP-2", Pic: { PicType: "string", PicCount: 201 }
//	}
type astBuilder struct {
	ast                 []*Record
	workingParentsStack workingParentsStack
}

type workingParentsStack []*Record

func createAndAddRecordToAST(ast, level, identifier, clauses any) error {
	treeBuilder, ok := ast.(*astBuilder)
	if !ok {
		return fmt.Errorf("ast is not a *astBuilder: %v", ast)
	}

	newRecord, err := createRecord(level, identifier, clauses)
	if err != nil {
		return fmt.Errorf("failed to create Record: %w", err)
	}

	if err := treeBuilder.addRecord(&newRecord); err != nil {
		return fmt.Errorf("failed to add Record to AST: %w", err)
	}

	return nil
}

func (ab *astBuilder) addRecord(rec *Record) error {
	if rec.Level < 1 {
		return fmt.Errorf("Record Level cannot be less than 1: %v", rec.Level)
	}

	// If a Level 1 Record, append to ast and reset working parents stack to this Record.
	if rec.Level == 1 {
		ab.ast = append(ab.ast, rec)
		ab.workingParentsStack = []*Record{rec}
		return nil
	}

	// Pop records from the working parents stack until the parent is at the top.
	for !ab.workingParentsStack.isTopParent(rec) {
		if len(ab.workingParentsStack) == 0 {
			return fmt.Errorf("stack is empty, should have at least one Level 1 Record")
		}
		ab.workingParentsStack.pop()
	}

	// Append Record to its parent.
	parent := ab.workingParentsStack.peek()
	parent.Children = append(parent.Children, rec)

	// If the Record is not a leaf node, append to the working parents stack.
	if !isLeafNode(rec) {
		ab.workingParentsStack.append(rec)
	}

	return nil
}

func isLeafNode(rec *Record) bool {
	// A Record with a Picture clause is a leaf node and will not have Children.
	return rec.Pic != Picture{}
}

func getAST(ast any) ([]*Record, error) {
	treeBuilder, ok := ast.(*astBuilder)
	if !ok {
		return nil, fmt.Errorf("ast is not a *astBuilder: %v", ast)
	}

	return treeBuilder.ast, nil
}

func (wp *workingParentsStack) isTopParent(rec *Record) bool {
	return len(*wp) > 0 && wp.peek().Level < rec.Level
}

func (wp *workingParentsStack) pop() {
	*wp = (*wp)[:len(*wp)-1]
}

func (wp *workingParentsStack) append(rec *Record) {
	*wp = append(*wp, rec)
}

func (wp *workingParentsStack) peek() *Record {
	return (*wp)[len(*wp)-1]
}
