package internal

import (
	"coding.net/kongchuanhujiao/server/internal/app/client/clientmsg"
	"coding.net/kongchuanhujiao/server/internal/app/kongchuanhujiao/public/wendapkg"
	"coding.net/kongchuanhujiao/server/internal/pkg/logger"

	"github.com/Mrs4s/MiraiGo/client"
	"go.uber.org/zap"
)

var loggerr = logger.Named("QQ客户端") // loggerr 日志

// QQ QQ 客户端
type QQ struct {
	client   *client.QQClient   // 客户端
	callback clientmsg.Callback // 回调
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
func (q *QQ) SendMessage(m *clientmsg.Message) {
	ms := q.transformToMiraiGO(m)
	if m.Target.Group != nil {
		q.client.SendGroupMessage(int64(m.Target.Group.ID), ms)
	} else {
		q.client.SendPrivateMessage(int64(m.Target.ID), ms)
	}
	loggerr.Info("发送消息", zap.Any("消息", m.Chain))
}

// ReceiveMessage 接收消息
func (q *QQ) ReceiveMessage(m *clientmsg.Message) {
	logger.Debug("接收消息", zap.Any("消息", m))
	if q.callback != nil && len(m.Chain) != 0 {
		q.callback(m)
	}
}

// SetCallback 设置回调
func (q *QQ) SetCallback(f clientmsg.Callback) {
	q.callback = f
}

// GetGroups 获取群
func (q *QQ) GetGroups() *wendapkg.Groups {
	g := wendapkg.Groups{}
	for _, v := range q.client.GroupList {
		g[uint64(v.Code)] = v.Name
	}
	return &g
}

// GetGroupMembers 获取群成员
func (q *QQ) GetGroupMembers(i uint64) *wendapkg.GroupMembers {
	m := wendapkg.GroupMembers{}
	for _, v := range q.client.FindGroup(int64(i)).Members {
		m[uint64(v.Uin)] = v.DisplayName()
	}
	return &m
}
