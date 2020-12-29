package clients

// Clients 客户端
type Clients interface {
	ReceiveMessage(message *Message) // ReceiveMessage 接收消息
	SendMessage(message *Message)                    // SendMessage 发送消息
}
