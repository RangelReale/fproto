package fdep

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/RangelReale/fproto"
)

type Dep struct {
	Files    map[string]*FileDep
	Packages map[string][]string
}

func NewDep() *Dep {
	return &Dep{
		Files:    make(map[string]*FileDep),
		Packages: make(map[string][]string),
	}
}

func (d *Dep) AddPath(dir string) error {
	return d.AddPathWithRoot("", dir)
}

func (d *Dep) AddPathWithRoot(currentpath, dir string) error {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, f := range files {
		if f.IsDir() {
			err = d.AddPathWithRoot(path.Join(currentpath, f.Name()), filepath.Join(dir, f.Name()))
		} else {
			if filepath.Ext(f.Name()) == ".proto" {
				err = d.AddFile(currentpath, filepath.Join(dir, f.Name()))
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

func (d *Dep) AddFile(currentpath string, filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("Error parsing file %s: %v", filename, err)
	}
	defer file.Close()

	pfile, err := fproto.Parse(file)
	if err != nil {
		return fmt.Errorf("Error parsing file %s: %v", filename, err)
	}

	fpath := path.Join(currentpath, filepath.Base(filename))
	d.Files[fpath] = &FileDep{
		Dep:       d,
		ProtoFile: pfile,
	}
	d.addPackage(fpath)
	return nil
}

func (d *Dep) addPackage(filepath string) {
	pkg := d.Files[filepath].ProtoFile.PackageName
	if _, ok := d.Packages[pkg]; !ok {
		d.Packages[pkg] = []string{}
	}

	d.Packages[pkg] = append(d.Packages[pkg], filepath)
}

func (d *Dep) GetType(name string) (*DepType, error) {
	t, err := d.GetTypes(name)
	if err != nil {
		return nil, err
	}

	if len(t) > 1 {
		return nil, fmt.Errorf("More than one type found for '%s'", name)
	}

	return t[0], nil
}

func (d *Dep) GetTypes(name string) ([]*DepType, error) {
	ret := make([]*DepType, 0)

	pkgs := make(map[string]string)

	sp := strings.Split(name, ".")

	for spi, _ := range sp {
		p := strings.Join(sp[:spi], ".")
		if _, ok := d.Packages[p]; ok {
			pkgs[p] = strings.Join(sp[spi:], ".")
		}
	}

	if len(pkgs) == 0 {
		return nil, fmt.Errorf("Package for type '%s' not found", name)
	}

	for sppkg, spname := range pkgs {
		for _, f := range d.Packages[sppkg] {
			for _, t := range d.Files[f].ProtoFile.FindName(spname) {
				ret = append(ret, &DepType{
					PackageName: sppkg,
					Name:        spname,
					Item:        t,
				})
			}
		}
	}

	return ret, nil
}
