package global

import (
	"glc/my/config"
	"log"
	"runtime"
	"strconv"

	"github.com/gotoeasy/glang/cmn"
)

// Init 初始化全局变量
func Init(path string) {
	// 初始化全局变量
	Path = path
	Conf = config.New()
	err := Conf.ReadConfig(Path, ConfFolder, AppConfFileName)
	if err != nil {
		log.Fatalf("读取配置文件失败:%v\n", err)
		return
	}
	// 设置glc配置
	enableConsoleLog := strconv.FormatBool(Conf.App.Debug)
	logLevel := "INFO"
	if Conf.App.Debug {
		logLevel = "DEBUG"
	}
	opt := cmn.GlcOptions{
		EnableConsoleLog: enableConsoleLog,
		LogLevel:         logLevel,
	}
	cmn.SetGlcClient(cmn.NewGlcClient(&opt))
	runtime.GOMAXPROCS(-1) // 使用的最大CPU数量，默认是最大CPU数量（设定值不在实际数量范围是按最大看待）

}
