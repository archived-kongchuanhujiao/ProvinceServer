package datahubpkg

import "coding.net/kongchuanhujiao/server/internal/app/datahub/internal/maria"

// ConnectDatabase 连接所有数据库
func ConnectDatabase() {
	maria.Connect()
}
