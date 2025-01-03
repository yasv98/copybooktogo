// Code generated by "enumer -type PicType -output pictype_enumer.generated.go -linecomment"; DO NOT EDIT.

package parse

import (
	"fmt"
)

const _PicTypeName = "unknownunsignedsigneddecimalalpha"

var _PicTypeIndex = [...]uint8{0, 7, 15, 21, 28, 33}

func (i PicType) String() string {
	if i < 0 || i >= PicType(len(_PicTypeIndex)-1) {
		return fmt.Sprintf("PicType(%d)", i)
	}
	return _PicTypeName[_PicTypeIndex[i]:_PicTypeIndex[i+1]]
}

var _PicTypeValues = []PicType{0, 1, 2, 3, 4}

var _PicTypeNameToValueMap = map[string]PicType{
	_PicTypeName[0:7]:   0,
	_PicTypeName[7:15]:  1,
	_PicTypeName[15:21]: 2,
	_PicTypeName[21:28]: 3,
	_PicTypeName[28:33]: 4,
}

// PicTypeString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func PicTypeString(s string) (PicType, error) {
	if val, ok := _PicTypeNameToValueMap[s]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to PicType values", s)
}

// PicTypeValues returns all values of the enum
func PicTypeValues() []PicType {
	return _PicTypeValues
}

// IsAPicType returns "true" if the value is listed in the enum definition. "false" otherwise
func (i PicType) IsAPicType() bool {
	for _, v := range _PicTypeValues {
		if i == v {
			return true
		}
	}
	return false
}
