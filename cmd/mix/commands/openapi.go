package commands

import "github.com/gozelle/cobra"

var openapiCmd = &cobra.Command{
	Use:   "openapi",
	Short: "基于 OpenAPI 文件生成 SDK",
	Long:  ``,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}
