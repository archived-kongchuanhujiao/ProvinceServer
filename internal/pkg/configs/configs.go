package configs

import (
	"os"
	"strconv"

	"coding.net/kongchuanhujiao/server/internal/pkg/logger"
)

// Configs 配置
type Configs struct {
	QQNumber   uint64 // QQ 号
	QQPassword string // QQ 密码
	DBURL      string // 数据库 URL
}

// GetConfigs 获取配置
func GetConfigs() *Configs {
	n, err := strconv.ParseUint(os.Getenv("KQQNum"), 10, 64)
	if err != nil {
		logger.Panic("无法读取配置信息")
	}

	logger.Debug("读取配置信息成功")
	return &Configs{n, os.Getenv("KQQPWA"), os.Getenv("KDBURL")}
}
