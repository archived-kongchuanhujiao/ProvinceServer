package clients

import (
	"coding.net/kongchuanhujiao/server/internal/app/internal/clients"
	"coding.net/kongchuanhujiao/server/internal/app/internal/clients/dingtalk"
	"coding.net/kongchuanhujiao/server/internal/app/internal/clients/qq"
	"coding.net/kongchuanhujiao/server/internal/pkg/configs"
)

var (
	qqClient       *qq.QQ             // qqClient QQ 客户端
	dingTalkClient *dingtalk.DingTalk // dingTalkClient 钉钉 客户端
)

// NewClients 新建客户端
func NewClients() {

	conf := configs.GetConfigs()

	qqClient = qq.NewQQClient(conf.QQNumber, conf.QQPassword)
	dingTalkClient = dingtalk.NewDingTalkClient()

}

// GetClient 获取客户端
func GetClient(t string) (c clients.Clients) {
	c = qqClient
	if t == clients.DingTalkClient {
		c = dingTalkClient
	}
	return
}
