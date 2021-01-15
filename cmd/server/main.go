package main

import (
	"coding.net/kongchuanhujiao/server/internal/app/apis"
	"coding.net/kongchuanhujiao/server/internal/app/client"
	"coding.net/kongchuanhujiao/server/internal/app/client/clientmsg"
	"coding.net/kongchuanhujiao/server/internal/app/datahub/datahubpkg"
	"coding.net/kongchuanhujiao/server/internal/app/kongchuanhujiao/wenda"
	"coding.net/kongchuanhujiao/server/internal/pkg/logger"
)

// main 启动函数
func main() {

	logger.Info("Copyright (C) 2020-present | version：")

	datahubpkg.ConnectAllDatabase()
	apis.NewApis()
	client.NewClients()
	client.SetCallback(func(m *clientmsg.Message) {
		if t := m.Chain[0]; t.TypeName() == "text" {
			c := t.(*clientmsg.Text)
			if c.Content == "你好" {
				client.GetClient().SendMessage(
					clientmsg.NewAtMessage(m.Target.ID).AddText("你好").SetGroupTarget(m.Target.Group),
				)
			}

		}
		wenda.HandleAnswer(m)
	})

	select {}

}
