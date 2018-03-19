package fdep

import (
	"fmt"
	"path"
	"strings"

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

func (dt FileDepType) String() string {
	switch dt {
	case DepType_Own:
		return "OWN"
	case DepType_Imported:
		return "IMPORTED"
	default:
		return "UNKNOWN"
	}
}

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
// If not found using the name, the current file's package is searched recursivelly
// appending the name.
//
// Use this method if there is a possibility that one name resolves to more than one type.
func (fd *FileDep) GetTypes(name string) ([]*DepType, error) {
	t, err := fd.Dep.internalGetTypes(name, fd)
	if err != nil {
		return nil, err
	}
	if len(t) > 0 {
		return t, nil
	}

	// if not found, search in each dotted scope of the current file's package
	if fd.ProtoFile != nil && fd.ProtoFile.PackageName != "" {
		scopes := strings.Split(fd.ProtoFile.PackageName, ".")
		for si := 1; si <= len(scopes); si++ {
			t, err = fd.Dep.internalGetTypes(fmt.Sprintf("%s.%s", strings.Join(scopes[:si], "."), name), fd)
			if err != nil {
				return nil, err
			}

			if t != nil {
				return t, nil
			}
		}
	}

	// Not found in any method
	return nil, nil
}

func (fd *FileDep) GetFileOfName(name string) (*FileDepOfName, error) {
	return fd.Dep.GetFileOfName(name)
}

func (fd *FileDep) GetFilesOfName(name string) ([]*FileDepOfName, error) {
	return fd.Dep.GetFilesOfName(name)
}

// Checks if the passed FileDep refers to the same file as this one.
func (fd *FileDep) IsSame(filedep *FileDep) bool {
	if filedep == nil {
		return false
	}

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
			return o.Value.String()
		}
	}
	return path.Dir(fd.ProtoFile.PackageName)
}

// Result of GetFilesOfName
type FileDepOfName struct {
	// File
	FileDep *FileDep
	// Package name
	Package string
	// Rest of name excluding the package name
	Name string
}
