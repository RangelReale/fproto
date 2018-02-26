package fdep

import (
	"fmt"

	"github.com/RangelReale/fproto"
)

type FileDep struct {
	FilePath  string
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
