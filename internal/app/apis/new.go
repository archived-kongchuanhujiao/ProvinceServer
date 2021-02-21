package apis

import (
	"github.com/kongchuanhujiao/server/internal/app/kongchuanhujiao/accounts"
	"github.com/kongchuanhujiao/server/internal/app/kongchuanhujiao/wenda"
	"github.com/kongchuanhujiao/server/internal/pkg/logger"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/recover"
	"github.com/kataras/iris/v12/mvc"
	"go.uber.org/zap"
)

// StartApis 启动 APIs
func StartApis() {

	app := iris.New()
	app.Use(recover.New())
	APIs := mvc.New(app.Party("apis/"))

	APIs.Party("accounts/").Handle(new(accounts.APIs))
	APIs.Party("wenda/").Handle(new(wenda.APIs))

	go func() {

		loggerr := logger.Named("APIs")
		loggerr.Info("启动服务中")
		if err := app.Listen(
			":5245",
			iris.WithoutInterruptHandler,
			iris.WithoutStartupLog,
			iris.WithCharset("utf-8"),
		); err != nil {
			loggerr.Panic("监听端口失败", zap.Error(err))
		}
	}()
}
