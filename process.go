package fproto

import (
	"strings"
)

// Parses the dot-separated string into the part before the first dot, and the part after it.
func NameSplit(name string) (first, rest string) {
	s := strings.Split(name, ".")
	if len(s) <= 0 {
		return "", ""
	} else if len(s) == 1 {
		return s[0], ""
	} else {
		return s[0], strings.Join(s[1:], ".")
	}
}

func ScopedName(element FProtoElement) string {
	return strings.Join(ScopedNameList(element), ".")
}

func ScopedNameList(element FProtoElement) []string {
	var ret []string
	cur := element
	for cur != nil {
		ename := cur.ElementName()
		if ename != "" {
			ret = append(ret, ename)
		}
		cur = cur.ParentElement()
	}

	// reverse the order
	return ReverseStr(ret)
}

func ScopedAlias(element FProtoElement) string {
	return strings.Join(ScopedAliasList(element), ".")
}

func ScopedAliasList(element FProtoElement) []string {
	var ret []string
	if element != nil {
		cur := element.ParentElement()
		for cur != nil {
			ename := cur.ElementName()
			if ename != "" {
				ret = append(ret, ename)
			}
			cur = cur.ParentElement()
		}
	}

	// reverse the order
	return ReverseStr(ret)
}

//
// PROCESS: ProtoFile
//

// Finds elements on the ProtoFile by name. Dots can be used to get an inner scope.
// Only Enum, Service and Message are searched.
// Ex: FindName("User.Address")
func (f *ProtoFile) FindName(name string) []FProtoElement {
	ret := make([]FProtoElement, 0)

	nfirst, nrest := NameSplit(name)

	// items that cannot nest
	if nrest == "" {
		for _, el := range f.Enums {
			if el.Name == nfirst {
				ret = append(ret, el)
			}
		}
		for _, el := range f.Services {
			if el.Name == nfirst {
				ret = append(ret, el)
			}
		}
	}

	// items that can nest
	for _, el := range f.Messages {
		if el.IsExtend && el.Name == name {
			ret = append(ret, el)
		} else if el.Name == nfirst {
			if nrest != "" {
				elr := el.FindName(nrest)
				if len(elr) > 0 {
					ret = append(ret, elr...)
				}
			} else {
				ret = append(ret, el)
			}
		}
	}

	return ret
}

// Finds an option by name.
func (f *ProtoFile) FindOption(name string) *OptionElement {
	for _, o := range f.Options {
		if o.Name == name || o.ParenthesizedName == name {
			return o
		}
	}
	return nil
}

func (f *ProtoFile) CollectEnums() []FProtoElement {
	var ret []FProtoElement

	for _, el := range f.Enums {
		ret = append(ret, el)
	}

	for _, el := range f.Messages {
		ret = append(ret, el.CollectEnums()...)
	}

	return ret
}

func (f *ProtoFile) CollectMessages() []FProtoElement {
	var ret []FProtoElement

	for _, el := range f.Messages {
		ret = append(ret, el)
		ret = append(ret, el.CollectMessages()...)
	}

	return ret
}

func (f *ProtoFile) CollectServices() []FProtoElement {
	var ret []FProtoElement

	for _, el := range f.Services {
		ret = append(ret, el)
	}

	return ret
}

//
// PROCESS: MessageElement
//

// Finds elements on the Message by name. Dots can be used to get an inner scope.
// Only Enum, Field, MapField, OneOf and inner Message are searched.
func (f *MessageElement) FindName(name string) []FProtoElement {
	ret := make([]FProtoElement, 0)

	nfirst, nrest := NameSplit(name)

	// items that cannot nest
	if nrest == "" {
		for _, el := range f.Enums {
			if el.Name == nfirst {
				ret = append(ret, el)
			}
		}
		for _, el := range f.Fields {
			if el.FieldName() == nfirst {
				ret = append(ret, el)
			}
		}
	}

	// items that can nest
	for _, el := range f.Messages {
		if el.Name == nfirst {
			if nrest != "" {
				elr := el.FindName(nrest)
				if len(elr) > 0 {
					ret = append(ret, elr...)
				}
			} else {
				ret = append(ret, el)
			}
		}
	}

	return ret
}

// Finds an option by name.
func (f *MessageElement) FindOption(name string) *OptionElement {
	for _, o := range f.Options {
		if o.Name == name || o.ParenthesizedName == name {
			return o
		}
	}
	return nil
}

func (f *MessageElement) CollectEnums() []FProtoElement {
	var ret []FProtoElement

	for _, el := range f.Enums {
		ret = append(ret, el)
	}

	for _, el := range f.Messages {
		ret = append(ret, el.CollectEnums()...)
	}

	return ret
}

func (f *MessageElement) CollectMessages() []FProtoElement {
	var ret []FProtoElement

	for _, el := range f.Messages {
		ret = append(ret, el)
		ret = append(ret, el.CollectMessages()...)
	}

	return ret
}

//
// PROCESS: OptionElement
//

// Finds an option by name.
func (f *OptionElement) FindOption(name string) *OptionElement {
	if f.Name == name || f.ParenthesizedName == name {
		return f
	}
	return nil
}

//
// PROCESS: EnumConstantElement
//

func (f *EnumConstantElement) FindOption(name string) *OptionElement {
	for _, o := range f.Options {
		if o.Name == name || o.ParenthesizedName == name {
			return o
		}
	}
	return nil
}

//
// PROCESS: EnumElement
//

func (f *EnumElement) FindOption(name string) *OptionElement {
	for _, o := range f.Options {
		if o.Name == name || o.ParenthesizedName == name {
			return o
		}
	}
	return nil
}

//
// PROCESS: RPCElement
//

func (f *RPCElement) FindOption(name string) *OptionElement {
	for _, o := range f.Options {
		if o.Name == name || o.ParenthesizedName == name {
			return o
		}
	}
	return nil
}

//
// PROCESS: ServiceElement
//

func (f *ServiceElement) FindOption(name string) *OptionElement {
	for _, o := range f.Options {
		if o.Name == name || o.ParenthesizedName == name {
			return o
		}
	}
	return nil
}

//
// PROCESS: FieldElement
//

func (f *FieldElement) FindOption(name string) *OptionElement {
	for _, o := range f.Options {
		if o.Name == name || o.ParenthesizedName == name {
			return o
		}
	}
	return nil
}

//
// PROCESS: MapFieldElement
//

func (f *MapFieldElement) FindOption(name string) *OptionElement {
	for _, o := range f.Options {
		if o.Name == name || o.ParenthesizedName == name {
			return o
		}
	}
	return nil
}

//
// PROCESS: OneOfElement
//

func (f *OneofFieldElement) FindOption(name string) *OptionElement {
	for _, o := range f.Options {
		if o.Name == name || o.ParenthesizedName == name {
			return o
		}
	}
	return nil
}

//
// PROCESS: ExtensionsElement
//

func (f *ExtensionsElement) FindOption(name string) *OptionElement {
	return nil
}

//
// PROCESS: ReservedRangeElement
//

func (f *ReservedRangeElement) FindOption(name string) *OptionElement {
	return nil
}
