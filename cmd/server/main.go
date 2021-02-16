package main

import (
	"coding.net/kongchuanhujiao/server/internal/app/apis"
	"coding.net/kongchuanhujiao/server/internal/app/client"
	"coding.net/kongchuanhujiao/server/internal/app/client/clientmsg"
	"coding.net/kongchuanhujiao/server/internal/app/datahub/pkg"
	"coding.net/kongchuanhujiao/server/internal/app/kongchuanhujiao/wenda"
	"coding.net/kongchuanhujiao/server/internal/pkg/logger"
)

// main 启动函数
func main() {

	logger.Named("主").Info("Copyright (C) 2020-present | version：21.02.XX")

	pkg.ConnectDatabase()
	client.NewClient()
	client.SetCallback(func(m *clientmsg.Message) {
		wenda.HandleTest(m)
		wenda.HandleAnswer(m)
	})
	apis.NewApis()

	select {}
}
