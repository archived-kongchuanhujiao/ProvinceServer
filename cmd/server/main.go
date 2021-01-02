package main

import (
	"coding.net/kongchuanhujiao/server/internal/app/apis"
	"coding.net/kongchuanhujiao/server/internal/app/client"
	"coding.net/kongchuanhujiao/server/internal/app/client/clientmsg"
	"coding.net/kongchuanhujiao/server/internal/app/datahub/datahubpkg"
	"coding.net/kongchuanhujiao/server/internal/pkg/logger"
)

// main 启动函数
func main() {

	logger.Info("Copyright (C) 2020-present | version：")

	datahubpkg.ConnectAllDatabase()
	apis.NewApis()
	client.NewClients()
	client.SetCallback(func(m *clientmsg.Message) {

		c, ok := m.Chain[0].(*clientmsg.Text)
		if !ok {
			return
		}

		if c.Content == "你好空传互教" {
			client.GetClient().SendMessage(m.AddText("\nReply: 你好。"))
		}

		if c.Content == "外部测试" {
			client.GetClient().SendMessage(clientmsg.NewTextMessage("你好"))
		}

	})

	select {}

}
