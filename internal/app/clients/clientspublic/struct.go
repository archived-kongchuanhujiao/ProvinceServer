package clientspublic

const (
	QQClient       = "qq"       // QQClient QQ 客户端
	DingTalkClient = "dingtalk" // DingTalkClient 钉钉 客户端
)

type (
	Callback func(*Message) // Callback 回调

	// Message 消息
	Message struct {
		Client Client    // 客户端
		Chain  []Element // 消息链
		Target *Target   // 目标
	}

	// Target 目标
	Target struct {
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
