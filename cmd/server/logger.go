package main

import (
	"os"
	"time"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var logger *zap.Logger

// init 初始化 logger
func init() {

	conf := zap.NewProductionEncoderConfig()
	conf.EncodeTime = zapcore.TimeEncoderOfLayout("01-02T15:04:05.999999999")

	console := conf
	file := conf

	console.EncodeLevel = zapcore.CapitalColorLevelEncoder
	file.EncodeCaller = zapcore.FullCallerEncoder

	lumber := &lumberjack.Logger{
		Filename:  time.Now().Format(".kongchuanhujiao/logs/2006-01-02.log"),
		MaxAge:    7,
		LocalTime: true,
		Compress:  true,
	}

	logger = zap.New(
		zapcore.NewTee(
			zapcore.NewCore(zapcore.NewJSONEncoder(file), zapcore.AddSync(lumber), zapcore.WarnLevel),
			zapcore.NewCore(zapcore.NewConsoleEncoder(console), zapcore.AddSync(os.Stdout), zapcore.DebugLevel),
		),
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)
	job(lumber)
	logger.Named("日志").Debug("初始化成功")

	zap.ReplaceGlobals(logger)
}

// job 日志轮转
func job(l *lumberjack.Logger) {

	c := cron.New()
	_, _ = c.AddFunc("0 0 * * *", func() {
		if err := l.Rotate(); err != nil {
			logger.Named("日志").Warn("无法轮转", zap.Error(err))
			return
		}
	})
	c.Start()

}
