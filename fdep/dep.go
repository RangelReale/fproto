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
			ret = append(ret, NewDepType(filedep, "", filedep.OriginalAlias(), name, t))
		}
	}

	// builds a list of possible package names from the dotted name.
	// for example, if name = "google.protobuf.Empty", this will search for
	// package "google", then "google.protobuf", but only "google.protobuf" will
	// be found and added to the list.
	//
	// The map item value will contain the rest of the type name, in the example case,
	// "Empty". It can also contain dots in case of nested items.
	pkgs := make(map[string]string)

	sp := strings.Split(name, ".")

	for spi, _ := range sp {
		p := strings.Join(sp[:spi], ".")
		if _, ok := d.Packages[p]; ok {
			pkgs[p] = strings.Join(sp[spi:], ".")
		}
	}

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
func (d *Dep) GetFileOfName(name string) (*FileDep, error) {
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

func (d *Dep) GetFilesOfName(name string) ([]*FileDep, error) {
	return d.internalGetFilesOfName(name, nil)
}

// Gets the files of a name. Try all package names until a file is found.
// The type itself that may be on the name is ignored.
func (d *Dep) internalGetFilesOfName(name string, filedep *FileDep) ([]*FileDep, error) {
	// builds a list of possible package names from the dotted name.
	// for example, if name = "google.protobuf.Empty", this will search for
	// package "google", then "google.protobuf", but only "google.protobuf" will
	// be found and added to the list.
	//
	// The map item value will contain the rest of the type name, in the example case,
	// "Empty". It can also contain dots in case of nested items.
	pkgs := make(map[string]string)

	sp := strings.Split(name, ".")

	for spi, _ := range sp {
		p := strings.Join(sp[:spi], ".")
		if _, ok := d.Packages[p]; ok {
			pkgs[p] = strings.Join(sp[spi:], ".")
		}
	}

	if len(pkgs) == 0 {
		return nil, nil
	}

	found := make(map[string]*FileDep)

	// Loop into the found packages.
	for sppkg, _ := range pkgs {
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
				found[f] = d.Files[f]
				/*
					if spname == "" {
						found[f] = d.Files[f]
					} else {
						// Search the name on the current proto file.
						if len(d.Files[f].ProtoFile.FindName(spname)) > 0 {
							found[f] = d.Files[f]
						}
					}
				*/
			}
		}
	}

	// build response
	var ret []*FileDep
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
