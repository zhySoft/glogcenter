package conf

import (
	"glc/my/global"
	"path/filepath"
	"strconv"
)

func SetMyConfig() {
	storeRoot = filepath.Join(global.Path, global.LogFolder)
	serverPort = strconv.Itoa(global.Conf.App.Port)
	saveDays = global.Conf.App.SaveDays
	mulitLineSearch = global.Conf.App.SearchMulitLint
	enableSecurityKey = true
	headerSecurityKey = "Token"
	securityKey = global.Token
	enableLogin = true
	username = "zhy"
	password = "jack110"
}
