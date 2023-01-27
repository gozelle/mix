package generateCmd

import (
	"fmt"
	"github.com/gozelle/cobra"
	"github.com/gozelle/mix/cmd/mix/commands"
	"os"
	"path/filepath"
)

var sdkCmd = &cobra.Command{
	Use:   "sdk",
	Short: "生成 API SDK",
	Long:  ``,
	Run:   generateSDK,
}

var (
	sdkOpenapi string
	sdkType    string
	sdkOutdir  string
	sdkOptions string
)

func init() {
	sdkCmd.Flags().StringVar(&sdkOpenapi, "openapi", "", "[必填] OpenAPI 文件")
	sdkCmd.Flags().StringVar(&sdkType, "sdk", "", "[必填] SDK 类型，如：axios")
	sdkCmd.Flags().StringVar(&sdkOutdir, "outdir", "", "[必填] SDK 存放目录")
	sdkCmd.Flags().StringVar(&sdkOptions, "options", "", "[可选]配置参数，请查看不同 SDK 配置选项")
	err := sdkCmd.MarkFlagsRequired("openapi", "sdk", "outdir")
	if err != nil {
		panic(err)
	}
}

func generateSDK(cmd *cobra.Command, args []string) {
	pwd, err := os.Getwd()
	if err != nil {
		commands.Fatal(err)
	}
	sdkOpenapi = filepath.Join(pwd, sdkOpenapi)
	sdkOutdir = filepath.Join(pwd, sdkOutdir)
	
	switch sdkType {
	case "axios":
	
	default:
		commands.Fatal(fmt.Errorf("sdk type: %s unsupported", sdkType))
	}
}
