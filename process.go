package fproto

import "strings"

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

// Finds an option on the ProtoFile by name.
func (f *ProtoFile) FindOption(name string) *OptionElement {
	for _, o := range f.Options {
		if o.Name == name {
			return o
		}
	}
	return nil
}

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
