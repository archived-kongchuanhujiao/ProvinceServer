package main

import (
	"github.com/kongchuanhujiao/server/internal/app/apis"
	"github.com/kongchuanhujiao/server/internal/app/client"
	"github.com/kongchuanhujiao/server/internal/app/client/message"
	"github.com/kongchuanhujiao/server/internal/app/datahub/pkg"
	"github.com/kongchuanhujiao/server/internal/app/kongchuanhujiao/wenda"
	"github.com/kongchuanhujiao/server/internal/pkg/configs"
	"github.com/kongchuanhujiao/server/internal/pkg/logger"
)

// main 启动函数
func main() {

	logger.Named("主").Info("Copyright (C) 2020-present | version：21.03.06+" + configs.Commit)

	configs.ReadConfigs()

	pkg.ConnectDatabase()
	client.NewClient()
	client.SetCallback(func(m *message.Message) {
		wenda.HandleTest(m)
		wenda.HandleAnswer(m)
	})
	apis.StartApis()

	select {}
}
