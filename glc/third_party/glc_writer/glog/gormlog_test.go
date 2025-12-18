package glog_test

import (
	"context"
	"fmt"
	"github.com/zhySoft/glogcenter/glc/third_party/glc_writer/glog"
	center "github.com/zhySoft/glogcenter/glc/third_party/glc_writer/writer"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

func ExampleNew() {
	var wg sync.WaitGroup
	zerolog.TimeFieldFormat = time.RFC3339Nano
	zerolog.LevelFieldName = "severity"
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	stdout := zerolog.NewConsoleWriter()
	stdout.TimeFormat = "15:04:05"
	url := "http://127.0.0.1:10021/glc/v1/log/add"
	logDir, _ := os.Getwd()
	logDir = filepath.Join(logDir, "log")

	center.Debug = true
	httpout := center.NewLogCenterWriter(&wg, url, logDir, 10)
	output := zerolog.MultiLevelWriter(stdout, httpout)
	logger := zerolog.New(output).
		With().Caller().Timestamp().
		Str(center.SystemFieldName, "Pivas").
		Str(center.ClientIpFieldName, "127.0.0.1").
		Str(center.UserFieldName, "admin").
		Logger()

	ctx := context.Background()
	ctx = context.WithValue(ctx, center.TraceidFieldName, "0123210")

	logger = logger.Hook(center.CtxHook{}).With().Ctx(ctx).Logger()

	conf := new(gorm.Config)
	conf.Logger = glog.New(&logger)
	connStr := fmt.Sprintf("server=%s%s;user id=%s;password=%s;database=%s;",
		"127.0.0.1", "\\MSSQLSERVER", "sa", "sa", "master")
	db, err := gorm.Open(sqlserver.Open(connStr), conf)
	if err != nil {
		return
	}
	err = db.WithContext(ctx).Exec("SELECT * FROM sys.configurations").Error
	wg.Wait()
}
