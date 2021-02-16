package pkg

import (
	"coding.net/kongchuanhujiao/server/internal/app/datahub/internal/maria"
	"coding.net/kongchuanhujiao/server/internal/pkg/configs"
)

// ConnectDatabase 连接所有数据库
func ConnectDatabase() {
	conf := configs.GetConfigs()
	maria.Connect(conf.Database)
}
