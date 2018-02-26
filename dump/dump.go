package dump

import (
	"fmt"
	"io"
	"strings"

	"github.com/emicklei/proto"
)

type Visitor struct {
	w     io.Writer
	ident int
}

func NewVisitor(w io.Writer) *Visitor {
	return &Visitor{w: w}
}

func (v *Visitor) PrintLine(format string, a ...interface{}) {
	if v.ident > 0 {
		fmt.Fprintf(v.w, strings.Repeat(" ", v.ident*4))
	}
	fmt.Fprintf(v.w, format, a...)
	fmt.Fprintf(v.w, "\n")
}

func (v *Visitor) PrintNested(nested []proto.Visitee) {
	nv := &Visitor{w: v.w, ident: v.ident + 1}
	for _, e := range nested {
		e.Accept(nv)
	}
}

func (v *Visitor) PrintOptions(options []*proto.Option) {
	nv := &Visitor{w: v.w, ident: v.ident + 1}
	for _, e := range options {
		e.Accept(nv)
	}
}

func (v *Visitor) VisitMessage(m *proto.Message) {
	var tags []string
	if m.IsExtend {
		tags = append(tags, "extend")
	}

	var stags string
	if len(tags) > 0 {
		stags = fmt.Sprintf(" {%s}", strings.Join(tags, ","))
	}

	v.PrintLine("* Message: %s%s", m.Name, stags)
	v.PrintNested(m.Elements)
}

func (v *Visitor) VisitService(s *proto.Service) {
	v.PrintLine("* Service: %s", s.Name)
	v.PrintNested(s.Elements)
}

func (v *Visitor) VisitSyntax(s *proto.Syntax) {
	v.PrintLine("* Syntax: %s", s.Value)
}

func (v *Visitor) VisitPackage(p *proto.Package) {
	v.PrintLine("* Package: %s", p.Name)
}

func (v *Visitor) VisitOption(o *proto.Option) {
	v.PrintLine("* Option: %s = %s", o.Name, o.Constant.SourceRepresentation())
}

func (v *Visitor) VisitImport(i *proto.Import) {
	v.PrintLine("* Import: %s", i.Filename)
}

func (v *Visitor) VisitNormalField(i *proto.NormalField) {
	var tags []string
	if i.Repeated {
		tags = append(tags, "repeated")
	}
	if i.Optional {
		tags = append(tags, "optional")
	}
	if i.Required {
		tags = append(tags, "required")
	}

	var stags string
	if len(tags) > 0 {
		stags = fmt.Sprintf(" {%s}", strings.Join(tags, ","))
	}

	v.PrintLine("* Normal Field: %s [%s]%s", i.Name, i.Type, stags)
	v.PrintOptions(i.Options)
}

func (v *Visitor) VisitEnumField(i *proto.EnumField) {
	v.PrintLine("* Enum Field: %s", i.Name)
	v.PrintNested(i.Elements)
}

func (v *Visitor) VisitEnum(e *proto.Enum) {
	v.PrintLine("* Enum: %s", e.Name)
	v.PrintNested(e.Elements)
}

func (v *Visitor) VisitComment(e *proto.Comment) {
	/*
		var c []string
		for _, cc := range e.Lines {
			c = append(c, strings.TrimSpace(cc))
		}
		v.PrintLine("* Comment: %s", strings.Join(c, ", "))
	*/
}

func (v *Visitor) VisitOneof(o *proto.Oneof) {
	v.PrintLine("* OneOf: %s", o.Name)
	v.PrintNested(o.Elements)
}

func (v *Visitor) VisitOneofField(o *proto.OneOfField) {
	v.PrintLine("* OneOfField: %s [%s]", o.Name, o.Type)
	v.PrintOptions(o.Options)
}

func (v *Visitor) VisitReserved(r *proto.Reserved) {
	var rs []string
	for _, rr := range r.Ranges {
		x := fmt.Sprintf("%d", rr.From)
		if rr.Max {
			x += " to MAX"
		} else if rr.To > rr.From {
			x += fmt.Sprintf(" to %d", rr.To)
		}

		rs = append(rs, x)
	}

	for _, rr := range r.FieldNames {
		rs = append(rs, rr)
	}

	v.PrintLine("* Reserved: %s", strings.Join(rs, ", "))
}

func (v *Visitor) VisitRPC(r *proto.RPC) {
	v.PrintLine("* RPC: %s (%s) returns %s", r.Name, r.RequestType, r.ReturnsType)
	v.PrintNested(r.Elements)
	v.PrintOptions(r.Options)
}

func (v *Visitor) VisitMapField(f *proto.MapField) {
	v.PrintLine("* MapField: %s [%s[%s]]", f.Name, f.KeyType, f.Type)
	v.PrintOptions(f.Options)
}

// proto2
func (v *Visitor) VisitGroup(g *proto.Group) {
	var tags []string
	if g.Repeated {
		tags = append(tags, "repeated")
	}
	if g.Optional {
		tags = append(tags, "optional")
	}
	if g.Required {
		tags = append(tags, "required")
	}

	var stags string
	if len(tags) > 0 {
		stags = fmt.Sprintf(" {%s}", strings.Join(tags, ","))
	}

	v.PrintLine("* Group: %s%s", g.Name, stags)
	v.PrintNested(g.Elements)
}

func (v *Visitor) VisitExtensions(e *proto.Extensions) {
	var rs []string
	for _, rr := range e.Ranges {
		x := fmt.Sprintf("%d", rr.From)
		if rr.Max {
			x += " to MAX"
		} else if rr.To > rr.From {
			x += fmt.Sprintf(" to %d", rr.To)
		}

		rs = append(rs, x)
	}

	v.PrintLine("* Extensions: %s", strings.Join(rs, ", "))
}
