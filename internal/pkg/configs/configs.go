package configs

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/base64"
	"io/ioutil"

	"github.com/kongchuanhujiao/server/internal/pkg/logger"

	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

// Configs 配置
type Configs struct {
	Number   uint64 `yaml:"number"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
	JWT      struct {
		Iss string            `yaml:"iss"`
		Key *ecdsa.PrivateKey `yaml:"key"`
	} `yaml:"jwt"`
}

// configs 配置
type configs struct {
	Number   uint64 `yaml:"number"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
	JWT      struct {
		Iss string `yaml:"iss"`
		Key string `yaml:"key"`
	} `yaml:"jwt"`
}

// GetConfigs 获取配置
func GetConfigs() (c *Configs) {

	var cc configs

	bytes, err := ioutil.ReadFile(".kongchuanhujiao/config.yml")
	if err != nil {
		logger.Error("读取配置文件失败", zap.Error(err))
	}

	err = yaml.Unmarshal(bytes, &cc)
	if err != nil {
		logger.Panic("无法读取配置信息", zap.Error(err))
	}

	k, err := base64.StdEncoding.DecodeString(cc.JWT.Key)
	if err != nil {
		return
	}

	kk, err := x509.ParseECPrivateKey(k)
	if err != nil {
		return
	}

	logger.Debug("读取配置信息成功")
	return &Configs{
		Number:   cc.Number,
		Password: cc.Password,
		Database: cc.Database,
		JWT: struct {
			Iss string            `yaml:"iss"`
			Key *ecdsa.PrivateKey `yaml:"key"`
		}{
			Iss: cc.JWT.Iss,
			Key: kk,
		},
	}
}
