package config

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	tomlv2 "github.com/pelletier/go-toml/v2"
)

func New() *AppConf {
	return &AppConf{}
}

// initExample 初始化配置文件示例
func (c *AppConf) initExample(path, folder, fileName string) (err error) {
	example := AppConf{
		Service: serviceSet{
			SuffixName:   "example",
			Dependencies: []string{""},
		},
		App: appSet{
			Port:            8080,
			Debug:           false,
			SaveDays:        7,
			SearchMulitLint: false,
		},
	}
	buf := bytes.Buffer{}
	enc := tomlv2.NewEncoder(&buf)
	enc.SetIndentTables(true) // 缩进
	err = enc.Encode(example)
	if err != nil {
		return
	}
	//  创建目录
	confPath := filepath.Join(path, folder)
	err = os.MkdirAll(confPath, os.ModePerm)
	if err != nil {
		return
	}
	// 写入示例
	filePath := filepath.Join(path, folder, fileName)
	filePath = fmt.Sprintf("%s.example", filePath)

	return os.WriteFile(filePath, buf.Bytes(), os.ModePerm)
}

// readConfig 读取配置文件
func (c *AppConf) readConfig(path, folder, fileName string) (err error) {
	filePath := filepath.Join(path, folder, fileName)
	data, err := os.ReadFile(filePath)
	if err != nil {
		return
	}
	c.path = path
	c.folder = folder
	c.fileName = fileName
	return tomlv2.Unmarshal(data, c)
}

// ReadConfig 读取配置文件
func (c *AppConf) ReadConfig(path, folder, fileName string) (err error) {
	err = c.readConfig(path, folder, fileName)
	if err != nil {
		return c.initExample(path, folder, fileName)
	}
	return
}
