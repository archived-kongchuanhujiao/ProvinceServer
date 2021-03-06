package pkg

import (
	"github.com/kongchuanhujiao/server/internal/app/datahub/internal/maria"
	"github.com/kongchuanhujiao/server/internal/pkg/configs"
)

// ConnectDatabase 连接所有数据库
func ConnectDatabase() {
	maria.Connect(configs.GetDatabaseConf())
}
