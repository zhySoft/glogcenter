package config

// App配置文件结构
type AppConf struct {
	//lock                   sync.RWMutex
	path, folder, fileName string     // 记录当前配置文件路径
	Service                serviceSet `toml:"Service" comment:"服务配置"`
	App                    appSet     `toml:"App" comment:"应用配置"`
}

// app配置
type appSet struct {
	Port            int  `toml:"Instance" comment:"端口"` // 端口
	Debug           bool `toml:"Debug" comment:"调试模式"`
	SaveDays        int  `toml:"SaveDays" comment:"日志保存天数"`
	SearchMulitLint bool `toml:"SearchMulitLine" comment:"是否对日志列的全部行进行索引检索"`
}

// service系统服务配置
type serviceSet struct {
	SuffixName   string   `toml:"SuffixName" comment:"服务后缀名"`  // 当安装多个服务时需要后缀名区分
	Dependencies []string `toml:"Dependencies" comment:"依赖项名"` // 依赖项名称
}
