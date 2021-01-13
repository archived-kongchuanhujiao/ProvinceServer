package datahubpkg

import (
	dt "coding.net/kongchuanhujiao/server/internal/app/datahub/internal/dingtalk"
	"github.com/CatchZeng/dingtalk"
)

// PushMessage 推送一个任意样式的消息至钉钉平台.
func PushMessage(accessToken string, secret string, md dingtalk.Message) error {
	return dt.Push(accessToken, secret, md)
}
