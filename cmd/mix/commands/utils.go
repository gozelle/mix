package commands

import (
	"github.com/gozelle/color"
	"os"
)

func warning(format string, a ...any) {
	color.Yellow(format, a...)
}

func info(format string, a ...any) {
	color.Green(format, a...)
}

func fatal(err error) {
	color.Red(err.Error())
	os.Exit(1)
}
