package writter

import (
	"github.com/gozelle/fs"
)

var (
	GenGoFileSuffix = ".mix.go"
)

type File struct {
	Name    string
	Content string
}

func WriteFiles(dir string, files []*File) (paths []string, err error) {
	for _, v := range files {
		if v.Name == "" {
			continue
		}
		path := fs.Join(dir, v.Name)
		err = fs.Write(path, []byte(v.Content))
		if err != nil {
			return
		}
		paths = append(paths, path)
	}
	return
}
