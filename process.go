package fproto

import "strings"

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

func (f *ProtoFile) FindName(name string) []interface{} {
	ret := make([]interface{}, 0)

	nfirst, nrest := NameSplit(name)

	// items that cannot nest
	if nrest != "" {
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

func (f *MessageElement) FindName(name string) []interface{} {
	ret := make([]interface{}, 0)

	nfirst, nrest := NameSplit(name)

	// items that cannot nest
	if nrest != "" {
		for _, el := range f.Enums {
			if el.Name == nfirst {
				ret = append(ret, el)
			}
		}
		for _, el := range f.Fields {
			if el.Name == nfirst {
				ret = append(ret, el)
			}
		}
		for _, el := range f.MapFields {
			if el.Name == nfirst {
				ret = append(ret, el)
			}
		}
		for _, el := range f.OneOfs {
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
