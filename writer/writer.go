package writer

import "github.com/gozelle/fs"

type Writer interface {
	Write(file, content string) (err error)
}

type FileWriter struct {
}

func (f FileWriter) Write(file, content string) error {
	return fs.Write(file, []byte(content))
}
