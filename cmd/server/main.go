package main

import (
	"coding.net/kongchuanhujiao/server/internal/app/clients/clientspkg"
	"coding.net/kongchuanhujiao/server/internal/app/clients/clientspublic"
	"coding.net/kongchuanhujiao/server/internal/pkg/logger"
)

// main 启动函数
func main() {

	logger.Info("Copyright (C) 2020-present | version：")

	// 启动 API

	clientspkg.NewClients()
	clientspkg.SetCallback(func(message *clientspublic.Message) {
		return
	})

	select {}

}
