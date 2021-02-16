package configs

import (
	"coding.net/kongchuanhujiao/server/internal/pkg/logger"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

// Configs 配置
type Configs struct {
	QQNumber   uint64 `yaml:"qqNumber"`    // QQ 号
	QQPassword string `yaml:"qqPassword"`  // QQ 密码
	DBURL      string `yaml:"databaseUrl"` // 数据库 URL
}

// GetConfigs 获取配置
func GetConfigs() (c *Configs) {

	bytes, err := ioutil.ReadFile("./config.yml")

	if err != nil {
		logger.Error("读取配置文件失败", zap.Error(err))
	}

	err = yaml.Unmarshal(bytes, &c)
	if err != nil {
		logger.Panic("无法读取配置信息", zap.Error(err))
	}

	logger.Debug("读取配置信息成功", zap.Any("配置", c))
	return
}
