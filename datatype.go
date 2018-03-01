package fproto

import (
	"strings"
)

// ScalarType is an enumeration which represents all known supported scalar
// field datatypes.
type ScalarType int

const (
	BoolScalar ScalarType = iota + 1
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

var scalarGoTypeLookupMap = map[ScalarType]string{
	BoolScalar:     "bool",
	BytesScalar:    "[]byte",
	DoubleScalar:   "float64",
	FloatScalar:    "float32",
	Fixed32Scalar:  "uint32",
	Fixed64Scalar:  "uint64",
	Int32Scalar:    "int32",
	Int64Scalar:    "int64",
	Sfixed32Scalar: "int32",
	Sfixed64Scalar: "int64",
	Sint32Scalar:   "int32",
	Sint64Scalar:   "int64",
	StringScalar:   "string",
	Uint32Scalar:   "uint32",
	Uint64Scalar:   "uint64",
}

// Returns the protobuf type string for the scalar type
func (s ScalarType) ProtoType() string {
	for n, v := range scalarLookupMap {
		if v == s {
			return n
		}
	}
	return ""
}

// Returns the go type string for the scalar type
func (s ScalarType) GoType() string {
	return scalarGoTypeLookupMap[s]
}

// Parses the protobuf type into ScalarType. The bool parameters indicates if
// the type is scalar or not
func ParseScalarType(s string) (ScalarType, bool) {
	key := strings.ToLower(s)
	if st, ok := scalarLookupMap[key]; ok {
		return st, true
	} else {
		return BoolScalar, false
	}
}
