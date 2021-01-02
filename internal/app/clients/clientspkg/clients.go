package clientspkg

import (
	"coding.net/kongchuanhujiao/server/internal/app/clients/clientspublic"
	"coding.net/kongchuanhujiao/server/internal/app/clients/internal/qq"
	"coding.net/kongchuanhujiao/server/internal/pkg/configs"
)

var (
	qqClient *qq.QQ // qqClient QQ 客户端
)

// NewClients 新建客户端
func NewClients() {
	conf := configs.GetConfigs()

	qqClient = qq.NewQQClient(conf.QQNumber, conf.QQPassword)
}

// GetClient 获取客户端。
// 执行函数：NewClients 前调用返回值为 nil
func GetClient(t string) (c clientspublic.Client) {
	c = qqClient
	return
}

// SetCallback 设置回调
func SetCallback(f clientspublic.Callback) {
	cs := []clientspublic.Client{qqClient}
	for _, v := range cs {
		v.SetCallback(f)
	}
}
