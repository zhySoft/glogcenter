package global

import (
	"glc/my/config"
)

// 全局变量
var (
	AppInfo appInfo         // 项目信息
	Path    string          // 程序运行时的路径
	Conf    *config.AppConf // 配置
)

type appInfo struct {
	AppName            string // 项目名称
	Version            string // 项目版本
	ServiceName        string // 服务名称
	ServiceDisplayName string // 服务显示名称
	ServiceDesc        string // 服务描述
}
