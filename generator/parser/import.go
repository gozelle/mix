package parser

type Import struct {
	Alias   string
	Path    string
	Package *Package
}
