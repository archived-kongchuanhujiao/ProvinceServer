package wenda

import (
	"coding.net/kongchuanhujiao/server/internal/app/datahub/datahubpkg/wenda"
	"github.com/CatchZeng/dingtalk"
	"github.com/Mrs4s/MiraiGo/message"
)

// convertToMarkDown 将问题数据转换为钉钉 Markdown 消息
// TODO 详见 apis.go@PostPushcenter()
func convertToMarkdown(tab *wenda.QuestionsTab) *dingtalk.MarkdownMessage {
	builder := dingtalk.NewMarkdownMessage()

	builder.Markdown.Text = ""

	return builder
}

// convertToChain 将问题数据转换为 MiraiGO 消息链
// TODO 详见 apis.go@PostPushcenter()
func convertToChain(tab *wenda.QuestionsTab) *message.SendingMessage {
	return nil
}
