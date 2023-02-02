package cmdInit

import "github.com/gozelle/cobra"

var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "初始化项目脚手架",
	Long:  ``,
	Run:   initProject,
}

type genFile struct {
	name   string
	tpl    string
	dir    string
	render func(project string) interface{}
}
