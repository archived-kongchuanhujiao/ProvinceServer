package wenda

import (
	"errors"
	"fmt"
	"github.com/kongchuanhujiao/server/internal/app/datahub/pkg/accounts"
	"github.com/kongchuanhujiao/server/internal/app/datahub/pkg/wenda"
	"github.com/kongchuanhujiao/server/internal/pkg/logger"
	"sort"
	"strconv"
	"strings"

	"github.com/kongchuanhujiao/server/internal/app/client/message"
	public "github.com/kongchuanhujiao/server/internal/app/datahub/public/wenda"

	"github.com/CatchZeng/dingtalk"
)

// PushDigestData 推送答题数据
func PushDigestData(platform string, tab *public.QuestionsTab) (err error) {
	switch platform {
	case "dingtalk":
		acc, errr := accounts.SelectAccount(tab.Creator, 0)

		if errr != nil {
			return errr
		}

		if acc[0].Token == "" || acc[0].Push == "" {
			return
		}

		if len(acc) < 1 {
			return errors.New("couldn't find account by given name")
		}

		err = PushDigestToDingtalk(acc[0].Token, acc[0].Push, convertToDTMessage(tab))

		return
	case "qq":
		return
	default:
		return errors.New("Unknown push platform: " + platform)
	}
}

// convertToMarkDown 将问题数据转换为钉钉 Markdown 消息
func convertToDTMessage(tab *public.QuestionsTab) *dingtalk.MarkdownMessage {
	builder := dingtalk.NewMarkdownMessage()
	builder.Markdown.Title = "答题数据："
	builder.Markdown.Text = digestQuestionData(tab, true)
	return builder
}

// convertToChain 将问题数据转换为消息链
func convertToChain(tab *public.QuestionsTab) *message.Message {
	return message.NewTextMessage(digestQuestionData(tab, false))
}

// digestQuestionData 摘要答题数据
func digestQuestionData(tab *public.QuestionsTab, isMarkdown bool) (sum string) {
	sum = digestQuestion(tab)
	template := ""

	calc, err := wenda.CalculateResult(tab.ID)

	if err != nil {
		return
	}

	if !isMarkdown {
		template = "## #%v 详细信息  \n  \n> 正确人数 > %v 人  \n> 正确率 > %v  \n> 易错选项 > %v  \n> 最快答对同学 %v"
	} else {
		template = "#%v 详细信息\n\n 正确人数 > %v 人\n 正确率 > %v\n 易错选项 > %v\n> 最快答对同学 %v"
	}
	sum += fmt.Sprintf(template, tab.ID, calc.Count, getRightRate(calc), getMostWrongOption(calc.Wrong), getFastestAnswerUser(tab))
	return
}

// digestQuestion 摘要题干
func digestQuestion(q *public.QuestionsTab) (s string) {

	var questionText string
	for _, v := range q.Topic.Question {
		if v.Type == "img" {
			questionText += "[图片]"
		}
		questionText += v.Data
	}

	var optionsText string

	abc := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I"} // FIXME 弄到全局public去
	for k, v := range q.Topic.Options {
		optionsText += abc[k] + ". " + v + "\n"
	}

	optionsText = strings.TrimSuffix(optionsText, "\n")

	s = "题目: " + questionText + " 选项：" + optionsText
	if len(s) > 25 {
		s = s[0:25] + "..."
	}
	return
}

// PushDigestToQQ TODO 推送摘要至QQ平台
func PushDigestToQQ() (err error) {
	return
}

// PushDigestToDingtalk 推送摘要至钉钉平台
func PushDigestToDingtalk(accessToken string, secret string, md dingtalk.Message) (err error) {
	logger.Info("正在推送答题概要至钉钉")

	client := dingtalk.Client{
		AccessToken: accessToken,
		Secret:      secret,
	}
	_, err = client.Send(md)
	return
}

// wrongFieldWrapper 包装类
type wrongFieldWrapper []public.ResultWrongField

func (a wrongFieldWrapper) Len() int {
	return len(a)
}
func (a wrongFieldWrapper) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a wrongFieldWrapper) Less(i, j int) bool {
	return len(a[j].Value) < len(a[i].Value)
}

// getRightRate 获取答题正答率
func getRightRate(result *public.Result) string {
	return fmt.Sprintf("%.2f", float32(len(result.Right))/float32(result.Count)*100) + "%"
}

// getMostWrongOption 获取易错选项
func getMostWrongOption(wrong []public.ResultWrongField) string {
	wrap := wrongFieldWrapper(wrong)
	sort.Sort(wrap)
	return wrap[0].Type
}

func getFastestAnswerUser(tab *public.QuestionsTab) (name string) {
	ans, err := wenda.SelectAnswers(tab.ID)

	if err != nil {
		return
	}

	for _, an := range ans {
		if an.Answer == tab.Topic.Key {
			return strconv.FormatUint(an.QQ, 10)
		}
	}

	return
}
