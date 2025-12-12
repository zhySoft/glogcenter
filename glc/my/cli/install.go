package cli

import (
	"log"

	"github.com/spf13/cobra"
)

// installCmd 表示 install 命令
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "安装为系统服务",
	Long:  "安装为系统服务",
	Run: func(cmd *cobra.Command, args []string) {
		sc, err := NewService()
		if err != nil {
			log.Fatalln(err)
			return
		}
		err = sc.Install()
		if err != nil {
			log.Fatalln(err)
			return
		} else {
			log.Println("安装成功")
		}
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}
