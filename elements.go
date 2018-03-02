package fproto

// OptionElement is a datastructure which models
// the option construct in a protobuf file. Option constructs
// exist at various levels/contexts like file, message etc.
type OptionElement struct {
	Parent interface{}
	Name   string
	Value  string
	//IsParenthesized bool
}

// EnumConstantElement is a datastructure which models
// the fields within an enum construct. Enum constants can
// also have inline options specified.
type EnumConstantElement struct {
	Parent        interface{}
	Name          string
	Documentation string
	Options       []*OptionElement
	Tag           int
}

// EnumElement is a datastructure which models
// the enum construct in a protobuf file. Enums are
// defined standalone or as nested entities within messages.
type EnumElement struct {
	Parent        interface{}
	Name          string
	Documentation string
	Options       []*OptionElement
	EnumConstants []*EnumConstantElement
}

// RPCElement is a datastructure which models
// the rpc construct in a protobuf file. RPCs are defined
// nested within ServiceElements.
type RPCElement struct {
	Parent          interface{}
	Name            string
	Documentation   string
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
	Parent        interface{}
	Name          string
	Documentation string
	Options       []*OptionElement
	RPCs          []*RPCElement
}

type FieldElementTag interface {
	FieldName() string
}

// FieldElement is a datastructure which models
// a field of a message, a field of a oneof element
// or an entry in the extend declaration in a protobuf file.
type FieldElement struct {
	Parent        interface{}
	Name          string
	Documentation string
	Options       []*OptionElement
	Repeated      bool
	Optional      bool // proto2
	Required      bool // proto2
	Type          string
	Tag           int
}

func (f *FieldElement) FieldName() string {
	return f.Name
}

// MapFieldElement is a datastructure which models
// a map field of a message
type MapFieldElement struct {
	Parent interface{}
	*FieldElement
	KeyType string
}

func (f *MapFieldElement) FieldName() string {
	return f.Name
}

// OneOfElement is a datastructure which models
// a oneoff construct in a protobuf file. All the fields in a
// oneof construct share memory, and at most one field can be
// set at any time.
type OneOfElement struct {
	Parent        interface{}
	Name          string
	Documentation string
	Options       []*OptionElement
	Fields        []FieldElementTag
}

func (f *OneOfElement) FieldName() string {
	return f.Name
}

// ExtensionsElement is a datastructure which models
// an extensions construct in a protobuf file. An extension
// is a placeholder for a field whose type is not defined by the
// original .proto file. This allows other .proto files to add
// to the original message definition by defining field ranges which
// can be used for extensions.
type ExtensionsElement struct {
	Parent        interface{}
	Documentation string
	Start         int
	End           int
	IsMax         bool
}

// ReservedRangeElement is a datastructure which models
// a reserved construct in a protobuf message.
type ReservedRangeElement struct {
	Parent        interface{}
	Documentation string
	Start         int
	End           int
	IsMax         bool
}

// MessageElement is a datastructure which models
// the message construct in a protobuf file.
type MessageElement struct {
	Parent         interface{}
	Name           string
	Documentation  string
	IsExtend       bool
	Options        []*OptionElement
	Fields         []FieldElementTag
	Enums          []*EnumElement
	Messages       []*MessageElement
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
	Services           []*ServiceElement
}
