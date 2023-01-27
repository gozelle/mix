package example

import "fmt"

type Examples []Example

func (e Examples) String() string {
	a := 0
	for _, v := range e {
		if b := len(v.Usage); b > a {
			a = b
		}
	}
	r := ""
	for _, v := range e {
		f := fmt.Sprintf("%%-%ds%%s\n", a)
		if v.Comment != "" {
			v.Comment = fmt.Sprintf(" # %s", v.Comment)
		}
		r += fmt.Sprintf(f, v.Usage, v.Comment)
	}
	return r
}

type Example struct {
	Usage   string
	Comment string
}
