package fproto

import (
	"fmt"
	"strings"

	"github.com/emicklei/proto"
)

// Internal classes to parse the .proto file into the ProtoFile struct.

type visitor struct {
	protofile *ProtoFile
	scope     FProtoElement
	err       error
}

func newVisitor(protofile *ProtoFile) *visitor {
	return &visitor{
		protofile: protofile,
		scope:     protofile,
	}
}

func (v *visitor) Err() error {
	return v.err
}

func (v *visitor) errInvalidScope(item, name string) {
	v.err = &InvalidScope{fmt.Sprintf("Invalid scope for item '%s' (%s)", item, name)}
}

func (v *visitor) visitElements(ml []proto.Visitee) {
	for _, m := range ml {
		m.Accept(v)
	}
}

func (v *visitor) visitOptions(ml []*proto.Option) {
	for _, m := range ml {
		m.Accept(v)
	}
}

func (v *visitor) VisitMessage(m *proto.Message) {
	if v.err != nil {
		return
	}

	// create new message element
	newm := &MessageElement{
		Parent:   v.scope,
		Name:     m.Name,
		IsExtend: m.IsExtend,
	}

	// visit children
	nv := &visitor{
		protofile: &ProtoFile{},
		scope:     newm,
	}
	nv.visitElements(m.Elements)
	if nv.Err() != nil {
		v.err = nv.Err()
		return
	}

	// add to scope
	if e, ok := v.scope.(iAddMessage); ok {
		e.addMessageElement(newm)
	} else {
		v.errInvalidScope("message", m.Name)
	}
}

func (v *visitor) VisitService(s *proto.Service) {
	if v.err != nil {
		return
	}

	// create new service element
	news := &ServiceElement{
		Parent: v.scope,
		Name:   s.Name,
	}

	// visit children
	nv := &visitor{
		protofile: &ProtoFile{},
		scope:     news,
	}
	nv.visitElements(s.Elements)
	if nv.Err() != nil {
		v.err = nv.Err()
		return
	}

	// add to scope
	if e, ok := v.scope.(iAddService); ok {
		e.addServiceElement(news)
	} else {
		v.errInvalidScope("service", s.Name)
	}

}

func (v *visitor) VisitSyntax(s *proto.Syntax) {
	if v.err != nil {
		return
	}

	v.protofile.Syntax = s.Value
}

func (v *visitor) VisitPackage(p *proto.Package) {
	if v.err != nil {
		return
	}

	v.protofile.PackageName = p.Name
}

func (v *visitor) VisitOption(o *proto.Option) {
	if v.err != nil {
		return
	}

	if el, ok := v.scope.(iAddOption); ok {
		oname := o.Name
		is_parenthesized := false

		if strings.HasPrefix(oname, "(") {
			is_parenthesized = true
			oname = strings.TrimPrefix(strings.TrimSuffix(oname, ")"), "(")
		}

		el.addOptionElement(&OptionElement{
			Parent:          v.scope,
			Name:            oname,
			Value:           o.Constant.Source,
			IsParenthesized: is_parenthesized,
		})
	} else {
		v.errInvalidScope("public dependency", o.Name)
	}

}

func (v *visitor) VisitImport(i *proto.Import) {
	if v.err != nil {
		return
	}

	if i.Kind == "public" {
		if el, ok := v.scope.(iAddPublicDependency); ok {
			el.addPublicDependency(i.Filename)
		} else {
			v.errInvalidScope("public dependency", i.Filename)
		}
	} else if i.Kind == "weak" {
		if el, ok := v.scope.(iAddWeakDependency); ok {
			el.addWeakDependency(i.Filename)
		} else {
			v.errInvalidScope("weak dependency", i.Filename)
		}
	} else if i.Kind == "" {
		if el, ok := v.scope.(iAddDependency); ok {
			el.addDependency(i.Filename)
		} else {
			v.errInvalidScope("dependency", i.Filename)
		}
	}
}

func (v *visitor) VisitNormalField(i *proto.NormalField) {
	if v.err != nil {
		return
	}

	// create field
	newf := &FieldElement{
		Parent:   v.scope,
		Name:     i.Name,
		Type:     i.Type,
		Repeated: i.Repeated,
		Optional: i.Optional,
		Required: i.Required,
		Tag:      i.Sequence,
	}

	// visit children
	nv := &visitor{
		protofile: &ProtoFile{},
		scope:     newf,
	}
	nv.visitOptions(i.Options)
	if nv.Err() != nil {
		v.err = nv.Err()
		return
	}

	// add to scope
	if el, ok := v.scope.(iAddField); ok {
		el.addField(newf)
	} else {
		v.errInvalidScope("field", i.Name)
	}
}

func (v *visitor) VisitEnumField(i *proto.EnumField) {
	if v.err != nil {
		return
	}

	// create enum constant
	newe := &EnumConstantElement{
		Parent: v.scope,
		Name:   i.Name,
		Tag:    i.Integer,
	}

	// visit children
	nv := &visitor{
		protofile: &ProtoFile{},
		scope:     newe,
	}
	nv.visitElements(i.Elements)
	if nv.Err() != nil {
		v.err = nv.Err()
		return
	}

	// add to scope
	if el, ok := v.scope.(iAddEnumConstant); ok {
		el.addEnumConstantElement(newe)
	} else {
		v.errInvalidScope("enum constant", i.Name)
	}

}

func (v *visitor) VisitEnum(e *proto.Enum) {
	if v.err != nil {
		return
	}

	// create enum
	newe := &EnumElement{
		Name:   e.Name,
		Parent: v.scope,
	}

	// visit children
	nv := &visitor{
		protofile: &ProtoFile{},
		scope:     newe,
	}
	nv.visitElements(e.Elements)
	if nv.Err() != nil {
		v.err = nv.Err()
		return
	}

	// add to scope
	if el, ok := v.scope.(iAddEnum); ok {
		el.addEnumElement(newe)
	} else {
		v.errInvalidScope("enum", e.Name)
	}
}

func (v *visitor) VisitComment(e *proto.Comment) {
	if v.err != nil {
		return
	}

}

func (v *visitor) VisitOneof(o *proto.Oneof) {
	if v.err != nil {
		return
	}

	// create oneof
	newo := &OneOfElement{
		Parent: v.scope,
		Name:   o.Name,
	}

	// visit children
	nv := &visitor{
		protofile: &ProtoFile{},
		scope:     newo,
	}
	nv.visitElements(o.Elements)
	if nv.Err() != nil {
		v.err = nv.Err()
		return
	}

	// add to scope
	if el, ok := v.scope.(iAddField); ok {
		el.addField(newo)
	} else {
		v.errInvalidScope("oneof", o.Name)
	}
}

func (v *visitor) VisitOneofField(o *proto.OneOfField) {
	if v.err != nil {
		return
	}

	// create field
	newf := &FieldElement{
		Parent: v.scope,
		Name:   o.Name,
		Type:   o.Type,
		Tag:    o.Sequence,
	}

	// visit children
	nv := &visitor{
		protofile: &ProtoFile{},
		scope:     newf,
	}
	nv.visitOptions(o.Options)
	if nv.Err() != nil {
		v.err = nv.Err()
		return
	}

	// add to scope
	if el, ok := v.scope.(iAddField); ok {
		el.addField(newf)
	} else {
		v.errInvalidScope("oneof field", o.Name)
	}
}

func (v *visitor) VisitReserved(r *proto.Reserved) {
	if v.err != nil {
		return
	}

	for _, rr := range r.Ranges {
		// add to scope
		if el, ok := v.scope.(iAddReservedRange); ok {
			el.addReservedRangeElement(&ReservedRangeElement{
				Parent: v.scope,
				Start:  rr.From,
				End:    rr.To,
				IsMax:  rr.Max,
			})
		} else {
			v.errInvalidScope("reserved range", "reserved")
		}
	}

	for _, rr := range r.FieldNames {
		// add to scope
		if el, ok := v.scope.(iAddReservedName); ok {
			el.addReservedName(rr)
		} else {
			v.errInvalidScope("reserved name", rr)
		}
	}
}

func (v *visitor) VisitRPC(r *proto.RPC) {
	if v.err != nil {
		return
	}

	// create RPC
	newr := &RPCElement{
		Parent:          v.scope,
		Name:            r.Name,
		RequestType:     r.RequestType,
		StreamsRequest:  r.StreamsRequest,
		ResponseType:    r.ReturnsType,
		StreamsResponse: r.StreamsReturns,
	}

	// visit children
	nv := &visitor{
		protofile: &ProtoFile{},
		scope:     newr,
	}
	nv.visitElements(r.Elements)
	nv.visitOptions(r.Options)
	if nv.Err() != nil {
		v.err = nv.Err()
		return
	}

	// add to scope
	if el, ok := v.scope.(iAddRPC); ok {
		el.addRPCElement(newr)
	} else {
		v.errInvalidScope("rpc", r.Name)
	}
}

func (v *visitor) VisitMapField(f *proto.MapField) {
	if v.err != nil {
		return
	}

	// create field
	newf := &MapFieldElement{
		Parent: v.scope,
		FieldElement: &FieldElement{
			Parent: v.scope,
			Name:   f.Name,
			Type:   f.Type,
			//Repeated: f.Repeated,
			//Optional: f.Optional,
			//Required: f.Required,
			Tag: f.Sequence,
		},
		KeyType: f.KeyType,
	}

	// visit children
	nv := &visitor{
		protofile: &ProtoFile{},
		scope:     newf,
	}
	nv.visitOptions(f.Options)
	if nv.Err() != nil {
		v.err = nv.Err()
		return
	}

	// add to scope
	if el, ok := v.scope.(iAddField); ok {
		el.addField(newf)
	} else {
		v.errInvalidScope("map field", f.Name)
	}
}

// proto2
func (v *visitor) VisitGroup(g *proto.Group) {
	if v.err != nil {
		return
	}

}

func (v *visitor) VisitExtensions(e *proto.Extensions) {
	if v.err != nil {
		return
	}

	for _, rr := range e.Ranges {
		// add to scope
		if el, ok := v.scope.(iAddExtensions); ok {
			el.addExtensionsElement(&ExtensionsElement{
				Parent: v.scope,
				Start:  rr.From,
				End:    rr.To,
				IsMax:  rr.Max,
			})
		} else {
			v.errInvalidScope("extensions", "extension")
		}
	}
}
