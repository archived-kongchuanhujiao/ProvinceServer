package dingtalk

import (
	"coding.net/kongchuanhujiao/server/internal/app/clients/clientspublic"
	"coding.net/kongchuanhujiao/server/internal/pkg/logger"

	"go.uber.org/zap"
)

var loggerr = logger.Named("钉钉客户端") // loggerr 日志

// DingTalk QQ 客户端
type DingTalk struct {
	callback clientspublic.Callback // 回调
}

// NewDingTalkClient 新建 钉钉 客户端
func NewDingTalkClient() (d *DingTalk) {
	d = &DingTalk{}
	return
}

func (d *DingTalk) SendMessage(m *clientspublic.Message) {
	loggerr.Info("发送群消息", zap.Any("消息", m.Chain))
}

func (d *DingTalk) ReceiveMessage(m *clientspublic.Message) {
	logger.Debug("接收消息", zap.Any("消息", m))
	if d.callback != nil {
		d.callback(m)
	}
}

func (d *DingTalk) SetCallback(f clientspublic.Callback) {
	d.callback = f
}
