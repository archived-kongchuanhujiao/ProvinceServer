package internal

import (
	"github.com/kongchuanhujiao/server/internal/app/client/message"

	"github.com/Mrs4s/MiraiGo/client"
	"go.uber.org/zap"
)

var loggerr = zap.L().Named("QQ") // loggerr 日志

// QQ 客户端
type QQ struct {
	client   *client.QQClient // 客户端
	callback message.Callback // 回调
}

// NewClient 新建 QQ 客户端
func NewClient(a uint64, p string) (q *QQ) {

	setProtocol()
	c := client.NewClient(int64(a), p)
	c.OnLog(setLogger)

	// 读取配置信息
	q = &QQ{client: c}

	if err := q.login(); err != nil {
		loggerr.Panic("登录失败", zap.Error(err))
	}
	q.setEventHandle()

	return
}

// SendMessage 发送消息
func (q *QQ) SendMessage(m *message.Message) {
	ms := q.transformToMiraiGO(m)
	if m.Target.Group != nil {
		q.client.SendGroupMessage(int64(m.Target.Group.ID), ms)
	} else {
		q.client.SendPrivateMessage(int64(m.Target.ID), ms)
	}
}

// ReceiveMessage 接收消息
func (q *QQ) ReceiveMessage(m *message.Message) {
	loggerr.Debug("接收消息", zap.Any("消息", m))
	if q.callback != nil && len(m.Chain) != 0 {
		q.callback(m)
	}
}

// SetCallback 设置回调
func (q *QQ) SetCallback(f message.Callback) {
	q.callback = f
}
