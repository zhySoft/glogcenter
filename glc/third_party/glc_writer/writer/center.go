package writer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

var (
	consoleBufPool = sync.Pool{
		New: func() interface{} {
			return bytes.NewBuffer(make([]byte, 0, 100))
		},
	}
)

const (
	TraceidFieldName    = "traceid"    // 追踪id 只有追踪id需要通过ctx传递，用hook从ctx取值
	KeywordFieldName    = "keyword"    // 关键字
	SystemFieldName     = "system"     // 系统名
	ServerNameFieldName = "servername" // 服务器名
	ServerIpFieldName   = "serverip"   // 服务器ip
	ClientIpFieldName   = "clientip"   // 客户端ip
	UserFieldName       = "user"       // 用户
)

// Formatter 将输入转换为格式化字符串。
type Formatter func(interface{}) string

var Debug bool
var LogAge time.Duration // 日志有效期

type LogCenterWriter struct {
	Out               io.Writer
	syncOut, asyncOut io.Writer // 同步记录日志 异步记录日志

	// TimeFormat 指定输出时间戳的格式。
	TimeFormat string

	// TimeLocation 告诉 LogCenterWriter 的默认格式时间戳如何本地化时间。
	TimeLocation *time.Location
}

// NewLogCenterWriter 创建并初始化一个新的 LogCenterWriter。
// 这是一个异步的日志，有一个 goroutine 单独用来写入日志。
//
// wg 是一个 sync.WaitGroup 用来等待日志写入完成
//
// url LogCenterWriter 的地址
//
// logDir 当写入日志中心失败时，写入本地日志目录
//
// storageChanLength 日志通道长度，根据日志写入频率适当修改，当为0时用于同步接收日志，大于0时异步接收日志
func NewLogCenterWriter(wg *sync.WaitGroup, url, logDir string, storageChanLength uint) LogCenterWriter {
	// 产生Fatal，Panic日志时，goroutine被强制终止，日志写入失败。
	// 同步写入，用于接收Fatal，Panic日志
	syncWriter := newCenterWriter(wg, url, logDir, 0)
	// 异步写入，用于接收其他日志
	asyncWriter := newCenterWriter(wg, url, logDir, storageChanLength)

	w := LogCenterWriter{
		syncOut:      syncWriter,
		asyncOut:     asyncWriter,
		TimeFormat:   time.RFC3339Nano,
		TimeLocation: time.Local,
	}

	return w
}

// 写写用格式化器转换JSON输入，并附加到 w.Out。
func (w LogCenterWriter) Write(p []byte) (n int, err error) {
	var buf = consoleBufPool.Get().(*bytes.Buffer)
	defer func() {
		buf.Reset()
		consoleBufPool.Put(buf)
	}()
	// byte to map
	var evt map[string]interface{}
	d := json.NewDecoder(bytes.NewReader(p))
	d.UseNumber()
	err = d.Decode(&evt)
	if err != nil {
		return n, fmt.Errorf("cannot decode event: %s", err)
	}
	// 如果是Fatal，Panic日志，重定向输出
	if evt[zerolog.LevelFieldName] == zerolog.LevelFatalValue || evt[zerolog.LevelFieldName] == zerolog.LevelPanicValue {
		w.Out = w.syncOut
	} else {
		w.Out = w.asyncOut
	}
	// 格式化字段值
	w.formatFieldsValue(evt)
	// 拼接结构化日志中不属于log要的字段到text
	w.joinTextLog(evt)
	// 修改字段名
	w.modifyFieldsName(evt)
	// map to byte
	o, err := json.Marshal(evt)
	if err != nil {
		return
	}
	_, err = bytes.NewBuffer(o).WriteTo(w.Out)
	return len(p), err
}

// Call the underlying writer's Close method if it is an io.Closer.
// Otherwise does nothing.
func (w LogCenterWriter) Close() error {
	if closer, ok := w.Out.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}

// joinTextLog 拼接text
func (w LogCenterWriter) joinTextLog(evt map[string]interface{}) {
	var text string
	if evt[zerolog.MessageFieldName] != nil {
		text = evt[zerolog.MessageFieldName].(string)
	}

	for k, v := range evt {
		switch k {
		case zerolog.MessageFieldName, zerolog.TimestampFieldName, zerolog.LevelFieldName, TraceidFieldName, KeywordFieldName, SystemFieldName, ServerNameFieldName, ServerIpFieldName, ClientIpFieldName, UserFieldName:
			continue
		default:
			delete(evt, k)
			text = fmt.Sprintf("%v\n%v=%v", text, k, v)
		}
	}
	// glc日志中心不记录空日志
	if text == "" {
		text = "nil"
	}
	evt[zerolog.MessageFieldName] = text
}

// formatFieldsValue 格式化字段值
func (w LogCenterWriter) formatFieldsValue(evt map[string]interface{}) {
	var f Formatter
	for k, v := range evt {
		switch k {
		case zerolog.TimestampFieldName:
			time.Parse(zerolog.TimeFieldFormat, v.(string))
			f = formatTimestamp(w.TimeFormat, w.TimeLocation)
		case zerolog.CallerFieldName:
			f = formatCaller
		default:
			f = nil
		}

		if f != nil {
			evt[k] = f(v)
		} else {
			evt[k] = fmt.Sprintf("%s", v)
		}
	}
	return
}

// modifyFieldsName 修改字段名称
func (w LogCenterWriter) modifyFieldsName(evt map[string]interface{}) {
	for k, v := range evt {
		switch k {
		case zerolog.TimestampFieldName:
			delete(evt, zerolog.TimestampFieldName)
			evt["date"] = v
		case zerolog.LevelFieldName:
			delete(evt, zerolog.LevelFieldName)
			evt["loglevel"] = v
		case zerolog.MessageFieldName:
			delete(evt, zerolog.MessageFieldName)
			evt["text"] = v
		case KeywordFieldName:
			delete(evt, KeywordFieldName)
			evt["keyword"] = v
		case TraceidFieldName:
			delete(evt, TraceidFieldName)
			evt["traceid"] = v
		case SystemFieldName:
			delete(evt, SystemFieldName)
			evt["system"] = v
		}
	}
	return
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

// formatTimestamp 格式化时间戳
func formatTimestamp(timeFormat string, location *time.Location) Formatter {
	if timeFormat == "" {
		timeFormat = time.RFC3339Nano
	}
	if location == nil {
		location = time.Local
	}

	return func(i interface{}) string {
		t := "<nil>"
		switch tt := i.(type) {
		case string:
			ts, err := time.ParseInLocation(zerolog.TimeFieldFormat, tt, location)
			if err != nil {
				t = tt
			} else {
				t = ts.In(location).Format(timeFormat)
			}
		case json.Number:
			i, err := tt.Int64()
			if err != nil {
				t = tt.String()
			} else {
				var sec, nsec int64

				switch zerolog.TimeFieldFormat {
				case zerolog.TimeFormatUnixNano:
					sec, nsec = 0, i
				case zerolog.TimeFormatUnixMicro:
					sec, nsec = 0, int64(time.Duration(i)*time.Microsecond)
				case zerolog.TimeFormatUnixMs:
					sec, nsec = 0, int64(time.Duration(i)*time.Millisecond)
				default:
					sec, nsec = i, 0
				}

				ts := time.Unix(sec, nsec)
				t = ts.In(location).Format(timeFormat)
			}
		}
		return t
	}
}
