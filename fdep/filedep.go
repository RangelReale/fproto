package fdep

import (
	"fmt"
	"path"

	"github.com/RangelReale/fproto"
)

type FileDepType int

const (
	DepType_Own FileDepType = iota
	DepType_Imported
)

type FileDep struct {
	FilePath  string
	DepType   FileDepType
	Dep       *Dep
	ProtoFile *fproto.ProtoFile
}

func (fd *FileDep) GetType(name string) (*DepType, error) {
	t, err := fd.GetTypes(name)
	if err != nil {
		return nil, err
	}

	if len(t) > 1 {
		return nil, fmt.Errorf("More than one type found for '%s'", name)
	}

	return t[0], nil
}

func (fd *FileDep) GetTypes(name string) ([]*DepType, error) {
	return fd.Dep.internalGetTypes(name, fd)
}

func (fd *FileDep) IsSame(filedep *FileDep) bool {
	if fd == filedep {
		return true
	}

	if fd.FilePath == filedep.FilePath && fd.ProtoFile.PackageName == filedep.ProtoFile.PackageName {
		return true
	}

	return false
}

func (fd *FileDep) GoPackage() string {
	for _, o := range fd.ProtoFile.Options {
		if o.Name == "go_package" {
			return o.Value
		}
	}
	return path.Dir(fd.ProtoFile.PackageName)
}
