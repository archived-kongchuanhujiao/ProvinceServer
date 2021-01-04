package dingtalk

import (
	"coding.net/kongchuanhujiao/server/internal/pkg/logger"
	"github.com/CatchZeng/dingtalk"
	"go.uber.org/zap"
)

// TODO 这里引入第三方包
// TODO 这里不能调用其他的包，只能调用第三方包，否则软件架构会混乱。所有需要的数据在进来前就要准备好，不能从这里获取其他数据信息
// 另请查询 [apis.go#PostPUSHCENTER()]

var (
	client      *dingtalk.Client
	accessToken string
	secret      string
)

func init() {
	client = dingtalk.NewClient(accessToken, secret)
}

func push(content string, atMobiles []string, isAtAll bool) {
	msg := dingtalk.NewTextMessage().SetContent(content).SetAt(atMobiles, isAtAll)
	_, err := client.Send(msg)

	if err != nil {
		logger.Warn("发送钉钉消息失败", zap.Error(err))
	}
}
