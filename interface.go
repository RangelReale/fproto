package fproto

//
// Tag interface
//

type FProtoElement interface {
	FProtoElement()
}

//
// Internal interfaces to help build the structs.
//

type iAddOption interface {
	addOptionElement(e *OptionElement)
}

type iAddEnumConstant interface {
	addEnumConstantElement(e *EnumConstantElement)
}

type iAddEnum interface {
	addEnumElement(e *EnumElement)
}

type iAddRPC interface {
	addRPCElement(e *RPCElement)
}

type iAddService interface {
	addServiceElement(e *ServiceElement)
}

type iAddField interface {
	addField(e FieldElementTag)
}

type iAddExtensions interface {
	addExtensionsElement(e *ExtensionsElement)
}

type iAddReservedRange interface {
	addReservedRangeElement(e *ReservedRangeElement)
}

type iAddReservedName interface {
	addReservedName(e string)
}

type iAddMessage interface {
	addMessageElement(e *MessageElement)
}

type iAddDependency interface {
	addDependency(e string)
}

type iAddPublicDependency interface {
	addPublicDependency(e string)
}

type iAddWeakDependency interface {
	addWeakDependency(e string)
}

//
// EnumConstantElement
//

func (el *EnumConstantElement) addOptionElement(e *OptionElement) {
	el.Options = append(el.Options, e)
}

//
// EnumElement
//

func (el *EnumElement) addOptionElement(e *OptionElement) {
	el.Options = append(el.Options, e)
}

func (el *EnumElement) addEnumConstantElement(e *EnumConstantElement) {
	el.EnumConstants = append(el.EnumConstants, e)
}

//
// RPCElement
//

func (el *RPCElement) addOptionElement(e *OptionElement) {
	el.Options = append(el.Options, e)
}

//
// ServiceElement
//

func (el *ServiceElement) addOptionElement(e *OptionElement) {
	el.Options = append(el.Options, e)
}

func (el *ServiceElement) addRPCElement(e *RPCElement) {
	el.RPCs = append(el.RPCs, e)
}

//
// FieldElement
//

func (el *FieldElement) addOptionElement(e *OptionElement) {
	el.Options = append(el.Options, e)
}

//
// MapFieldElement
//

func (el *MapFieldElement) addOptionElement(e *OptionElement) {
	el.Options = append(el.Options, e)
}

//
// OneOfElement
//

func (el *OneOfElement) addOptionElement(e *OptionElement) {
	el.Options = append(el.Options, e)
}

func (el *OneOfElement) addField(e FieldElementTag) {
	el.Fields = append(el.Fields, e)
}

//
// MessageElement
//

func (el *MessageElement) addOptionElement(e *OptionElement) {
	el.Options = append(el.Options, e)
}

func (el *MessageElement) addField(e FieldElementTag) {
	el.Fields = append(el.Fields, e)
}

func (el *MessageElement) addEnumElement(e *EnumElement) {
	el.Enums = append(el.Enums, e)
}

func (el *MessageElement) addMessageElement(e *MessageElement) {
	el.Messages = append(el.Messages, e)
}

func (el *MessageElement) addExtensionsElement(e *ExtensionsElement) {
	el.Extensions = append(el.Extensions, e)
}

func (el *MessageElement) addReservedRangeElement(e *ReservedRangeElement) {
	el.ReservedRanges = append(el.ReservedRanges, e)
}

func (el *MessageElement) addReserverName(e string) {
	el.ReservedNames = append(el.ReservedNames, e)
}

//
// ProtoFile
//

func (el *ProtoFile) addDependency(e string) {
	el.Dependencies = append(el.Dependencies, e)
}

func (el *ProtoFile) addPublicDependency(e string) {
	el.PublicDependencies = append(el.PublicDependencies, e)
}

func (el *ProtoFile) addWeakDependency(e string) {
	el.WeakDependencies = append(el.WeakDependencies, e)
}

func (el *ProtoFile) addOptionElement(e *OptionElement) {
	el.Options = append(el.Options, e)
}

func (el *ProtoFile) addEnumElement(e *EnumElement) {
	el.Enums = append(el.Enums, e)
}

func (el *ProtoFile) addMessageElement(e *MessageElement) {
	el.Messages = append(el.Messages, e)
}

func (el *ProtoFile) addServiceElement(e *ServiceElement) {
	el.Services = append(el.Services, e)
}
