/**
 * 日志模型
 * 1）面向日志接口，设定常用属性方便扩充
 */
package logdata

import (
	"encoding/json"

	"github.com/gotoeasy/glang/cmn"
)

// Text是必须有的日志内容，Id自增，内置其他属性可选
// 其中Tags是空格分隔的标签，日期外各属性值会按空格分词
// 对应的json属性统一全小写
type LogDataModel struct {
	Id         string `json:"id,omitempty" xml:"id"`                 // 从1开始递增
	Text       string `json:"text,omitempty" xml:"text"`             // 【必须】日志内容，多行时仅为首行，直接显示用，是全文检索对象
	Date       string `json:"date,omitempty" xml:"date"`             // 日期（格式YYYY-MM-DD HH:MM:SS.SSS）
	System     string `json:"system,omitempty" xml:"system"`         // 系统名
	ServerName string `json:"servername,omitempty" xml:"servername"` // 服务器名
	ServerIp   string `json:"serverip,omitempty" xml:"serverip"`     // 服务器IP
	ClientIp   string `json:"clientip,omitempty" xml:"clientip"`     // 客户端IP
	TraceId    string `json:"traceid,omitempty" xml:"traceid"`       // 跟踪ID
	LogLevel   string `json:"loglevel,omitempty" xml:"loglevel"`     // 日志级别（debug、info、warn、error）
	User       string `json:"user,omitempty" xml:"user"`             // 用户
	Detail     string `json:"detail,omitempty" xml:"detail"`         // 【内部字段】多行时的详细日志信息，通常是包含错误堆栈等的日志内容
	StoreName  string `json:"storename,omitempty" xml:"storename"`   // 日志仓名称（未存储，仅赋值给前端使用）
	Keyword    string `json:"keyword,omitempty" xml:"keyword"`       // 关键字
}

func (d *LogDataModel) ToJson() string {
	bt, _ := json.Marshal(d)
	return cmn.BytesToString(bt)
}

func (d *LogDataModel) LoadJson(jsonstr string) error {
	return json.Unmarshal(cmn.StringToBytes(jsonstr), d)
}
