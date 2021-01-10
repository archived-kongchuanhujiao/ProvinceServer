package wenda

import (
	"coding.net/kongchuanhujiao/server/internal/app/client/clientmsg"
	"coding.net/kongchuanhujiao/server/internal/app/datahub/datahubpkg/wenda"
	"fmt"
	"github.com/CatchZeng/dingtalk"
	"strconv"
)

// convertToMarkDown 将问题数据转换为钉钉 Markdown 消息
// FIXME 详见 apis.go@PostPushcenter()
func convertToMarkdown(tab *wenda.QuestionsTab) *dingtalk.MarkdownMessage {
	builder := dingtalk.NewMarkdownMessage()

	t, c := getQuestionDetail(tab, true)

	builder.Markdown.Title = t
	builder.Markdown.Text = c

	return builder
}

// convertToChain 将问题数据转换为消息链
// TODO 详见 apis.go@PostPushcenter()
func convertToChain(tab *wenda.QuestionsTab) (m *clientmsg.Message) {
	t, c := getQuestionDetail(tab, false)
	m = clientmsg.NewTextMessage(t + "\n" + c)

	return
}

// getQuestionDetail 获取问题信息的字符串形式
func getQuestionDetail(tab *wenda.QuestionsTab, useMarkdown bool) (title string, content string) {
	title = "问题 #" + strconv.Itoa(int(tab.ID)) + " 的数据"
	// FIXME 等待数据库 API
	content = getQuestionSummary(tab, useMarkdown)

	return
}

// getQuestionSummary 获取题目概要
func getQuestionSummary(tab *wenda.QuestionsTab, isMarkdown bool) (sum string) {
	sum = tab.Question + " 选项: " + tab.Options

	if len(sum) > 20 {
		sum = sum[0:20] + "..."
	}

	if !isMarkdown {
		template := "## #{id} 详细信息\n\n> 正确人数 > {number} 人\n> 正确率 > {percent}%\n> 易错选项 > {option}\n> 最快答对的同学 > {name}"
		sum += fmt.Sprintf(template, tab.ID, "人数", "正确率", "易错选项", "最速同学")
	} else {
		template := "#{id} 详细信息\n\n 正确人数 > {number} 人\n 正确率 > {percent}%\n 易错选项 > {option}\n 最快答对的同学 > {name}"
		sum += fmt.Sprintf(template, tab.ID, "人数", "正确率", "易错选项", "最速同学")
	}

	return
}

/*
Markdown 模板

## #{id} 详细信息

> 正确人数 > {number} 人
> 正确率 > {percent}%
> 易错选项 > {option}
> 最快答对的同学 > {name}
*/
