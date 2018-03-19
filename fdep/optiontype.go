package fdep

import "github.com/RangelReale/fproto"

type OptionItem int

const (
	FILE_OPTION OptionItem = iota
	MESSAGE_OPTION
	FIELD_OPTION
	ENUM_OPTION
	ENUMVALUE_OPTION
	SERVICE_OPTION
	METHOD_OPTION
)

func (ot OptionItem) MessageName() string {
	switch ot {
	case FILE_OPTION:
		return "google.protobuf.FileOptions"
	case MESSAGE_OPTION:
		return "google.protobuf.MessageOptions"
	case FIELD_OPTION:
		return "google.protobuf.FieldOptions"
	case ENUM_OPTION:
		return "google.protobuf.EnumOptions"
	case ENUMVALUE_OPTION:
		return "google.protobuf.EnumValueOptions"
	case SERVICE_OPTION:
		return "google.protobuf.ServiceOptions"
	case METHOD_OPTION:
		return "google.protobuf.MethodOptions"
	}
	return "Unknown"
}

type OptionType struct {
	// Requested option name
	OptionName string
	// Type of the root option
	SourceOption *DepType
	// Option can be nil if the type is one of the root option types
	Option *DepType
	// Type of option
	OptionItem OptionItem
	// Name of option after removing package name
	Name string
	// The proto field if available
	FieldItem fproto.FieldElementTag
}
