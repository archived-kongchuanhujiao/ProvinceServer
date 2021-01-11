package datahubpkg

import (
	"coding.net/kongchuanhujiao/server/internal/app/datahub/internal/dingtalk"
	dingtalk2 "github.com/CatchZeng/dingtalk"
)

// PushMessage 推送一个纯文本样式的消息到钉钉.
func PushMessage(accessToken string, secret string, content string, atMobiles []string, isAtAll bool) error {
	return dingtalk.Push(accessToken, secret, content, atMobiles, isAtAll)
}

// PushMessageMD 推送一个 Markdown 样式的消息到钉钉.
func PushMessageMD(accessToken string, secret string, md *dingtalk2.MarkdownMessage) error {
	return dingtalk.PushMD(accessToken, secret, md)
}
