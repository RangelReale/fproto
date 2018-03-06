package fdep

import "github.com/RangelReale/fproto"

// DepType represents one type into one .proto file.
type DepType struct {
	// The file where the type is defined.
	FileDep *FileDep

	// The alias of the type. Can vary depending of how the type was requested.
	// When returned on the FileDep scope, it can be blank if it is contained on
	// the file itself.
	Alias string

	// The name of the type.
	Name string

	// The *fproto.XXXElement corresponding to the type.
	Item fproto.FProtoElement
}

func (d *DepType) IsPointer() bool {
	switch d.Item.(type) {
	case *fproto.MessageElement:
		return true
	default:
		return false
	}
}
