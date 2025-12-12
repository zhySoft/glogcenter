package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"glc/my/global"
)

// rootCmd 表示调用时不带任何子命令的基本命令
var rootCmd = &cobra.Command{
	Use:     global.AppInfo.ServiceName,
	Short:   global.AppInfo.ServiceDisplayName,
	Long:    global.AppInfo.ServiceDesc,
	Example: "NotifyTask run --path=.",
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
		_ = cmd.Help()
	},
}

// Execute 执行 root 命令
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
