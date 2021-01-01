package main

import (
	"coding.net/kongchuanhujiao/server/internal/app/apis"
	"coding.net/kongchuanhujiao/server/internal/app/clients/clientspkg"
	"coding.net/kongchuanhujiao/server/internal/app/clients/clientspublic"
	"coding.net/kongchuanhujiao/server/internal/app/datahub/datahubpkg"
	"coding.net/kongchuanhujiao/server/internal/pkg/logger"
)

// main 启动函数
func main() {

	logger.Info("Copyright (C) 2020-present | version：")

	datahubpkg.ConnectAllDatabase()
	apis.NewApis()
	clientspkg.NewClients()
	clientspkg.SetCallback(func(m *clientspublic.Message) {

		c, ok := m.Chain[0].(*clientspublic.Text)
		if !ok {
			return
		}

		if c.Content == "你好空传互教" {
			m.QuickMessage(m.AddText("\nReply: 你好。"))
		}
	})

	select {}

}
