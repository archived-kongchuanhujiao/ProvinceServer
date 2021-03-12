package config

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/base64"
	"io/ioutil"

	"github.com/kongchuanhujiao/server/internal/pkg/logger"

	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

var conf *Configs

type (
	// Configs 配置
	Configs struct {
		Number   uint64
		Password string
		Database string
		JWT      JWT
	}

	// JWT 配置
	JWT struct {
		Iss string
		Key *ecdsa.PrivateKey
	}

	// config 配置
	configs struct {
		Number   uint64 `yaml:"number"`
		Password string `yaml:"password"`
		Database string `yaml:"database"`
		JWT      struct {
			Iss string `yaml:"iss"`
			Key string `yaml:"key"`
		} `yaml:"jwt"`
	}
)

// ReadConfigs 读取配置
func ReadConfigs() {

	var c configs

	bytes, err := ioutil.ReadFile(".kongchuanhujiao/config.yml")
	if err != nil {
		logger.Panic("读取配置文件失败", zap.Error(err))
	}

	err = yaml.Unmarshal(bytes, &c)
	if err != nil {
		logger.Panic("无法读取配置信息", zap.Error(err))
	}

	logger.Debug("读取配置信息成功")

	conf = &Configs{
		Number: c.Number, Password: c.Password, Database: c.Database,
		JWT: JWT{Iss: c.JWT.Iss},
	}

	err = parseJWTKey(c.JWT.Key)
	if err != nil {
		logger.Panic("解析 JWT 私钥失败", zap.Error(err))
	}

	return
}

// parseJWTKey 解析 JWT Key
func parseJWTKey(k string) (err error) {

	bytes, err := base64.StdEncoding.DecodeString(k)
	if err != nil {
		return
	}

	conf.JWT.Key, err = x509.ParseECPrivateKey(bytes)

	return
}
