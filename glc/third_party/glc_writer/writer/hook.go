package writer

import (
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

// CtxHook 用于处理日志中心 ctx 中传递的字段。
//
// 当前只处理 traceid 字段
type CtxHook struct {
}

func (h CtxHook) Run(e *zerolog.Event, level zerolog.Level, message string) {
	ctx := e.GetCtx()
	ctxType := reflect.TypeOf(ctx).String()
	switch ctxType {
	case "gin.Context":
		v := reflect.ValueOf(ctx)
		t := reflect.TypeOf(gin.Context{})
		gtx := v.Convert(t).Interface().(gin.Context)
		traceid := gtx.GetString(TraceidFieldName)
		e.Str(TraceidFieldName, traceid)
	case "*gin.Context":
		v := reflect.ValueOf(ctx)
		t := reflect.TypeOf(&gin.Context{})
		gtx := v.Convert(t).Interface().(*gin.Context)
		traceid := gtx.GetString(TraceidFieldName)
		e.Str(TraceidFieldName, traceid)
	default:
		traceid := ctx.Value(TraceidFieldName)
		if traceid == nil {
			return
		}
		if reflect.TypeOf(traceid).Kind() == reflect.String {
			e.Str(TraceidFieldName, traceid.(string))
		}
	}
}
