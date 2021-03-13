package api

import (
	"github.com/kongchuanhujiao/server/internal/app/kongchuanhujiao/account"
	"github.com/kongchuanhujiao/server/internal/app/kongchuanhujiao/wenda"
	"github.com/kongchuanhujiao/server/internal/pkg/config"
	"github.com/kongchuanhujiao/server/internal/pkg/logger"

	"github.com/iris-contrib/middleware/jwt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/recover"
	"github.com/kataras/iris/v12/mvc"
	"go.uber.org/zap"
)

var jc = jwt.New(jwt.Config{
	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
		return config.GetJWTConf().Key.Public(), nil
	},
	Extractor:     jwt.FromFirst(jwt.FromAuthHeader, jwt.FromParameter("token")),
	SigningMethod: jwt.SigningMethodES256,
})

// jwtMiddleware 处理鉴权
func jwtMiddleware(c iris.Context) {

	if err := jc.CheckJWT(c); err != nil {
		c.StatusCode(401)
		logger.Warn("未授权的访问", zap.Error(err), zap.String("客户", c.RemoteAddr()))
		return
	}

	t := c.Values().Get("jwt").(*jwt.Token)
	err := t.Claims.Valid()
	if err != nil {
		c.StatusCode(401)
		logger.Warn("无效的 Token", zap.Error(err), zap.String("客户", c.RemoteAddr()))
		return
	}

	cla := t.Claims.(jwt.MapClaims)
	if cla["iss"] != config.GetJWTConf().Iss {
		c.StatusCode(401)
		logger.Warn("危险的 Token", zap.Error(err), zap.String("客户", c.RemoteAddr()))
		return
	}

	c.Values().Set("account", cla["sub"])
	c.Next()
}

// StartApis 启动 APIs
func StartApis() {

	app := iris.New()
	app.Use(recover.New())
	APIs := mvc.New(app.Party("apis/"))

	APIs.Party("accounts/").Handle(new(account.APIs))
	APIs.Party("wenda/", jwtMiddleware).Handle(new(wenda.APIs))

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
