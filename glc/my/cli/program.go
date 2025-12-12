package cli

import (
	"fmt"
	"os"

	"github.com/kardianos/service"
	"glc/my/config"
	"glc/my/global"
)

// 创建系统服务
type program struct{}

func (p *program) Start(s service.Service) error {
	// Start should not block. Do the actual work async.
	go p.run()
	return nil
}
func (p *program) run() {
	// Do work here
	Execute()
	os.Exit(1)
}
func (p *program) Stop(s service.Service) error {
	// Stop should not block. Return with a few seconds.
	return nil
}

// NewService 创建服务配置
func NewService() (sc service.Service, err error) {
	// 初始化全局变量
	filePath, err := os.Getwd() // 取程序所在路径
	if err != nil {
		return
	}
	global.Path = filePath
	global.Conf = config.New()
	err = global.Conf.ReadConfig(global.Path, global.ConfFolder, global.AppConfFileName)
	if err != nil {
		return
	}

	// 创建服务配置
	name := global.AppInfo.ServiceName
	displayName := global.AppInfo.ServiceDisplayName
	desc := fmt.Sprintf("%v Port:%v", global.AppInfo.ServiceDesc, global.Conf.App.Port)
	if global.Conf.Service.SuffixName != "" {
		name = fmt.Sprintf("%v$%v", name, global.Conf.Service.SuffixName)
		displayName = fmt.Sprintf("%v（%v）", displayName, global.Conf.Service.SuffixName)
	}

	config := &service.Config{
		Name:         name,
		DisplayName:  displayName,
		Description:  desc,
		Arguments:    []string{"run", fmt.Sprintf("--path=%v", filePath)},
		Dependencies: global.Conf.Service.Dependencies, // 依赖项
	}
	return service.New(&program{}, config)
}
