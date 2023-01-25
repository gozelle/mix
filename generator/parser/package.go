package parser

import (
	"fmt"
	"github.com/gozelle/fs"
	"github.com/gozelle/logging"
	"strings"
)

var log = logging.Logger("parser")

type Packages map[string]*Package

type Package struct {
	Name       string
	Path       string
	Interfaces map[string]*Interface
	Defs       map[string]*Def
	Files      []*File
}

func (p *Package) getDef(name string) *Def {
	if p.Defs == nil {
		return nil
	}
	return p.Defs[name]
}

func (p *Package) addType(name string, item *Def) *Def {
	if p.Defs == nil {
		p.Defs = map[string]*Def{}
	}
	if v, ok := p.Defs[name]; ok {
		return v
	}
	p.Defs[name] = item
	return item
}

func (p *Package) addInterface(name string, item *Interface) *Interface {
	if p.Interfaces == nil {
		p.Interfaces = map[string]*Interface{}
	}
	if v, ok := p.Interfaces[name]; ok {
		return v
	}
	p.Interfaces[name] = item
	return item
}

func (p *Package) GetInterface(name string) *Interface {
	if p.Interfaces == nil {
		return nil
	}
	v, ok := p.Interfaces[name]
	if ok {
		return v
	}
	return nil
}

func (p *Package) loadFiles(mod *Mod, files []string) (err error) {
	for _, v := range files {
		if !strings.HasSuffix(v, "_test.go") {
			f := &File{pkg: p, mod: mod, path: v}
			err = f.load(v)
			if err != nil {
				return
			}
		}
	}
	return
}

func (p *Package) load(mod *Mod, dir string) error {
	
	if mod.loaded == nil {
		mod.loaded = map[string]bool{}
	}
	if _, ok := mod.loaded[dir]; ok {
		return nil
	}
	mod.loaded[dir] = true
	//log.Debugf("load dir: %s", dir)
	err := fs.IsDir(dir)
	if err != nil {
		return fmt.Errorf("only accept dir")
	}
	files, err := fs.Files(dir, ".go")
	if err != nil {
		return err
	}
	
	err = p.loadFiles(mod, files)
	if err != nil {
		return err
	}
	
	return nil
}
