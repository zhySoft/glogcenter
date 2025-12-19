package glog

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/rs/zerolog"
	gormLog "gorm.io/gorm/logger"
)

type Glog struct {
	logger *zerolog.Logger
}

func New(logger *zerolog.Logger) *Glog {
	return &Glog{logger: logger}
}

func (g *Glog) LogMode(level gormLog.LogLevel) gormLog.Interface {
	var l zerolog.Level
	switch level {
	case gormLog.Silent:
		l = zerolog.Disabled
	case gormLog.Error:
		l = zerolog.ErrorLevel
	case gormLog.Warn:
		l = zerolog.WarnLevel
	case gormLog.Info:
		l = zerolog.InfoLevel
	}
	g.logger.Level(l)
	return g
}
func (g *Glog) Info(ctx context.Context, format string, args ...interface{}) {
	// 去除gorm日志中的caller敏感信息
	if format == "replacing callback `%s` from %s\n" && len(args) >= 2 {
		args[1] = formatCaller(args[1])
	}
	g.logger.Info().Ctx(ctx).Msgf(format, args...)
}
func (g *Glog) Warn(ctx context.Context, format string, args ...interface{}) {
	g.logger.Warn().Ctx(ctx).Msgf(format, args...)
}
func (g *Glog) Error(ctx context.Context, format string, args ...interface{}) {
	g.logger.Error().Ctx(ctx).Msgf(format, args...)
}
func (g *Glog) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if g.logger.GetLevel() > zerolog.PanicLevel {
		return
	}
	// 执行时间
	elapsed := time.Since(begin)
	elapsedStr := fmt.Sprintf("%.3fms", float64(elapsed.Nanoseconds())/1e6)

	sql, rows := fc()

	//ctx = context.WithValue(ctx, "elapsed", elapsedStr) // 执行时间
	//ctx = context.WithValue(ctx, "rows", rows)          // 影响行数

	g.logger.Debug().Ctx(ctx).Str("elapsed", elapsedStr).Int64("rows", rows).Msg(sql)
}

// formatCaller
func formatCaller(i interface{}) string {
	var c string
	if cc, ok := i.(string); ok {
		c = cc
	}
	if len(c) > 0 {
		if cwd, err := os.Getwd(); err == nil {
			if rel, err := filepath.Rel(cwd, c); err == nil {
				c = rel
			}
		}
	}
	return c
}
