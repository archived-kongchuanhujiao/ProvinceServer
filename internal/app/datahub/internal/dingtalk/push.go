package dingtalk

import (
	"github.com/CatchZeng/dingtalk"
)

// 这里不能调用其他的包，只能调用第三方包，否则软件架构会混乱。所有需要的数据在进来前就要准备好，不能从这里获取其他数据信息/
// 另请查询 [apis.go#PostPushcenter()]

// Push 推送一个任意样式的消息至钉钉平台.
func Push(accessToken string, secret string, md *dingtalk.Message) error {

	client := dingtalk.Client{
		AccessToken: accessToken,
		Secret:      secret,
	}
	_, err := client.Send(*md)

	return err
}
