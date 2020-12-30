package clients

const (
	QQClient       = "qq"       // QQClient QQ 客户端
	DingTalkClient = "dingtalk" // DingTalkClient 钉钉 客户端
)

type (
	// Message 消息
	Message struct {
		Client string    // 客户端
		Chain  []Element // 消息链
		Sender *Sender   // 发送者
	}

	// Sender 发送者
	Sender struct {
		ID    uint64 // 唯一识别码
		Name  string // 名称
		Group *Group // 群
	}

	// Group 群
	Group struct {
		ID   uint64 // 唯一识别码
		Name string // 名称
	}
)
