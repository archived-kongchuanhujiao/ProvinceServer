package logger

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var logger *zap.Logger

// fileCores 文件核心
type fileCores []struct {
	v    *zapcore.Core        // 变量
	name string               // 文件名
	age  int                  // 最大保存天数
	l    zapcore.LevelEnabler // 等级
}

// config 配置
func config() (console zapcore.EncoderConfig, file zapcore.EncoderConfig) {
	console = zap.NewProductionEncoderConfig()
	console.EncodeTime = zapcore.RFC3339TimeEncoder

	file = console

	console.EncodeLevel = zapcore.CapitalColorLevelEncoder
	file.EncodeCaller = zapcore.FullCallerEncoder
	return
}

// fileName 时间
func fileName(t time.Time) string { return t.Format("logs/2006年01月02日/") }

// level 等级
func level() (debug zap.LevelEnablerFunc, info zap.LevelEnablerFunc) {
	debug = func(level zapcore.Level) bool { return level == zapcore.DebugLevel }
	info = func(level zapcore.Level) bool {
		return level == zapcore.InfoLevel || level == zapcore.WarnLevel
	}
	return
}

// core 核心
func core(
	consoleConf zapcore.EncoderConfig, fileConf zapcore.EncoderConfig) (
	console zapcore.Core, fileDebug zapcore.Core, fileInfo zapcore.Core, fileError zapcore.Core,
) {

	console = zapcore.NewCore(
		zapcore.NewConsoleEncoder(consoleConf),
		zapcore.AddSync(os.Stdout),
		zap.InfoLevel,
	)

	t := time.Now()
	jsonEncoder := zapcore.NewJSONEncoder(fileConf)
	debugLevel, infoLevel := level()

	list := fileCores{
		{&fileDebug, "debug.log", 3, debugLevel},
		{&fileInfo, "info.log", 31, infoLevel},
		{&fileError, "error.log", 64, zapcore.ErrorLevel},
	}

	for _, v := range list {
		*v.v = zapcore.NewCore(
			jsonEncoder,
			zapcore.AddSync(&lumberjack.Logger{Filename: fileName(t) + v.name, MaxAge: v.age, LocalTime: true}),
			v.l,
		)
	}
	return
}

// init 设置 Logger
func init() {
	logger = zap.New(zapcore.NewTee(core(config())), zap.AddCaller(), zap.AddCallerSkip(1))
	defer logger.Sync()
}
