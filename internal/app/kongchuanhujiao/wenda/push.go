package wenda

import (
	"fmt"
	"strings"

	"coding.net/kongchuanhujiao/server/internal/app/client/message"
	"coding.net/kongchuanhujiao/server/internal/app/datahub/pkg/wenda"

	"github.com/CatchZeng/dingtalk"
)

// convertToMarkDown 将问题数据转换为钉钉 Markdown 消息
// FIXME 详见 apis.go@PostPushcenter()
func ConvertToDTMessage(tab *wenda.QuestionsTab) *dingtalk.MarkdownMessage {
	builder := dingtalk.NewMarkdownMessage()
	builder.Markdown.Title = "答题数据："
	builder.Markdown.Text = digestQuestionData(tab, true)
	return builder
}

// convertToChain 将问题数据转换为消息链
// TODO 详见 apis.go@PostPushcenter()
func ConvertToChain(tab *wenda.QuestionsTab) *message.Message {
	return message.NewTextMessage(digestQuestionData(tab, false))
}

// digestQuestionData 摘要答题数据
func digestQuestionData(tab *wenda.QuestionsTab, isMarkdown bool) (sum string) {
	sum = digestQuestion(tab)
	template := ""
	if !isMarkdown {
		template = "## #%v 详细信息  \n  \n> 正确人数 > %v 人  \n> 正确率 > %v  \n> 易错选项 > %v  \n> 最快答对的同学 > %v"
	} else {
		template = "#%v 详细信息\n\n 正确人数 > %v 人\n 正确率 > %v\n 易错选项 > %v\n 最快答对的同学 > %v"
	}
	sum += fmt.Sprintf(template, tab.ID, "人数", "正确率", "易错选项", "最速同学")
	return
}

// digestQuestion 摘要题干
func digestQuestion(q *wenda.QuestionsTab) (s string) {

	var questionText string
	for _, v := range q.Question {
		if v.Type == "img" {
			questionText += "[图片]"
		}
		questionText += v.Data
	}

	var optionsText string

	abc := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I"} // FIXME 弄到全局public去
	for k, v := range q.Options {
		optionsText += abc[k] + ". " + v + "\n"
	}

	optionsText = strings.TrimSuffix(optionsText, "\n")

	s = "题目: " + questionText + " 选项：" + optionsText
	if len(s) > 20 {
		s = s[0:20] + "..."
	}
	return
}

// PushDigestToQQ TODO 推送摘要至QQ平台
func PushDigestToQQ() (err error) {
	return
}

// PushDigestToDingtalk 推送摘要至钉钉平台
func PushDigestToDingtalk(accessToken string, secret string, md dingtalk.Message) (err error) {
	client := dingtalk.Client{
		AccessToken: accessToken,
		Secret:      secret,
	}
	_, err = client.Send(md)
	return
}
