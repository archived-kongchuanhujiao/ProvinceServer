package wenda

import (
	"coding.net/kongchuanhujiao/server/internal/app/client/clientmsg"
	"coding.net/kongchuanhujiao/server/internal/app/kongchuanhujiao/public/wendapkg"
	"encoding/json"
	"fmt"
	"github.com/CatchZeng/dingtalk"
	"strings"
)

// convertToMarkDown 将问题数据转换为钉钉 Markdown 消息
// FIXME 详见 apis.go@PostPushcenter()
func ConvertToDTMessage(tab *wendapkg.QuestionsTab) *dingtalk.MarkdownMessage {
	builder := dingtalk.NewMarkdownMessage()
	builder.Markdown.Title = "答题数据："
	builder.Markdown.Text = digestQuestionData(tab, true)
	return builder
}

// convertToChain 将问题数据转换为消息链
// TODO 详见 apis.go@PostPushcenter()
func ConvertToChain(tab *wendapkg.QuestionsTab) *clientmsg.Message {
	return clientmsg.NewTextMessage(digestQuestionData(tab, false))
}

// digestQuestionData 摘要答题数据
func digestQuestionData(tab *wendapkg.QuestionsTab, isMarkdown bool) (sum string) {
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
func digestQuestion(q *wendapkg.QuestionsTab) (s string) {
	// FIXME Question 和 Options 均为json，需要特殊解析
	qs := wendapkg.Question{}
	var os []wendapkg.Option

	err := json.Unmarshal([]byte(q.Question), &qs)
	err = json.Unmarshal([]byte(q.Options), &os)

	if err != nil {
		return "无法解析"
	}

	var optionsText string

	for _, op := range os {
		optionsText = op.Type + ": " + op.Body + "\n"
	}

	optionsText = strings.TrimSuffix(optionsText, "\n")

	s = "题目: " + qs.Text + " 选项：" + optionsText
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
