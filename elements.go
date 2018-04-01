package fproto

import "bytes"

// Comment one or more comment text lines, either in c- or c++ style.
type Comment struct {
	// Lines are comment text lines without prefixes //, ///, /* or suffix */
	Lines      []string
	Cstyle     bool // refers to /* ... */,  C++ style is using //
	ExtraSlash bool // is true if the comment starts with 3 slashes
}

// Literal value from source
type Literal struct {
	Source   string
	IsString bool
}

// SourceRepresentation returns the source (if quoted then use double quote).
func (l Literal) SourceRepresentation() string {
	var buf bytes.Buffer
	if l.IsString {
		buf.WriteRune('"')
	}
	buf.WriteString(l.Source)
	if l.IsString {
		buf.WriteRune('"')
	}
	return buf.String()
}

// Raw string content
func (l Literal) String() string {
	return l.Source
}

// OptionElement is a datastructure which models
// the option construct in a protobuf file. Option constructs
// exist at various levels/contexts like file, message etc.
type OptionElement struct {
	Parent            FProtoElement
	Name              string
	ParenthesizedName string
	IsParenthesized   bool
	Value             Literal
	AggregatedValues  map[string]Literal
	Comment           *Comment
}

// EnumConstantElement is a datastructure which models
// the fields within an enum construct. Enum constants can
// also have inline options specified.
type EnumConstantElement struct {
	Parent  FProtoElement
	Name    string
	Comment *Comment
	Options []*OptionElement
	Tag     int
}

// EnumElement is a datastructure which models
// the enum construct in a protobuf file. Enums are
// defined standalone or as nested entities within messages.
type EnumElement struct {
	Parent        FProtoElement
	Name          string
	Comment       *Comment
	Options       []*OptionElement
	EnumConstants []*EnumConstantElement
}

// RPCElement is a datastructure which models
// the rpc construct in a protobuf file. RPCs are defined
// nested within ServiceElements.
type RPCElement struct {
	Parent          FProtoElement
	Name            string
	Comment         *Comment
	Options         []*OptionElement
	RequestType     string
	StreamsRequest  bool
	ResponseType    string
	StreamsResponse bool
}

// ServiceElement is a datastructure which models
// the service construct in a protobuf file. Service
// construct defines the rpcs (apis) for the service.
type ServiceElement struct {
	Parent  FProtoElement
	Name    string
	Comment *Comment
	Options []*OptionElement
	RPCs    []*RPCElement
}

// Can be:
// - FieldElement
// - MapFieldElement
// - OneofFieldElement
type FieldElementTag interface {
	FProtoElement
	FieldName() string
	// Returns the smallest field tag (oneof can have more than one tag
	FirstFieldTag() int
}

// FieldElement is a datastructure which models
// a field of a message, a field of a oneof element
// or an entry in the extend declaration in a protobuf file.
type FieldElement struct {
	Parent   FProtoElement
	Name     string
	Comment  *Comment
	Options  []*OptionElement
	Repeated bool
	Optional bool // proto2
	Required bool // proto2
	Type     string
	Tag      int
}

func (f *FieldElement) FieldName() string {
	return f.Name
}

func (f *FieldElement) FirstFieldTag() int {
	return f.Tag
}

// MapFieldElement is a datastructure which models
// a map field of a message
type MapFieldElement struct {
	Parent FProtoElement
	*FieldElement
	KeyType string
}

func (f *MapFieldElement) FieldName() string {
	return f.Name
}

func (f *MapFieldElement) FirstFieldTag() int {
	return f.Tag
}

// OneOfElement is a datastructure which models
// a oneoff construct in a protobuf file. All the fields in a
// oneof construct share memory, and at most one field can be
// set at any time.
type OneOfFieldElement struct {
	Parent  FProtoElement
	Name    string
	Comment *Comment
	Options []*OptionElement
	Fields  []FieldElementTag
}

func (f *OneOfFieldElement) FieldName() string {
	return f.Name
}

func (f *OneOfFieldElement) FirstFieldTag() int {
	smallest := -1
	for _, fld := range f.Fields {
		if smallest == -1 || fld.FirstFieldTag() < smallest {
			smallest = fld.FirstFieldTag()
		}
	}
	return smallest
}

// ExtensionsElement is a datastructure which models
// an extensions construct in a protobuf file. An extension
// is a placeholder for a field whose type is not defined by the
// original .proto file. This allows other .proto files to add
// to the original message definition by defining field ranges which
// can be used for extensions.
type ExtensionsElement struct {
	Parent  FProtoElement
	Comment *Comment
	Start   int
	End     int
	IsMax   bool
}

// ReservedRangeElement is a datastructure which models
// a reserved construct in a protobuf message.
type ReservedRangeElement struct {
	Parent  FProtoElement
	Comment *Comment
	Start   int
	End     int
	IsMax   bool
}

// MessageElement is a datastructure which models
// the message construct in a protobuf file.
type MessageElement struct {
	Parent         FProtoElement
	Name           string
	Comment        *Comment
	IsExtend       bool
	Options        []*OptionElement
	Fields         []FieldElementTag
	Enums          []*EnumElement
	Messages       []*MessageElement
	ExtendMessages []*MessageElement
	Extensions     []*ExtensionsElement
	ReservedRanges []*ReservedRangeElement
	ReservedNames  []string
}

// ProtoFile is a datastructure which represents the parsed model
// of the given protobuf file.
//
// It includes the package name, the syntax, the import dependencies,
// any public import dependencies, any options, enums, messages, services,
// extension declarations etc.
//
// This is populated by the parser and returned to the
// client code.
type ProtoFile struct {
	PackageName        string
	Syntax             string
	Dependencies       []string
	PublicDependencies []string
	WeakDependencies   []string
	Options            []*OptionElement
	Enums              []*EnumElement
	Messages           []*MessageElement
	ExtendMessages     []*MessageElement
	Services           []*ServiceElement
}

// Tag interfaces

func (e *OptionElement) FProtoElement()        {}
func (e *EnumConstantElement) FProtoElement()  {}
func (e *EnumElement) FProtoElement()          {}
func (e *RPCElement) FProtoElement()           {}
func (e *ServiceElement) FProtoElement()       {}
func (e *FieldElement) FProtoElement()         {}
func (e *MapFieldElement) FProtoElement()      {}
func (e *OneOfFieldElement) FProtoElement()    {}
func (e *ExtensionsElement) FProtoElement()    {}
func (e *ReservedRangeElement) FProtoElement() {}
func (e *MessageElement) FProtoElement()       {}
func (e *ProtoFile) FProtoElement()            {}

func (e *OptionElement) ParentElement() FProtoElement        { return e.Parent }
func (e *EnumConstantElement) ParentElement() FProtoElement  { return e.Parent }
func (e *EnumElement) ParentElement() FProtoElement          { return e.Parent }
func (e *RPCElement) ParentElement() FProtoElement           { return e.Parent }
func (e *ServiceElement) ParentElement() FProtoElement       { return e.Parent }
func (e *FieldElement) ParentElement() FProtoElement         { return e.Parent }
func (e *MapFieldElement) ParentElement() FProtoElement      { return e.Parent }
func (e *OneOfFieldElement) ParentElement() FProtoElement    { return e.Parent }
func (e *ExtensionsElement) ParentElement() FProtoElement    { return e.Parent }
func (e *ReservedRangeElement) ParentElement() FProtoElement { return e.Parent }
func (e *MessageElement) ParentElement() FProtoElement       { return e.Parent }
func (e *ProtoFile) ParentElement() FProtoElement            { return nil }

func (e *OptionElement) ElementName() string        { return e.Name }
func (e *EnumConstantElement) ElementName() string  { return e.Name }
func (e *EnumElement) ElementName() string          { return e.Name }
func (e *RPCElement) ElementName() string           { return e.Name }
func (e *ServiceElement) ElementName() string       { return e.Name }
func (e *FieldElement) ElementName() string         { return e.Name }
func (e *MapFieldElement) ElementName() string      { return e.Name }
func (e *OneOfFieldElement) ElementName() string    { return e.Name }
func (e *ExtensionsElement) ElementName() string    { return "" }
func (e *ReservedRangeElement) ElementName() string { return "" }
func (e *MessageElement) ElementName() string       { return e.Name }
func (e *ProtoFile) ElementName() string            { return "" }

func (e *OptionElement) ElementTypeName() string        { return "OPTION" }
func (e *EnumConstantElement) ElementTypeName() string  { return "ENUM CONSTANT" }
func (e *EnumElement) ElementTypeName() string          { return "ENUM" }
func (e *RPCElement) ElementTypeName() string           { return "RPC" }
func (e *ServiceElement) ElementTypeName() string       { return "SERVICE" }
func (e *FieldElement) ElementTypeName() string         { return "FIELD" }
func (e *MapFieldElement) ElementTypeName() string      { return "MAP FIELD" }
func (e *OneOfFieldElement) ElementTypeName() string    { return "ONEOF FIELD" }
func (e *ExtensionsElement) ElementTypeName() string    { return "EXTENSION" }
func (e *ReservedRangeElement) ElementTypeName() string { return "RESERVED RANGE" }
func (e *MessageElement) ElementTypeName() string       { return "MESSAGE" }
func (e *ProtoFile) ElementTypeName() string            { return "PROTO FILE" }
