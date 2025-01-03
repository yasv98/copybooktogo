// Package generate provides functionality for generating Go struct definitions from COBOL copybook ASTs.
package generate

import (
	"bytes"
	"fmt"
	"slices"
	"text/template"

	"copybooktogo/util/generic"

	"github.com/kenshaw/snaker"
	"golang.org/x/tools/imports"

	"copybooktogo/parse"
)

// TODO: Add version info to the generated attribute once tool is out of development.
const goStructsGenTemplate = `
// This file is generated by copybooktogo. DO NOT EDIT.

package {{.Package}}

import (
    "github.com/anzx/fabric-go-pic/pkg/decimal"
)

{{ range .Structs }}
// {{ .StructVarName }} contains a representation of {{ .Identifier }}
type {{ .StructVarName }} struct {
    {{- range .Fields }}
    {{ .FieldVarName }} {{ .VarType }} ` + "`pic:\"{{ .PicTag }}\"`" + ` // start:{{ .PicGlobalStart }} end:{{ .PicGlobalEnd }}{{if .RedefinesVarName}} REDEFINES {{ .RedefinesVarName }}{{end}}
    {{- end }}
}
{{ end }}
`

type templateParams struct {
	Package string
	Structs []StructData
}

// StructData represents a Go struct definition.
type StructData struct {
	StructVarName string
	Identifier    string
	Fields        []FieldData
}

// FieldData represents a field in a Go struct.
type FieldData struct {
	FieldVarName     string
	VarType          string
	RedefinesVarName string
	PicSize          int
	PicTag           string
	PicGlobalStart   int
	PicGlobalEnd     int
}

type goGenerator struct {
	pos            *positionTracker
	picTypeMapping map[parse.PicType]string
}

type positionInfo struct {
	localStart  int
	globalStart int
}

type positionTracker struct {
	localPos    int
	globalPos   int
	recordStore map[string]positionInfo
}

// ToGoStructsData generates Go struct definitions from a COBOL copybook AST.
func ToGoStructsData(ast []*parse.Record, copybookName, packageName string, typeOverrides map[parse.PicType]string) ([]byte, error) {
	if len(ast) == 0 {
		return nil, fmt.Errorf("ast is empty")
	}

	goGen := goGenerator{
		pos: newPositionTracker(),
		// Merge default PIC type mappings with any configured overrides.
		picTypeMapping: generic.MergeMaps(defaultTypeMapping(), typeOverrides),
	}

	data := templateParams{
		Package: packageName,
		Structs: goGen.buildStructData(copybookName, ast),
	}

	generatedCode, err := executeTemplate(goStructsGenTemplate, data)
	if err != nil {
		return nil, err
	}

	return imports.Process("", generatedCode, nil)
}

func executeTemplate(genTemplate string, data templateParams) ([]byte, error) {
	t, err := template.New("").Parse(genTemplate)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}
	var buf bytes.Buffer
	if err = t.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.Bytes(), nil
}

func (g *goGenerator) buildStructData(parentName string, records []*parse.Record) []StructData {
	currentStruct := StructData{
		StructVarName: toGoName(parentName),
		Identifier:    parentName,
		Fields:        g.buildFieldsData(records, parentName),
	}

	// Recursively process nested struct fields.
	var nestedStructs []StructData
	for _, field := range records {
		if len(field.Children) > 0 {
			// A field's children will start from the same global position as the parent field.
			_, g.pos.globalPos = g.pos.getStoredPos(field.Identifier)
			nestedStructs = slices.Concat(nestedStructs, g.buildStructData(field.Identifier, field.Children))
		}
	}

	return slices.Concat([]StructData{currentStruct}, nestedStructs)
}

func (g *goGenerator) buildFieldsData(records []*parse.Record, parentName string) []FieldData {
	fields := make([]FieldData, 0, len(records))
	fillerCount := 0
	g.pos.localPos = 1

	for _, rec := range records {
		fillerCount = handleFillerName(rec, parentName, fillerCount)

		// Build and store field data
		fieldData := g.buildFieldData(rec)
		fieldData = handleRedefines(rec, fieldData, g.pos)
		fields = append(fields, fieldData)

		// Update position tracking
		g.pos.storeAndAdvancePos(rec.Identifier, fieldData.PicSize)
	}

	return fields
}

func (g *goGenerator) buildFieldData(rec *parse.Record) FieldData {
	varName := toGoName(rec.Identifier)
	size := calculateSize(rec)

	return FieldData{
		FieldVarName:   varName,
		VarType:        getVarType(rec, varName, g.picTypeMapping),
		PicSize:        size,
		PicTag:         getPicTag(rec, size, g.pos.localPos),
		PicGlobalStart: g.pos.globalPos,
		PicGlobalEnd:   g.pos.globalPos + size - 1,
	}
}

func handleRedefines(rec *parse.Record, fieldData FieldData, pos *positionTracker) FieldData {
	if rec.Redefines != "" {
		// If a field redefines another field, its local and global
		// start position will be the same as the redefined field.
		pos.localPos, pos.globalPos = pos.getStoredPos(rec.Redefines)

		fieldData.RedefinesVarName = toGoName(rec.Redefines)
		fieldData.PicTag = getPicTag(rec, fieldData.PicSize, pos.localPos)
		fieldData.PicGlobalStart = pos.globalPos
		fieldData.PicGlobalEnd = pos.globalPos + fieldData.PicSize - 1
	}

	return fieldData
}

func handleFillerName(rec *parse.Record, parentName string, fillerCount int) int {
	// FILLER is a keyword that can be used multiple times at a group and individual
	// level in COBOL. This logic ensure there is no clashes with struct and field
	// names when generating the Go code.
	const fillerKeyWord = "FILLER"
	if rec.Identifier == fillerKeyWord {
		fillerCount++
		rec.Identifier = fmt.Sprint(parentName, "-", fillerKeyWord, fillerCount)
	}

	return fillerCount
}

func toGoName(s string) string {
	return snaker.SnakeToCamelIdentifier(s)
}

func getVarType(rec *parse.Record, varName string, picTypeMappings map[parse.PicType]string) string {
	switch {
	case len(rec.Children) == 0:
		goType, ok := picTypeMappings[rec.Pic.PicType]
		if !ok {
			// Default to string if no mapping is found.
			goType = "string"
		}

		if rec.OccursCount > 1 {
			return fmt.Sprint("[", rec.OccursCount, "]", goType)
		}
		return goType
	case rec.OccursCount > 1:
		return fmt.Sprint("[", rec.OccursCount, "]", varName)
	default:
		return varName
	}
}

func getPicTag(rec *parse.Record, fieldSize, localStartPos int) string {
	picTag := fmt.Sprint(localStartPos, ",", localStartPos+fieldSize-1)
	if rec.OccursCount > 1 {
		picTag += fmt.Sprint(",", rec.OccursCount)
	}

	if len(rec.Children) == 0 {
		picTag += fmt.Sprint(",clause=", rec.Pic.PicString)
	} else {
		// To account for group fields with occurs, we need to calculate the size of the group for
		// one occurrence.
		singleFieldSize := fieldSize / max(1, rec.OccursCount)
		// Pad single digit clause values with a leading zero to match common COBOL convention.
		picTag += fmt.Sprintf(",clause=X(%02d)", singleFieldSize)
	}

	return picTag
}

func calculateSize(rec *parse.Record) int {
	if len(rec.Children) == 0 {
		return rec.Pic.PicCount * max(1, rec.OccursCount)
	}

	return calculateGroupSize(rec)
}

func calculateGroupSize(rec *parse.Record) int {
	size := 0
	sizeStore := make(map[string]int)

	for _, child := range rec.Children {
		childSize := calculateSize(child)

		// Store the size of the child for redefines handling.
		sizeStore[child.Identifier] = childSize

		if child.Redefines == "" {
			size += childSize
		} else {
			// A field that redefines another can be of a different size.
			// To account for this, we subtract the size of the redefined field
			// and add the size of the redefining field as it takes precedent.
			redefinedChildSize, ok := sizeStore[child.Redefines]
			if !ok {
				// This should never happen as a redefined field should always be
				// processed at the same level before the redefining field.
				panic(fmt.Sprint("redefined field ", child.Redefines, " not found"))
			}
			size -= redefinedChildSize
			size += childSize

			// A redefined field size will be updated to the new field size for
			// future references.
			sizeStore[child.Redefines] = childSize
		}
	}

	return size * max(1, rec.OccursCount)
}

func newPositionTracker() *positionTracker {
	return &positionTracker{localPos: 1, globalPos: 1, recordStore: make(map[string]positionInfo)}
}

func (p *positionTracker) storeAndAdvancePos(identifier string, size int) {
	p.recordStore[identifier] = positionInfo{localStart: p.localPos, globalStart: p.globalPos}
	p.localPos += size
	p.globalPos += size
}

func (p *positionTracker) getStoredPos(identifier string) (int, int) {
	pos, ok := p.recordStore[identifier]
	if !ok {
		// This should never happen as a record should never be accessed before being stored.
		panic(fmt.Sprint("position for ", identifier, " not found"))
	}
	return pos.localStart, pos.globalStart
}

func defaultTypeMapping() map[parse.PicType]string {
	return map[parse.PicType]string{
		parse.Unsigned: "uint",
		parse.Signed:   "int",
		parse.Decimal:  "decimal.Decimal",
		parse.Alpha:    "string",
		parse.Unknown:  "string",
	}
}
