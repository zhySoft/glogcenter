package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"glc/my/global"
)

// versionCmd 表示 version 命令
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "版本号",
	Long:  "版本号",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(global.AppInfo.Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
