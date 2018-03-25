package fdep

import (
	"fmt"

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

// Creates a new DepType
func NewDepType(filedep *FileDep, alias string, originalAlias string, name string, item fproto.FProtoElement) *DepType {
	return &DepType{
		FileDep:       filedep,
		Alias:         alias,
		OriginalAlias: originalAlias,
		Name:          name,
		Item:          item,
	}
}

// Creates a new DepType for a scalar.
func NewDepTypeScalar(scalarType fproto.ScalarType) *DepType {
	return &DepType{
		Name:       scalarType.ProtoType(),
		ScalarType: &scalarType,
	}
}

// Creates a new DepType from a file's element.
func NewDepTypeFromElement(filedep *FileDep, element fproto.FProtoElement) *DepType {
	return NewDepType(filedep, filedep.OriginalAlias(), filedep.OriginalAlias(), fproto.ScopedName(element), element)
}

// Returns the name plus alias, if available
func (d *DepType) IsSame(od *DepType) bool {
	if d.IsScalar() != od.IsScalar() ||
		(d.IsScalar() && od.IsScalar() && *d.ScalarType != *od.ScalarType) {
		return false
	}

	if d.FileDep == nil || od.FileDep == nil {
		return false
	}

	if d.FileDep.FilePath != od.FileDep.FilePath {
		return false
	}

	if d.OriginalAlias != od.OriginalAlias || d.Name != od.Name {
		return false
	}

	return true
}

// Returns the parent deptype, or nil if root
func (d *DepType) Parent() *DepType {
	if d.Item != nil && d.Item.ParentElement() != nil {
		return NewDepTypeFromElement(d.FileDep, d.Item.ParentElement())
	}
	return nil
}

// Returns up to the n-th parent if possible, excluding the ProtFile.
// The second return value is the amount found.
func (d *DepType) SkipParents(n int) (*DepType, int) {
	cur := d.Item
	ct := 0
	for n > 0 {
		if cur != nil && cur.ParentElement() != nil {
			if _, ispfile := cur.ParentElement().(*fproto.ProtoFile); !ispfile {
				cur = cur.ParentElement()
				ct++
				n--
			} else {
				break
			}
		} else {
			break
		}
	}

	if cur == nil {
		return nil, 0
	}

	return NewDepTypeFromElement(d.FileDep, cur), ct
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

// Returns whether the field can be a pointer. (the scalar []byte cannot)
func (d *DepType) CanPointer() bool {
	if d.ScalarType != nil && *d.ScalarType == fproto.BytesScalar {
		return false
	}

	return true
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
	t, err := d.GetTypes(name)
	if err != nil {
		return nil, err
	}

	if len(t) == 0 {
		return nil, nil
	} else if len(t) > 1 {
		return nil, fmt.Errorf("More than one type found for '%s'", name)
	}

	return t[0], nil
}

// Like GetType, but returns an error if not found
func (d *DepType) MustGetType(name string) (*DepType, error) {
	t, err := d.GetType(name)
	if err != nil {
		return nil, err
	}
	if t == nil {
		return nil, fmt.Errorf("Type %s not found in %s", name, d.FullOriginalName())
	}
	return t, nil
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

// Returns a list of extension packages for this type.
func (d *DepType) ExtensionPackages() []string {
	if d.FileDep != nil {
		return d.FileDep.Dep.GetExtensions(d.FileDep, d.OriginalAlias, d.Name)
	}
	return nil
}

// Returns a list of DepTypes for the extensions of this type.
func (d *DepType) GetTypeExtensions() (map[string]*DepType, error) {
	pkgs := d.ExtensionPackages()
	if len(pkgs) == 0 {
		return nil, nil
	}

	ret := make(map[string]*DepType)

	for _, p := range pkgs {
		edt, err := d.GetTypeExtension(p)
		if err != nil {
			return nil, err
		}
		if edt != nil {
			ret[p] = edt
		}
	}

	return ret, nil
}

// Returns a single Deptype for an extensions of this type named by a package.
func (d *DepType) GetTypeExtension(extensionPkg string) (*DepType, error) {
	if d.FileDep != nil {
		return d.FileDep.Dep.GetTypeExtension(d.FullOriginalName(), extensionPkg)
	}
	return nil, nil
}
