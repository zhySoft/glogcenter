package writer_test

import (
	"context"
	"fmt"
	center "github.com/zhySoft/glogcenter/glc/third_party/glc_writer/writer"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

func ExampleNewLogCenterWriter() {
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

	logger := zerolog.New(output).With().Caller().Timestamp().
		Str(center.SystemFieldName, "Pivas").
		Str(center.ClientIpFieldName, "127.0.0.1").
		Str(center.UserFieldName, "admin").
		Logger()
	ctx := context.Background()
	ctx = context.WithValue(ctx, center.TraceidFieldName, "0123210")

	logger = logger.Hook(center.CtxHook{}).With().Ctx(ctx).Logger()
	logger.Trace().Msg("trace")
	logger.Debug().Msg("debug")
	logger.Info().Msg("info")
	logger.Warn().Msg("warn")
	logger.Error().Msg("error")
	//logger.Fatal().Msg("fatal")
	//logger.Panic().Msg("panic")
	fmt.Println("log start")
	for i := 0; i < 100; i++ {
		logger.Info().Str(center.KeywordFieldName, "test").Msg(fmt.Sprintf("message %v", i))
	}
	fmt.Println("log end")
	wg.Wait()
}
