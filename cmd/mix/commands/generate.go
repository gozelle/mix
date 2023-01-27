package commands

import (
	"github.com/gozelle/cobra"
)

var generateExamples = cobra.Examples{
	{
		Usage:   "mix generate client -h",
		Comment: "查看生成 Client 帮助",
	},
	{
		Usage:   "mix generate sdk -h",
		Comment: "查看生成 API SDK 帮助",
	},
}

var GenerateCmd = &cobra.Command{
	Use:     "generate",
	Example: generateExamples.String(),
	Short:   "生成 RPC Client、API SDK",
	Long:    ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

func init() {
	GenerateCmd.AddCommand(
		clientCmd,
		sdkCmd,
		openapiCmd,
	)
}
