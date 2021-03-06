package kongchuanhujiao

type Response struct { // Response 响应
	Status  uint16      `json:"status"`  // 状态码
	Message string      `json:"message"` // 消息
	Data    interface{} `json:"data"`    // 数据
}

var ErrRes = &Response{Status: 1, Message: "服务器错误"}
