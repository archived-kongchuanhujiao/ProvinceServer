package main

import (
	"github.com/kongchuanhujiao/server/internal/app/api"
	"github.com/kongchuanhujiao/server/internal/app/client"
	"github.com/kongchuanhujiao/server/internal/app/client/message"
	"github.com/kongchuanhujiao/server/internal/app/datahub/pkg"
	"github.com/kongchuanhujiao/server/internal/app/kongchuanhujiao/wenda"
	"github.com/kongchuanhujiao/server/internal/pkg/config"
	"github.com/kongchuanhujiao/server/internal/pkg/logger"
)

// main 启动函数
func main() {

	logger.Named("主").Info("Copyright (C) 2020-present | " + config.Commit)

	config.ReadConfigs()

	pkg.ConnectDatabase()
	client.NewClient()
	client.SetCallback(func(m *message.Message) {
		wenda.HandleTest(m)
		wenda.HandleAnswer(m)
		wenda.HandleWordStat(m)
	})
	api.StartApis()

	select {}
}
