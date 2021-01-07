package dingtalk

import (
	"github.com/CatchZeng/dingtalk"
)

// TODO 这里引入第三方包
// TODO 这里不能调用其他的包，只能调用第三方包，否则软件架构会混乱。所有需要的数据在进来前就要准备好，不能从这里获取其他数据信息
// 另请查询 [apis.go#PostPUSHCENTER()]

func Push(accessToken string, secret string, content string, atMobiles []string, isAtAll bool) error {
	msg := dingtalk.NewTextMessage().SetContent(content)

	if len(atMobiles) > 0 {
		msg.SetAt(atMobiles, isAtAll)
	}

	client := dingtalk.Client{
		AccessToken: accessToken,
		Secret:      secret,
	}
	_, err := client.Send(msg)

	return err
}
