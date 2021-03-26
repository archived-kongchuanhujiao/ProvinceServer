package kongchuanhujiao

// Response 响应
type Response struct {
	Status  uint16      `json:"status"`  // 状态码
	Message string      `json:"message"` // 消息
	Data    interface{} `json:"data"`    // 数据
}

var (
	DefaultSuccResp = &Response{Message: "ok"}               // DefaultSuccResp 默认成功响应
	DefaultErrResp  = &Response{Status: 1, Message: "服务器错误"} // DefaultErrResp 默认错误响应
)

// GenerateErrResp 生成成功响应 返回带数据的响应
func GenerateSuccResp(d interface{}) *Response { return &Response{Data: d, Message: "ok"} }

// GenerateErrResp 生成错误响应 返回带状态码和信息的响应
func GenerateErrResp(s uint16, msg string) *Response { return &Response{Status: s, Message: msg} }
