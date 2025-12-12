package main

import (
	"fmt"
	"glc/my/cli"
	"glc/my/global"
)

/*
func main() {
	cmn.SetGlcClient(cmn.NewGlcClient(&cmn.GlcOptions{
		EnableConsoleLog: cmn.GetEnvStr("GLC_ENABLE_CONSOLE_LOG", "false"), // 关闭控制台日志输出
		LogLevel:         cmn.GetEnvStr("GLC_LOG_LEVEL", "INFO"),           // 控制台INFO日志级别输出
	}))

	runtime.GOMAXPROCS(conf.GetGoMaxProcess()) // 使用最大CPU数量
	onstart.Run()
}
*/

func main() {
	Run()
}
func Run() {
	global.AppInfo.AppName = "LogCenter"          // 项目名称
	global.AppInfo.Version = "v25.12.12.001"      // 项目版本
	global.AppInfo.ServiceName = "LogCenter"      // 服务名称
	global.AppInfo.ServiceDisplayName = "智慧园日志中心" // 服务显示名称
	global.AppInfo.ServiceDesc = "用于接收日志"         // 服务描述

	fmt.Println(global.AppInfo.ServiceDisplayName)
	fmt.Printf("version: %v\n", global.AppInfo.Version)
	fmt.Println("作者：王飞")

	//cli.Execute()
	sc, err := cli.NewService()
	if err != nil {
		panic(err)
	}
	err = sc.Run()
	if err != nil {
		panic(err)
	}
}
