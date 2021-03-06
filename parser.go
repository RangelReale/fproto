package fproto

import (
	"io"

	"github.com/emicklei/proto"
)

// Parses an io.Reader corresponding to a .proto file into a ProtoFile struct
func Parse(r io.Reader) (*ProtoFile, error) {
	parser := proto.NewParser(r)
	definition, err := parser.Parse()
	if err != nil {
		return nil, err
	}

	protofile := &ProtoFile{}

	v := newVisitor(protofile)

	for _, element := range definition.Elements {
		element.Accept(v)
	}

	if v.Err() != nil {
		return nil, err
	}

	return protofile, nil
}
