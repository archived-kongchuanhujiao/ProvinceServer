package accounts

import (
	"errors"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"coding.net/kongchuanhujiao/server/internal/app/client"
	"coding.net/kongchuanhujiao/server/internal/app/client/clientmsg"
	"coding.net/kongchuanhujiao/server/internal/app/datahub/datahubpkg/accounts"
	"coding.net/kongchuanhujiao/server/internal/app/kongchuanhujiao"
	"coding.net/kongchuanhujiao/server/internal/pkg/logger"

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

	c.SetCookie(&http.Cookie{
		Name: "account", Value: v.ID, Path: "/", Expires: time.Now().AddDate(0, 1, 0),
	})

	// TODO 生成Token

	return &kongchuanhujiao.Response{Message: "ok"}

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

	m := clientmsg.NewTextMessage("您的验证码是： " + c + " 。任何人都不会索要验证码！")
	client.GetClient().SendMessage(m.SetTarget(&clientmsg.Target{ID: a[0].QQ}))
	code[id] = c
	return
}
