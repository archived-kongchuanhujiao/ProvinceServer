package client

import (
	"github.com/kongchuanhujiao/server/internal/app/client/internal"
	"github.com/kongchuanhujiao/server/internal/app/client/message"
	"github.com/kongchuanhujiao/server/internal/pkg/configs"
)

var cli *internal.QQ // client 客户端

// NewClient 新建客户端
func NewClient() {
	conf := configs.GetConfigs()
	cli = internal.NewClient(conf.Number, conf.Password)
}

// GetClient 获取客户端。
// 执行函数：NewClient 前调用返回值为 nil
func GetClient() *internal.QQ { return cli }

// SetCallback 设置回调
func SetCallback(f message.Callback) { cli.SetCallback(f) }
