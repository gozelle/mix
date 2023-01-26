package parser

import "fmt"

func Parse(mod *Mod, dir string) (pkg *Package, err error) {
	
	pkg = &Package{}
	err = pkg.Parse(mod, dir)
	if err != nil {
		return
	}
	
	for _, v := range pkg.interfaces {
		err = v.Load(pkg, v.file)
		if err != nil {
			panic(fmt.Errorf("load interface: %s error: %s", v.Name, err))
		}
	}
	
	//for name, v := range pkg.defs {
	//log.Infof("def: %s  => %s", name, v.String())
	//}
	//}
	
	return
}
