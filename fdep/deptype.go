package fdep

import (
	"fmt"
	"strings"

	"github.com/RangelReale/fproto"
)

// DepType represents one type into one .proto file.
type DepType struct {
	// The file where the type is defined. Can be nil if scalar.
	FileDep *FileDep

	// The alias of the type. Can vary depending of how the type was requested.
	// When returned on the FileDep scope, it can be blank if it is contained on
	// the file itself.
	Alias string

	// The original alias of the type, independently of how it was requested.
	OriginalAlias string

	// The name of the type.
	Name string

	// Scalar type if scalar
	ScalarType *fproto.ScalarType

	// The *fproto.XXXElement corresponding to the type. Can be nil if scalar.
	Item fproto.FProtoElement
}

// Creates a new DepType from a file's element.
func NewDepTypeFromElement(filedep *FileDep, element fproto.FProtoElement) *DepType {
	return &DepType{
		FileDep:       filedep,
		Alias:         filedep.OriginalAlias(),
		OriginalAlias: filedep.OriginalAlias(),
		Name:          fproto.ScopedName(element),
		Item:          element,
	}
}

// Returns the name plus alias, if available
func (d *DepType) FullName() string {
	if d.Alias != "" {
		return fmt.Sprintf("%s.%s", d.Alias, d.Name)
	} else {
		return d.Name
	}
}

// Returns the name plus alias, if available
func (d *DepType) FullOriginalName() string {
	if d.OriginalAlias != "" {
		return fmt.Sprintf("%s.%s", d.OriginalAlias, d.Name)
	} else {
		return d.Name
	}
}

// Returns whether the field is pointer.
func (d *DepType) IsPointer() bool {
	switch d.Item.(type) {
	case *fproto.MessageElement:
		return true
	default:
		return false
	}
}

// Returns whether the field is scalar.
func (d *DepType) IsScalar() bool {
	return d.ScalarType != nil
}

// Returns one named type from the dependency, in relation to the current type.
//
// If multiple types are found for the same name, an error is issued.
// If there is this possibility, use the GetTypes method instead.
func (d *DepType) GetType(name string) (*DepType, error) {
	if d.FileDep == nil {
		return nil, nil
	}

	// first find normally in the full file
	dt, err := d.FileDep.GetType(name)
	if err != nil {
		return nil, err
	}

	if dt != nil {
		return dt, nil
	}

	// if not found, search in the current scope
	dt, err = d.FileDep.GetType(fmt.Sprintf("%s.%s", d.Name, name))
	if err != nil {
		return nil, err
	}

	if dt != nil {
		return dt, nil
	}

	// if not found, search in each dotted scope of the current alias
	if d.OriginalAlias != "" {
		scopes := strings.Split(d.OriginalAlias, ".")
		for si := 1; si <= len(scopes); si++ {
			dt, err = d.FileDep.GetType(fmt.Sprintf("%s.%s", strings.Join(scopes[:si], "."), name))
			if err != nil {
				return nil, err
			}

			if dt != nil {
				return dt, nil
			}
		}
	}

	// Not found in any method
	return nil, nil
}

// Returns all named types from the dependency, in relation to the current type.
//
// Use this method if there is a possibility that one name resolves to more than one type.
func (d *DepType) GetTypes(name string) ([]*DepType, error) {
	if d.FileDep == nil {
		return nil, nil
	}

	// first find normally in the full file
	dt, err := d.FileDep.GetTypes(name)
	if err != nil {
		return nil, err
	}

	if dt != nil && len(dt) > 0 {
		return dt, nil
	}

	// if not found, search in the current item scope
	return d.FileDep.GetTypes(fmt.Sprintf("%s.%s", d.Name, name))
}
