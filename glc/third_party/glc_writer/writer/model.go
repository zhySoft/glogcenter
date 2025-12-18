package writer

// responseModel 日志中心返回的响应模型
type responseModel struct {
	Code    int    `json:"code" xml:"code"`
	Success bool   `json:"success" xml:"success"`
	Message string `json:"message" xml:"message"`
}
