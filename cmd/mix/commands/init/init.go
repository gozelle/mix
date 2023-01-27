package initCmd

import "github.com/gozelle/cobra"

var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "初始化项目脚手架",
	Long:  ``,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}
