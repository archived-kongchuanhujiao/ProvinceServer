package main

import (
	_ "github.com/kongchuanhujiao/server/cmd/server/logger"
	"github.com/kongchuanhujiao/server/internal/app/api"
	"github.com/kongchuanhujiao/server/internal/app/client"
	"github.com/kongchuanhujiao/server/internal/app/client/message"
	"github.com/kongchuanhujiao/server/internal/app/datahub/pkg"
	"github.com/kongchuanhujiao/server/internal/app/kongchuanhujiao/ciyun"
	"github.com/kongchuanhujiao/server/internal/app/kongchuanhujiao/wenda"
	"github.com/kongchuanhujiao/server/internal/pkg/config"

	"go.uber.org/zap"
)

// main 启动函数
func main() {

	zap.L().Named("主").Info("Copyright (C) 2020-present | " + config.Commit)

	config.ReadConfigs()

	pkg.ConnectDatabase()
	client.NewClient()
	client.SetCallback(func(m *message.Message) {
		wenda.HandleTest(m)
		wenda.HandleAnswer(m)
		ciyun.HandleWordStat(m)
	})
	api.StartApis()

	select {}
}
