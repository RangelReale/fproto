package fdep

import (
	"fmt"
	"path"

	"github.com/RangelReale/fproto"
)

// The dependency file type.
type FileDepType int

const (
	// Your own project's proto files.
	DepType_Own FileDepType = iota

	// Imported proto file, that are not part of your project.
	DepType_Imported
)

// FileDep represents one .proto file into the dependency.
type FileDep struct {
	// The INTERNAL path of the .proto file, for example "google/protobuf/empty.proto"
	// This is NOT the filesystem path
	FilePath string

	// The type of the file dependency, whether it is your own file, or an imported one.
	DepType FileDepType

	// The parent dependency list this file is contained.
	Dep *Dep

	// The parsed proto file. Can be nil it was from an ignored dependency.
	ProtoFile *fproto.ProtoFile
}

// Returns one named type from the dependency, in relation to the current file.
// If the type is from the current file, the "Alias" field is blank.
//
// If multiple types are found for the same name, an error is issued.
// If there is this possibility, use the GetTypes method instead.
func (fd *FileDep) GetType(name string) (*DepType, error) {
	t, err := fd.GetTypes(name)
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

// Returns all named types from the dependency, in relation to the current file.
// If the type is from the current file, the "Alias" field is blank.
//
// Use this method if there is a possibility that one name resolves to more than one type.
func (fd *FileDep) GetTypes(name string) ([]*DepType, error) {
	return fd.Dep.internalGetTypes(name, fd)
}

// Checks if the passed FileDep refers to the same file as this one.
func (fd *FileDep) IsSame(filedep *FileDep) bool {
	if fd == filedep {
		return true
	}

	if fd.FilePath == filedep.FilePath && fd.ProtoFile.PackageName == filedep.ProtoFile.PackageName {
		return true
	}

	return false
}

func (fd *FileDep) OriginalAlias() string {
	if fd.ProtoFile != nil {
		return fd.ProtoFile.PackageName
	}
	return ""
}

// Checks if the passed FileDep refers to the same package as this one.
func (fd *FileDep) IsSamePackage(filedep *FileDep) bool {
	if fd == filedep {
		return true
	}

	if path.Dir(fd.FilePath) == path.Dir(filedep.FilePath) && fd.ProtoFile.PackageName == filedep.ProtoFile.PackageName {
		return true
	}

	return false
}

// Returns the go package of the file. If there is no "go_package" option, returns the "path" part of the package name.
func (fd *FileDep) GoPackage() string {
	for _, o := range fd.ProtoFile.Options {
		if o.Name == "go_package" {
			return o.Value
		}
	}
	return path.Dir(fd.ProtoFile.PackageName)
}
