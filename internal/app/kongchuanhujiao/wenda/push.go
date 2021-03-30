package wenda

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/kongchuanhujiao/server/internal/app/client"
	"github.com/kongchuanhujiao/server/internal/app/client/message"
	"github.com/kongchuanhujiao/server/internal/app/datahub/pkg/account"
	"github.com/kongchuanhujiao/server/internal/app/datahub/pkg/wenda"
	public "github.com/kongchuanhujiao/server/internal/app/datahub/public/wenda"

	"github.com/CatchZeng/dingtalk"
	"go.uber.org/zap"
)

// PushDigestData 推送答题数据
func PushDigestData(tab *public.QuestionsTab) (err error) {

	acc, err := account.SelectAccount(tab.Creator, 0)

	if err != nil {
		return
	}
	if len(acc) < 1 {
		return errors.New("couldn't find account by given name")
	}

	for _, p := range acc[0].Push {

		switch p.Platform {
		case "dingtalk":
			err = PushDigestToDingtalk(p.Key, p.Secret, convertToDTMessage(tab))
			// FIXME 错误处理
		case "qq":
			t, err := strconv.ParseUint(p.Key, 10, 64)
			if err != nil {
				return err
			}

			PushDigestToQQ(t, convertToChain(tab))
		}
	}
	return
}

// convertToDTMessage 将问题数据转换为钉钉 Markdown 消息
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
	var temp string

	calc, err := wenda.CalculateResult(tab.ID)

	if err != nil {
		return
	}

	if isEmptyQuestion(calc) {
		return fmt.Sprintf("题目 #%v 暂无数据", tab.ID)
	}

	temp = "题目 #%v 详细信息\n\n 正确人数 > %v 人\n 正确率 > %v\n 易错选项 > %v\n最快答对同学 > %v"

	if isMarkdown {
		temp = strings.ReplaceAll(temp, "\n", "  \n")
	}

	sum += fmt.Sprintf(temp, tab.ID, calc.Count, getRightRate(calc), getMostWrongOption(calc.Wrong), getFastestAnswerUser(tab))
	return
}

// PushDigestToQQ 推送摘要至QQ平台
func PushDigestToQQ(target uint64, data *message.Message) {
	zap.L().Info("正在推送答题概要至QQ")

	client.GetClient().SendMessage(data.SetTarget(&message.Target{ID: target}))
}

// PushDigestToDingtalk 推送摘要至钉钉平台
func PushDigestToDingtalk(accessToken string, secret string, md dingtalk.Message) (err error) {

	zap.L().Info("正在推送答题概要至钉钉")

	c := dingtalk.Client{
		AccessToken: accessToken,
		Secret:      secret,
	}
	_, err = c.Send(md)
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
	if wrap.Len() == 0 {
		return "无"
	} else {
		return wrap[0].Type
	}
}

// getFastestAnswerUser 获取最快答对用户
func getFastestAnswerUser(tab *public.QuestionsTab) (name string) {
	ans, err := wenda.SelectAnswers(tab.ID)

	if err != nil {
		return
	}

	g := client.GetClient().GetGroupMembers(tab.Topic.Target)

	for _, an := range ans {
		if an.Answer == tab.Topic.Key {
			return (*g)[an.QQ]
		}
	}

	name = "无"

	return
}

// isEmptyQuestion 该问题是否无人回答过
func isEmptyQuestion(r *public.Result) bool {
	return r.Count == 0
}
