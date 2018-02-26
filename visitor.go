package fproto

import (
	"fmt"

	"github.com/emicklei/proto"
)

type Visitor struct {
	protofile *ProtoFile
	scope     interface{}
	err       error
}

func NewVisitor(protofile *ProtoFile) *Visitor {
	return &Visitor{
		protofile: protofile,
		scope:     protofile,
	}
}

func (v *Visitor) Err() error {
	return v.err
}

func (v *Visitor) errInvalidScope(item, name string) {
	v.err = &InvalidScope{fmt.Sprint("Invalid scope for item '%s' (%s)")}
}

func (v *Visitor) visitElements(ml []proto.Visitee) {
	for _, m := range ml {
		m.Accept(v)
	}
}

func (v *Visitor) visitOptions(ml []*proto.Option) {
	for _, m := range ml {
		m.Accept(v)
	}
}

func (v *Visitor) VisitMessage(m *proto.Message) {
	if v.err != nil {
		return
	}

	// create new message element
	newm := &MessageElement{
		Name:     m.Name,
		IsExtend: m.IsExtend,
	}

	// visit children
	nv := &Visitor{
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

func (v *Visitor) VisitService(s *proto.Service) {
	if v.err != nil {
		return
	}

	// create new service element
	news := &ServiceElement{
		Name: s.Name,
	}

	// visit children
	nv := &Visitor{
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

func (v *Visitor) VisitSyntax(s *proto.Syntax) {
	if v.err != nil {
		return
	}

	v.protofile.Syntax = s.Value
}

func (v *Visitor) VisitPackage(p *proto.Package) {
	if v.err != nil {
		return
	}

	v.protofile.PackageName = p.Name
}

func (v *Visitor) VisitOption(o *proto.Option) {
	if v.err != nil {
		return
	}

	if el, ok := v.scope.(iAddOption); ok {
		el.addOptionElement(&OptionElement{
			Name:  o.Name,
			Value: o.Constant.Source,
		})
	} else {
		v.errInvalidScope("public dependency", o.Name)
	}

}

func (v *Visitor) VisitImport(i *proto.Import) {
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

func (v *Visitor) VisitNormalField(i *proto.NormalField) {
	if v.err != nil {
		return
	}

	// create field
	newf := &FieldElement{
		Name:     i.Name,
		Type:     i.Type,
		Repeated: i.Repeated,
		Optional: i.Optional,
		Required: i.Required,
		Tag:      i.Sequence,
	}

	// visit children
	nv := &Visitor{
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
		el.addFieldElement(newf)
	} else {
		v.errInvalidScope("field", i.Name)
	}
}

func (v *Visitor) VisitEnumField(i *proto.EnumField) {
	if v.err != nil {
		return
	}

	// create enum constant
	newe := &EnumConstantElement{
		Name: i.Name,
		Tag:  i.Integer,
	}

	// visit children
	nv := &Visitor{
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

func (v *Visitor) VisitEnum(e *proto.Enum) {
	if v.err != nil {
		return
	}

	// create enum
	newe := &EnumElement{
		Name: e.Name,
	}

	// visit children
	nv := &Visitor{
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

func (v *Visitor) VisitComment(e *proto.Comment) {
	if v.err != nil {
		return
	}

}

func (v *Visitor) VisitOneof(o *proto.Oneof) {
	if v.err != nil {
		return
	}

	// create oneof
	newo := &OneOfElement{
		Name: o.Name,
	}

	// visit children
	nv := &Visitor{
		protofile: &ProtoFile{},
		scope:     newo,
	}
	nv.visitElements(o.Elements)
	if nv.Err() != nil {
		v.err = nv.Err()
		return
	}

	// add to scope
	if el, ok := v.scope.(iAddOneOf); ok {
		el.addOneOfElement(newo)
	} else {
		v.errInvalidScope("oneof", o.Name)
	}

}

func (v *Visitor) VisitOneofField(o *proto.OneOfField) {
	if v.err != nil {
		return
	}

	// create field
	newf := &FieldElement{
		Name: o.Name,
		Type: o.Type,
		Tag:  o.Sequence,
	}

	// visit children
	nv := &Visitor{
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
		el.addFieldElement(newf)
	} else {
		v.errInvalidScope("oneof field", o.Name)
	}
}

func (v *Visitor) VisitReserved(r *proto.Reserved) {
	if v.err != nil {
		return
	}

	for _, rr := range r.Ranges {
		// add to scope
		if el, ok := v.scope.(iAddReservedRange); ok {
			el.addReservedRangeElement(&ReservedRangeElement{
				Start: rr.From,
				End:   rr.To,
				IsMax: rr.Max,
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

func (v *Visitor) VisitRPC(r *proto.RPC) {
	if v.err != nil {
		return
	}

	// create RPC
	newr := &RPCElement{
		Name:            r.Name,
		RequestType:     r.RequestType,
		StreamsRequest:  r.StreamsRequest,
		ResponseType:    r.ReturnsType,
		StreamsResponse: r.StreamsReturns,
	}

	// visit children
	nv := &Visitor{
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

func (v *Visitor) VisitMapField(f *proto.MapField) {
	if v.err != nil {
		return
	}

	// create field
	newf := &MapFieldElement{
		FieldElement: &FieldElement{
			Name: f.Name,
			Type: f.Type,
			//Repeated: f.Repeated,
			//Optional: f.Optional,
			//Required: f.Required,
			Tag: f.Sequence,
		},
		KeyType: f.KeyType,
	}

	// visit children
	nv := &Visitor{
		protofile: &ProtoFile{},
		scope:     newf,
	}
	nv.visitOptions(f.Options)
	if nv.Err() != nil {
		v.err = nv.Err()
		return
	}

	// add to scope
	if el, ok := v.scope.(iAddMapField); ok {
		el.addMapFieldElement(newf)
	} else {
		v.errInvalidScope("map field", f.Name)
	}
}

// proto2
func (v *Visitor) VisitGroup(g *proto.Group) {
	if v.err != nil {
		return
	}

}

func (v *Visitor) VisitExtensions(e *proto.Extensions) {
	if v.err != nil {
		return
	}

	for _, rr := range e.Ranges {
		// add to scope
		if el, ok := v.scope.(iAddExtensions); ok {
			el.addExtensionsElement(&ExtensionsElement{
				Start: rr.From,
				End:   rr.To,
				IsMax: rr.Max,
			})
		} else {
			v.errInvalidScope("extensions", "extension")
		}
	}
}
