package datahubpkg

import "coding.net/kongchuanhujiao/server/internal/app/datahub/internal/maria"

// ConnectAllDatabase 连接所有数据库
func ConnectAllDatabase() {
	maria.Connect()
}
