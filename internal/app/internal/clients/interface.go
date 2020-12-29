package clients

type (
	// Clients 客户端
	Clients interface {
		ReceiveMessage(message *Message) // ReceiveMessage 接收消息
		SendMessage(message *Message)    // SendMessage 发送消息
	}

	// Element 消息链
	Element interface {
		TypeName() string // TypeName 类型名
	}
)
