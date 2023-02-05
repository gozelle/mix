package cmdGenerate

import (
	"github.com/gozelle/cobra"
)

var sentryCmd = &cobra.Command{
	Use:   "sentry",
	Short: "文件生成 sentry 文件",
	Long:  ``,
	Run:   generatesentry,
}

var (
	sentryTpl       string
	sentryPath      string
	sentryInterface string
	sentryOutfile   string
)

func init() {
	sentryCmd.Flags().StringVar(&sentryPath, "path", "", "[必填]源目录")
	sentryCmd.Flags().StringVar(&sentryInterface, "interface", "", "[必填]源 Interface 名")
	sentryCmd.Flags().StringVar(&sentryOutfile, "outfile", "", "[必填]生成文件路径")
	sentryCmd.Flags().StringVar(&sentryTpl, "template", "", "[可选] sentry 文件模板")
	err := sentryCmd.MarkFlagsRequired("path", "interface", "outfile")
	if err != nil {
		panic(err)
	}
}

func generatesentry(cmd *cobra.Command, args []string) {

	//pwd, err := os.Getwd()
	//if err != nil {
	//	commands.Fatal(err)
	//}
	//
	//sentryPath = fs.Join(pwd, sentryPath)
	//sentryOutfile = fs.Join(pwd, sentryOutfile)
	//if sentryTpl != "" {
	//	sentryTpl = fs.Join(pwd, sentryTpl)
	//}
	//
	//doc, err := parser.Parse(sentryTpl, sentryPath, sentryInterface)
	//if err != nil {
	//	commands.Fatal(err)
	//}
	//
	//c, err := doc.MarshalJSON()
	//if err != nil {
	//	commands.Fatal(err)
	//}
	//
	//err = fs.Write(sentryOutfile, c)
	//if err != nil {
	//	commands.Fatal(err)
	//}
	//
	//commands.Info("write file: %s", sentryOutfile)
}
