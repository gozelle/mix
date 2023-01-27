package commands

import "github.com/gozelle/cobra"

var GenerateCmd = &cobra.Command{
	Use: "generate",
	Example: `mix generate client -h   # 查看生成 Client 帮助
mix generate sdk -h      # 查看生成 API SDK 帮助
`,
	Short: "生成 RPC Client、API SDK",
	Long:  ``,
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
