package accounts

import (
	"errors"
	"math/rand"
	"strconv"
	"time"

	"coding.net/kongchuanhujiao/server/internal/app/client"
	"coding.net/kongchuanhujiao/server/internal/app/client/clientmsg"
	"coding.net/kongchuanhujiao/server/internal/app/datahub/datahubpkg/accounts"
	"coding.net/kongchuanhujiao/server/internal/app/kongchuanhujiao"
)

// TODO https://qianjunakasumi.coding.net/p/kongchuanhujiao/requirements/issues/15/detail
/*
TODO
   MVC API 参考 wenda 包接入
   注册账号 API 实现步骤：
   接收POST注册账号 POST /apis/accounts/register/verifier 要求QQ字段 QQ：uint64
   向QQ号发送好友消息（消息内容：验证码,4位，有效期5min）
   客户端 POST /apis/accounts/register  要求 QQ 、验证码和ID 字段  验证码：uint8 id:string
   验证 验证码：通过map[uint64(QQ)]uint8(验证码) // 对比，如果一致则注册成功
   向数据库写入，数据API我后面会写
*/

type (
	APIs struct{} // APIs 账号 APIs

	PostCodeReq struct{ ID string } // PostCodeReq 验证码发送
)

var code map[string]string = map[string]string{} // code 验证码

// 验证码。
// POST apis/accounts/code
func (a *APIs) PostCode(v *PostCodeReq) *kongchuanhujiao.Response {
	err := sendCode(v.ID)
	if err != nil {
		return &kongchuanhujiao.Response{Status: 1, Message: "服务器错误"}
	}
	return &kongchuanhujiao.Response{Message: "ok"}
}

// 登录。
// POST apis/accounts/login
func (a *APIs) PostLogin() {

}

func (a *APIs) GetLogin() {

}

// sendCode 发送验证码
func sendCode(id string) error {
	a, err := accounts.SelectAccount(id, 0)
	if err != nil {
		return err
	}

	if len(a) == 0 {
		return errors.New("账号不存在")
	}

	rand.Seed(time.Now().UnixNano())
	c := strconv.FormatFloat(rand.Float64(), 'f', -1, 64)[2:6]

	m := clientmsg.NewTextMessage("您的验证码是： " + c + " 。任何人都不会索要验证码！")
	client.GetClient().SendMessage(m.SetTarget(&clientmsg.Target{ID: a[0].QQ}))
	code[id] = c
	return nil
}
