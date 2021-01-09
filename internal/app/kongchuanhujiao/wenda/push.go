package wenda

import (
	"coding.net/kongchuanhujiao/server/internal/app/datahub/datahubpkg/wenda"
	"github.com/CatchZeng/dingtalk"
	"github.com/Mrs4s/MiraiGo/message"
	"strconv"
)

// convertToMarkDown 将问题数据转换为钉钉 Markdown 消息
// FIXME 详见 apis.go@PostPushcenter()
func convertToMarkdown(tab *wenda.QuestionsTab) *dingtalk.MarkdownMessage {
	builder := dingtalk.NewMarkdownMessage()

	t, c := getQuestionDetail(tab)

	builder.Markdown.Title = t
	builder.Markdown.Text = c

	return builder
}

// convertToChain 将问题数据转换为 MiraiGO 消息链
// TODO 详见 apis.go@PostPushcenter()
func convertToChain(tab *wenda.QuestionsTab) (m *message.SendingMessage) {
	m = &message.SendingMessage{
		Elements: []message.IMessageElement{},
	}

	t, c := getQuestionDetail(tab)

	m.Elements = append(m.Elements, message.NewText(t+"\n"+c))

	return
}

// getQuestionDetail 获取问题信息的字符串形式
func getQuestionDetail(tab *wenda.QuestionsTab) (title string, content string) {
	title = "问题 #" + strconv.Itoa(int(tab.ID)) + " 的数据"
	// FIXME 等待数据库 API
	content = "问题状态: " + getStatusName(tab.Status) + "\n"

	return
}

// getStatusName 获取状态名
func getStatusName(status uint8) string {
	switch status {
	case 0:
		return "准备作答"
	case 1:
		return "允许作答"
	case 2:
		return "停止作答"
	default:
		return "未知"
	}
}
