package qq

import (
	"coding.net/kongchuanhujiao/server/internal/app/internal/clients"
	"coding.net/kongchuanhujiao/server/internal/pkg/logger"

	"github.com/Mrs4s/MiraiGo/client"
	"go.uber.org/zap"
)

var loggerr = logger.Named("QQ客户端") // loggerr 日志

// QQ QQ 客户端
type QQ struct {
	client *client.QQClient // 客户端
}

// NewQQClient 新建 QQ 客户端
func NewQQClient(a uint64, p string) (q *QQ) {

	setProtocol()
	c := client.NewClient(int64(a), p)
	c.OnLog(setLogger)

	// 读取配置信息
	q = &QQ{c}

	if err := q.login(); err != nil {
		loggerr.Panic("登录失败", zap.Error(err))
	}
	return
}

func (q *QQ) SendMessage(message *clients.Message) {
	ms := q.transformToMiraiGO(message)
	if message.Target.Group == nil {
		q.client.SendGroupMessage(int64(message.Target.Group.ID), ms)
	} else {
		q.client.SendPrivateMessage(int64(message.Target.ID), ms)
	}
	loggerr.Info("发送群消息", zap.Any("消息", message.Chain))
}

func (q *QQ) ReceiveMessage(message *clients.Message) {

}
