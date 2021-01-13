package datahubpkg

import (
	"coding.net/kongchuanhujiao/server/internal/app/datahub/internal/dingtalk"
	dingtalk2 "github.com/CatchZeng/dingtalk"
)

// PushMessage 推送一个任意样式的消息至钉钉平台.
func PushMessage(accessToken string, secret string, md *dingtalk2.Message) error {
	return dingtalk.Push(accessToken, secret, md)
}
