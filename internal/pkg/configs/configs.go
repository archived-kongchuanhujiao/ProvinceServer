package configs

import (
	"io/ioutil"

	"coding.net/kongchuanhujiao/server/internal/pkg/logger"

	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

// Configs 配置
type Configs struct {
	Number   uint64 `yaml:"number"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

// GetConfigs 获取配置
func GetConfigs() (c *Configs) {

	bytes, err := ioutil.ReadFile(".kongchuanhujiao/config.yml")
	if err != nil {
		logger.Error("读取配置文件失败", zap.Error(err))
	}

	err = yaml.Unmarshal(bytes, &c)
	if err != nil {
		logger.Panic("无法读取配置信息", zap.Error(err))
	}

	logger.Debug("读取配置信息成功")
	return
}
