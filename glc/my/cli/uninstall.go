package cli

import (
	"log"

	"github.com/spf13/cobra"
)

// uninstallCmd 表示卸载命令
var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "卸载服务",
	Long:  "卸载服务",
	Run: func(cmd *cobra.Command, args []string) {
		sc, err := NewService()
		if err != nil {
			log.Fatalln(err)
			return
		}
		err = sc.Uninstall()
		if err != nil {
			log.Fatalln(err)
			return
		} else {
			log.Println("服务已卸载")
		}
	},
}

func init() {
	rootCmd.AddCommand(uninstallCmd)
}
