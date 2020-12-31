package clientspublic

type (
	// Client 客户端
	Client interface {
		ReceiveMessage(*Message) // ReceiveMessage 接收消息
		SendMessage(*Message)    // SendMessage 发送消息
		SetCallback(Callback)    // SetCallback 设置回调
	}

	// Element 消息链
	Element interface {
		TypeName() string // TypeName 类型名
	}
)
