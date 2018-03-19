package fdep

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	gofilepath "path/filepath"
	"strings"

	"github.com/RangelReale/fproto"
)

// Dep represents an .proto file hierarchy with dependencies between files.
type Dep struct {
	// List of files parsed. The file names are the INTERNAL name, like "google/protobuf/empty.proto".
	Files map[string]*FileDep

	// List of packages of the parsed files, with a list of files if a package have more than one.
	// The list of files should be used on the Files member to find the file itself.
	Packages map[string][]string

	// Extensions for a given type. Each item contains a package name.
	Extensions map[string][]string

	// Directories to look for unknown includes
	IncludeDirs []string

	// Ignore not found dependencies
	IgnoreNotFoundDependencies bool
}

// Creates a new Dep struct.
func NewDep() *Dep {
	return &Dep{
		Files:      make(map[string]*FileDep),
		Packages:   make(map[string][]string),
		Extensions: make(map[string][]string),
	}
}

// Add one include dir to be searched for an unknown import.
func (d *Dep) AddIncludeDir(dir string) error {
	if s, err := os.Stat(dir); err != nil {
		return fmt.Errorf("Invalid directory %s: %v", dir, err)
	} else if !s.IsDir() {
		return fmt.Errorf("Path %s isn't a directory", dir)
	}

	d.IncludeDirs = append(d.IncludeDirs, dir)

	return nil
}

// Returns a FileDep given a ProtoFile
func (d *Dep) FileDepFromProtofile(pfile *fproto.ProtoFile) *FileDep {
	for _, fd := range d.Files {
		if fd.ProtoFile == pfile {
			return fd
		}
	}
	return nil
}

// Returns a FileDep given an fproto element
func (d *Dep) FileDepFromElement(element fproto.FProtoElement) *FileDep {
	root := fproto.GetRootElement(element)
	if pfile, pfileok := root.(*fproto.ProtoFile); pfileok {
		return d.FileDepFromProtofile(pfile)
	}
	return nil
}

// Returns a DepType given an fproto element
func (d *Dep) DepTypeFromElement(element fproto.FProtoElement) *DepType {
	fd := d.FileDepFromElement(element)
	if fd != nil {
		return NewDepTypeFromElement(fd, element)
	}
	return nil
}

// Add files from one directory recursively, assuming this is a .protobuf root path.
// Ex: dep.AddPath("/protoc-3.5.1/include", fdep.DepType_Imported)
// This will add files from google/protobuf directory.
func (d *Dep) AddPath(dir string, deptype FileDepType) error {
	return d.AddPathWithRoot("", dir, deptype)
}

// Add files from one directory recursively, using "currentpath" as the root path of this directory.
// Ex: dep.AddPathWithRoot("google", "/protoc-3.5.1/include/google", fdep.DepType_Imported)
// This will add files from protobuf directory, assuming they are on the "google" path.
//
// These 2 commands have exactly the same effect:
// 		dep.AddPath("/protoc-3.5.1/include", fdep.DepType_Imported)
// 		dep.AddPathWithRoot("google", "/protoc-3.5.1/include/google", fdep.DepType_Imported)
func (d *Dep) AddPathWithRoot(currentpath, dir string, deptype FileDepType) error {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, f := range files {
		if f.IsDir() {
			err = d.AddPathWithRoot(path.Join(currentpath, f.Name()), gofilepath.Join(dir, f.Name()), deptype)
		} else {
			if gofilepath.Ext(f.Name()) == ".proto" {
				err = d.AddFile(currentpath, gofilepath.Join(dir, f.Name()), deptype)
			} else {
				err = nil
			}
		}
		if err != nil {
			return err
		}
	}

	return nil
}

// Adds a single file to the dependency, assuming the file's path as "currentpath".
// Ex: dep.AddFile("google/protobuf", "/protoc-3.5.1/include/google/protobuf/empty.proto", fdep.DepType_Imported)
func (d *Dep) AddFile(currentpath string, filename string, deptype FileDepType) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("Error parsing file %s: %v", filename, err)
	}
	defer file.Close()

	// builds the file path
	fpath := path.Join(currentpath, gofilepath.Base(filename))

	// reads the file
	return d.AddReader(fpath, file, deptype)
}

// Adds a single file to the dependency, using an reader.
// Ex: dep.AddReader("google/protobuf/Empty.proto", reader, fdep.DepType_Imported)
func (d *Dep) AddReader(filepath string, r io.Reader, deptype FileDepType) error {
	// parses the file
	pfile, err := fproto.Parse(r)
	if err != nil {
		return fmt.Errorf("Error parsing file %s: %v", filepath, err)
	}

	// adds the file to the list
	d.Files[filepath] = &FileDep{
		FilePath:  filepath,
		DepType:   deptype,
		Dep:       d,
		ProtoFile: pfile,
	}

	// load file dependencies
	for _, fd := range pfile.Dependencies {
		err = d.AddIncludeFile(fd)
		if err != nil {
			return err
		}
	}

	// add to the package list
	d.addPackage(filepath)

	// add to the extension list
	d.addExtensions(filepath)

	return nil
}

// Adds an include file
func (d *Dep) AddIncludeFile(filepath string) error {
	if _, ok := d.Files[filepath]; ok {
		// File already exists
		return nil
	}

	for _, inc := range d.IncludeDirs {
		inc_file := gofilepath.Join(inc, gofilepath.FromSlash(filepath))

		// check if file exists
		_, err := os.Stat(inc_file)
		if err != nil && !os.IsNotExist(err) {
			return err
		} else if err == nil {
			return d.AddFile(path.Dir(filepath), inc_file, DepType_Imported)
		}
	}

	if !d.IgnoreNotFoundDependencies {
		return fmt.Errorf("File not found in include path: %s", filepath)
	}

	// Add file as if it was found, but without a parsed file and without package references
	d.Files[filepath] = &FileDep{
		FilePath:  filepath,
		DepType:   DepType_Imported,
		Dep:       d,
		ProtoFile: nil,
	}

	return nil
}

// Adds files from a provider
func (d *Dep) AddFileProvider(fp FileProvider) error {
	for fp.HasNext() {
		filepath, r, deptype, err := fp.GetNext()
		if err != nil {
			return err
		}

		err = d.AddReader(filepath, r, deptype)
		if err != nil {
			return err
		}
	}
	return nil
}

// Adds the package of the file to the Packages list.
func (d *Dep) addPackage(filepath string) {
	pkg := d.Files[filepath].ProtoFile.PackageName
	if _, ok := d.Packages[pkg]; !ok {
		d.Packages[pkg] = []string{}
	}

	d.Packages[pkg] = append(d.Packages[pkg], filepath)
}

// Add message extensions
func (d *Dep) addExtensions(filepath string) {
	d.addMessageExtensions(d.Files[filepath].ProtoFile, d.Files[filepath].ProtoFile.Messages)
}

// Add message extensions
func (d *Dep) addMessageExtensions(prfile *fproto.ProtoFile, messages []*fproto.MessageElement) {
	for _, m := range messages {
		if m.IsExtend {
			if _, ok := d.Extensions[m.Name]; !ok {
				d.Extensions[m.Name] = make([]string, 0)
			}
			d.Extensions[m.Name] = append(d.Extensions[m.Name], prfile.PackageName)
		}

		d.addMessageExtensions(prfile, m.Messages)
	}
}

// Builds a list of valid package names from the dotted name.
// for example, if name = "google.protobuf.Empty", this will search for
// package "google", then "google.protobuf", but only "google.protobuf" will
// be found and added to the list.
//
// The map item value will contain the rest of the type name, in the example case,
// "Empty". It can also contain dots in case of nested items.
func (d *Dep) FindPackagesOfName(name string) map[string]string {
	pkgs := make(map[string]string)

	nameparts := strings.Split(name, ".")

	for nameidx, _ := range nameparts {
		p := strings.Join(nameparts[:nameidx], ".")
		if _, ok := d.Packages[p]; ok {
			pkgs[p] = strings.Join(nameparts[nameidx:], ".")
		}
	}

	return pkgs
}

// Returns one named type from the dependency.
//
// If multiple types are found for the same name, an error is issued.
// If there is this possibility, use the GetTypes method instead.
func (d *Dep) GetType(name string) (*DepType, error) {
	t, err := d.GetTypes(name)
	if err != nil {
		return nil, err
	}

	if len(t) > 1 {
		return nil, fmt.Errorf("More than one type found for '%s'", name)
	} else if len(t) == 0 {
		return nil, nil
	}

	return t[0], nil
}

// Like GetType, but returns an error if not found
func (d *Dep) MustGetType(name string) (*DepType, error) {
	t, err := d.GetType(name)
	if err != nil {
		return nil, err
	}
	if t == nil {
		return nil, fmt.Errorf("Type %s not found", name)
	}
	return t, nil
}

// Gets an extensions for a type from a source package
func (d *Dep) GetTypeExtension(name string, extensionPkg string) (*DepType, error) {
	t, err := d.GetType(name)
	if err != nil {
		return nil, err
	}
	if t == nil {
		return nil, nil
	}

	for _, ext := range t.ExtensionPackages() {
		if ext == extensionPkg {
			return d.GetType(fmt.Sprintf("%s.%s", ext, name))
		}
	}

	return nil, nil
}

// Returns all named types from the dependency.
//
// Use this method if there is a possibility that one name resolves to more than one type.
func (d *Dep) GetTypes(name string) ([]*DepType, error) {
	return d.internalGetTypes(name, nil)
}

// This functions is the one that really does the type finding.
// If filedep is not-nil, the type is returned in relation to it.
func (d *Dep) internalGetTypes(name string, filedep *FileDep) ([]*DepType, error) {
	ret := make([]*DepType, 0)

	// check if is scalar
	if scalar, is_scalar := fproto.ParseScalarType(name); is_scalar {
		ret = append(ret, NewDepTypeScalar(scalar))
	}

	// locate the name into the own filedep
	if filedep != nil {
		for _, t := range filedep.ProtoFile.FindName(name) {
			switch t.(type) {
			case fproto.FieldElementTag:
				// ignore fields
			default:
				ret = append(ret, NewDepType(filedep, "", filedep.OriginalAlias(), name, t))
			}
		}
	}

	pkgs := d.FindPackagesOfName(name)

	if len(pkgs) == 0 {
		if len(ret) > 0 {
			return ret, nil
		}

		return nil, nil
	}

	// Loop into the found packages.
	for sppkg, spname := range pkgs {
		// Loop into the files of these packages.
		for _, f := range d.Packages[sppkg] {
			include_file := false

			if filedep != nil {
				// If a file was passed, only check on the dependencies of the file.
				for _, ffdep := range filedep.ProtoFile.Dependencies {
					if ffdep == f {
						include_file = true
						break
					}
				}
			} else {
				// Else check all files
				include_file = true
			}

			if include_file {
				// Search the name on the current proto file.
				for _, t := range d.Files[f].ProtoFile.FindName(spname) {
					ret = append(ret, NewDepType(d.Files[f], sppkg, sppkg, spname, t))
				}
			}
		}
	}

	return ret, nil
}

// Gets a file of a name. Try all package names until a file is found.
// The type itself that may be on the name is ignored.
// For example, GetFilesOfName("google.protobuf.Empty") returns:
//		FileDep: *FileDep{"google/protobuf/empty.proto"}
//		Package: google.protobuf
//		Name: Empty
func (d *Dep) GetFileOfName(name string) (*FileDepOfName, error) {
	t, err := d.internalGetFilesOfName(name, nil)
	if err != nil {
		return nil, err
	}

	if len(t) > 1 {
		return nil, fmt.Errorf("More than one file found for '%s'", name)
	} else if len(t) == 0 {
		return nil, nil
	}

	return t[0], nil
}

// Gets the files of a name. Try all package names until a file is found.
// The type itself that may be on the name is ignored.
func (d *Dep) GetFilesOfName(name string) ([]*FileDepOfName, error) {
	return d.internalGetFilesOfName(name, nil)
}

// Gets the files of a name. Try all package names until a file is found.
// The type itself that may be on the name is ignored.
func (d *Dep) internalGetFilesOfName(name string, filedep *FileDep) ([]*FileDepOfName, error) {
	pkgs := d.FindPackagesOfName(name)

	if len(pkgs) == 0 {
		return nil, nil
	}

	found := make(map[string]*FileDepOfName)

	// Loop into the found packages.
	for sppkg, spname := range pkgs {
		// Loop into the files of these packages.
		for _, f := range d.Packages[sppkg] {
			include_file := false

			if filedep != nil {
				// If a file was passed, only check on the dependencies of the file.
				for _, ffdep := range filedep.ProtoFile.Dependencies {
					if ffdep == f {
						include_file = true
						break
					}
				}
			} else {
				// Else check all files
				include_file = true
			}

			if include_file {
				found[f] = &FileDepOfName{
					FileDep: d.Files[f],
					Package: sppkg,
					Name:    spname,
				}
			}
		}
	}

	// build response
	var ret []*FileDepOfName
	for _, fd := range found {
		ret = append(ret, fd)
	}

	return ret, nil
}

// Get a list for extension packages for a type.
func (d *Dep) GetExtensions(filedep *FileDep, originalAlias string, name string) []string {
	var ret []string

	fname := name
	if originalAlias != "" {
		fname = originalAlias + "." + fname
	}

	if ext, extok := d.Extensions[fname]; extok {
		for _, e := range ext {
			ret = append(ret, e)
		}
	}

	return ret
}

func (d *Dep) GetOption(optionItem OptionItem, name string) (*OptionType, error) {
	t, err := d.internalGetOptions(optionItem, name, nil)
	if err != nil {
		return nil, err
	}

	if len(t) > 1 {
		return nil, fmt.Errorf("More than one option found for '%s'", name)
	} else if len(t) == 0 {
		return nil, nil
	}

	return t[0], nil
}

func (d *Dep) GetOptions(optionItem OptionItem, name string) ([]*OptionType, error) {
	return d.internalGetOptions(optionItem, name, nil)
}

func (d *Dep) internalGetOptions(optionItem OptionItem, name string, filedep *FileDep) ([]*OptionType, error) {
	// Get packages of passed name
	depnames, err := d.internalGetFilesOfName(name, filedep)
	if err != nil {
		return nil, err
	}

	// find the source type from descriptor.proto
	srcTypeName := optionItem.MessageName()

	sourceType, err := d.GetType(srcTypeName)
	if err != nil {
		return nil, fmt.Errorf("Error gettint the source type '%s': %v", srcTypeName, err)
	}

	var ret []*OptionType
	for _, dn := range depnames {
		// checks if there is an extension message of the source type in the root of the proto file
		for _, m := range dn.FileDep.ProtoFile.Messages {
			if m.IsExtend && m.Name == srcTypeName {
				include_file := false

				var field_item fproto.FieldElementTag

				if dn.Name != "" {
					// checks if the message has the first part of the name
					field_item, _ = m.FindFieldPartial(dn.Name)
					if field_item != nil {
						include_file = true
					}
				} else {
					include_file = true
				}

				if include_file {
					ret = append(ret, &OptionType{
						OptionName:   name,
						SourceOption: sourceType,
						Option:       NewDepTypeFromElement(dn.FileDep, m),
						OptionItem:   optionItem,
						Name:         dn.Name,
						FieldItem:    field_item,
					})
				}
			}
		}
	}

	if len(ret) == 0 {
		// check if the field isn't on the core type
		fld, _ := sourceType.Item.(*fproto.MessageElement).FindFieldPartial(name)
		if fld != nil {
			ret = append(ret, &OptionType{
				SourceOption: sourceType,
				Option:       nil,
				OptionItem:   optionItem,
				Name:         name,
			})
		}
	}

	return ret, nil
}
