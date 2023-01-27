package commands

import (
	"github.com/gozelle/cobra"
	jsonrpc_client "github.com/gozelle/mix/generator/clients/jsonrpc-client"
	"os"
	"path/filepath"
	"strings"
)

var clientExamples = cobra.Examples{
	{
		Usage:   "mix generate client --path ./example --pkg example",
		Comment: "简单用法",
	},
	{
		Usage:   "mix generate client --path ./example --pkg example --outpkg example --outfile ./example/proxy_gen.go",
		Comment: "自定义路径",
	},
}

var clientCmd = &cobra.Command{
	Use:     "client",
	Short:   "生成 RPC 客户端",
	Long:    ``,
	Example: clientExamples.String(),
	Run:     generateClient,
}

var (
	clientPath    string
	clientPkg     string
	clientOutPkg  string
	clientOutfile string
)

func init() {
	clientCmd.Flags().StringVar(&clientPath, "path", "", "[必填]源路径")
	clientCmd.Flags().StringVar(&clientPkg, "pkg", "", "[必填]源包名")
	clientCmd.Flags().StringVar(&clientOutPkg, "outpkg", "", "[可选]指定生成 package 名")
	clientCmd.Flags().StringVar(&clientOutfile, "outfile", "", "[可选]指定存放生成文件路径")
	err := clientCmd.MarkFlagsRequired("path", "pkg")
	if err != nil {
		panic(err)
	}
}

func generateClient(cmd *cobra.Command, args []string) {
	if clientOutPkg == "" {
		clientOutPkg = clientPkg
		warning("modify outpkg: %s", clientPkg)
	}
	pwd, err := os.Getwd()
	if err != nil {
		fatal(err)
	}

	clientPath = filepath.Join(pwd, clientPath)

	if clientOutfile == "" {
		clientOutfile = filepath.Join(clientPath, "proxy_gen.go")
		warning("modify outfile: %s", clientOutfile)
	} else if !strings.HasSuffix(clientOutfile, ".go") {
		clientOutfile = filepath.Join(pwd, clientOutfile, "proxy_gen.go")
		warning("modify outfile: %s", clientOutfile)
	} else if !strings.HasSuffix(clientOutfile, "gen_.go") {
		clientOutfile = filepath.Join(pwd, strings.TrimSuffix(clientOutfile, ".go")+"_gen.go")
		warning("modify outfile: %s", clientOutfile)
	} else {
		clientOutfile = filepath.Join(pwd, clientOutfile)
	}

	err = jsonrpc_client.Generate(clientPath, clientPkg, clientOutPkg, clientOutfile)
	if err != nil {
		fatal(err)
	}
	info("write file: %s", clientOutfile)
}
