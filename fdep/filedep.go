package fdep

import (
	"github.com/RangelReale/fproto"
)

type FileDep struct {
	Dep       *Dep
	ProtoFile *fproto.ProtoFile
}
