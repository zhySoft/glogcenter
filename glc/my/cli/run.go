package cli

import (
	"glc/conf"
	"glc/ldb/backup"
	"glc/ldb/sysmnt"
	"glc/my/global"
	"glc/onstart"
	"glc/www/controller"

	"github.com/spf13/cobra"
)

// runCmd 表示 Run 命令
var runCmd = &cobra.Command{
	Use:     "run",
	Short:   "运行程序",
	Long:    "运行程序",
	Example: "path=.",
	Run: func(cmd *cobra.Command, args []string) {
		// 初始化全局变量
		global.Init(path)
		// 因为使用了配置文件，需要修改init顺序
		conf.Init()
		conf.SetMyConfig()
		backup.Init()
		sysmnt.Init()
		onstart.Init()
		controller.Init()

		onstart.Run()
	},
}
var path string

func init() {
	runCmd.Flags().StringVarP(&path, "path", "p", ".", "配置文件路径")
	_ = runCmd.MarkFlagRequired("path")
	rootCmd.AddCommand(runCmd)
}
