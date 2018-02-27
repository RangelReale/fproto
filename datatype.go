package fproto

import (
	"strings"
)

// ScalarType is an enumeration which represents all known supported scalar
// field datatypes.
type ScalarType int

const (
	AnyScalar ScalarType = iota + 1
	BoolScalar
	BytesScalar
	DoubleScalar
	FloatScalar
	Fixed32Scalar
	Fixed64Scalar
	Int32Scalar
	Int64Scalar
	Sfixed32Scalar
	Sfixed64Scalar
	Sint32Scalar
	Sint64Scalar
	StringScalar
	Uint32Scalar
	Uint64Scalar
)

var scalarLookupMap = map[string]ScalarType{
	"any":      AnyScalar,
	"bool":     BoolScalar,
	"bytes":    BytesScalar,
	"double":   DoubleScalar,
	"float":    FloatScalar,
	"fixed32":  Fixed32Scalar,
	"fixed64":  Fixed64Scalar,
	"int32":    Int32Scalar,
	"int64":    Int64Scalar,
	"sfixed32": Sfixed32Scalar,
	"sfixed64": Sfixed64Scalar,
	"sint32":   Sint32Scalar,
	"sint64":   Sint64Scalar,
	"string":   StringScalar,
	"uint32":   Uint32Scalar,
	"uint64":   Uint64Scalar,
}

func (s ScalarType) String() string {
	for n, v := range scalarLookupMap {
		if v == s {
			return n
		}
	}
	return ""
}

func ParseScalarType(s string) (ScalarType, bool) {
	key := strings.ToLower(s)
	if st, ok := scalarLookupMap[key]; ok {
		return st, true
	} else {
		return AnyScalar, false
	}
}
