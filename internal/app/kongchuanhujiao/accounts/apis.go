package accounts

import (
	"errors"
	"math/rand"
	"strconv"
	"time"

	"github.com/kongchuanhujiao/server/internal/app/client"
	"github.com/kongchuanhujiao/server/internal/app/client/message"
	"github.com/kongchuanhujiao/server/internal/app/datahub/pkg/accounts"
	"github.com/kongchuanhujiao/server/internal/app/kongchuanhujiao"
	"github.com/kongchuanhujiao/server/internal/pkg/configs"
	"github.com/kongchuanhujiao/server/internal/pkg/logger"

	"github.com/iris-contrib/middleware/jwt"
	"github.com/kataras/iris/v12/context"
	"go.uber.org/zap"
)

type (
	APIs struct{} // APIs 账号 APIs

	PostCodeReq struct{ ID string } // PostCodeReq 验证码发送

	PostLoginReq struct { // PostLoginReq 登录验证
		ID   string // 唯一识别码
		Code string // 验证码
	}
)

var code = map[string]string{} // code 验证码
var conf = configs.GetConfigs()

// 验证码。
// POST apis/accounts/code
func (a *APIs) PostCode(v *PostCodeReq) *kongchuanhujiao.Response {
	err := sendCode(v.ID)
	if err != nil {
		return &kongchuanhujiao.Response{Status: 1, Message: err.Error()}
	}
	return &kongchuanhujiao.Response{Message: "ok"}
}

// 登录。
// POST apis/accounts/login
func (a *APIs) PostLogin(v *PostLoginReq, c *context.Context) *kongchuanhujiao.Response {

	if v.Code != code[v.ID] || v.Code == "" {
		return &kongchuanhujiao.Response{Status: 1, Message: "验证码有误"}
	}

	now := time.Now()
	t, err := jwt.NewTokenWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"iss": conf.JWT.Iss,
		"sub": v.ID,
		"exp": now.AddDate(0, 1, 0).Unix(),
		"nbf": now.Unix(),
		"iat": now.Unix(),
	}).SignedString(conf.JWT.Key)
	if err != nil {
		logger.Error("生成 JWT Token 失败", zap.Error(err))
		return &kongchuanhujiao.Response{Status: 1, Message: "服务器错误"}
	}

	return &kongchuanhujiao.Response{Message: t}
}

// sendCode 发送验证码
func sendCode(id string) (err error) {
	a, err := accounts.SelectAccount(id, 0)
	if err != nil {
		logger.Error("发送验证码错误", zap.Error(err))
		return errors.New("服务器错误")
	}

	if len(a) == 0 {
		return errors.New("账号不存在")
	}

	rand.Seed(time.Now().UnixNano())
	c := strconv.FormatFloat(rand.Float64(), 'f', -1, 64)[2:6]

	m := message.NewTextMessage("您的验证码是：" + c + "，请勿泄露给他人。")
	client.GetClient().SendMessage(m.SetTarget(&message.Target{ID: a[0].QQ}))
	code[id] = c
	return
}
