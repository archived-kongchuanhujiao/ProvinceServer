package wenda

import (
	"fmt"
	"strings"

	"github.com/kongchuanhujiao/server/internal/app/client"
	"github.com/kongchuanhujiao/server/internal/app/client/message"
	"github.com/kongchuanhujiao/server/internal/app/datahub/pkg/wenda"
	public "github.com/kongchuanhujiao/server/internal/app/datahub/public/wenda"
)

// HandleAnswer 处理回答
func HandleAnswer(m *message.Message) {

	qid := wenda.GetActiveGroup(m.Target.Group.ID)
	if qid == 0 {
		return
	}

	ans, ok := m.Chain[0].(*message.Text)
	if !ok {
		return
	}
	answer := ans.Content

	q := wenda.GetCaches(qid)
	for _, v := range q.Answers {
		if v.QQ == m.Target.ID {
			return
		}
	}

	switch q.Questions.Type {
	case 0: // 选择题
		if !checkAnswerForSelect(answer) {
			return
		}

		var (
			ans  = strings.ToUpper(answer)
			mark string
		)
		if ans != wenda.GetCaches(qid).Questions.Topic.Key {
			mark = ans
		}
		_ = wenda.InsertAnswer(&public.AnswersTab{
			Question: qid,
			QQ:       m.Target.ID,
			Answer:   ans,
			Mark:     mark,
		})

		calc, _ := wenda.CalculateResult(wenda.GetCaches(qid).Questions.ID)
		wenda.PushData(q.Questions.ID, calc)

	case 1: // 填空题

	case 2: // 多选题

	case 3: // 简答题
		if !checkAnswerForFill(answer) {
			return
		}
		_ = wenda.InsertAnswer(&public.AnswersTab{
			Question: qid,
			QQ:       m.Target.ID,
			Answer:   strings.TrimPrefix(answer, "#"),
		})
	}
}

// HandleTest 处理测试
func HandleTest(m *message.Message) {

	t, ok := m.Chain[0].(*message.Text)
	if !ok {
		return
	}

	switch t.Content {
	case "你好":
		client.GetClient().SendMessage(
			message.NewAtMessage(m.Target.ID).AddText("你好").SetGroupTarget(m.Target.Group),
		)
	case "活动的群":
		client.GetClient().SendMessage(
			message.NewTextMessage(fmt.Sprintln(wenda.GetAllActiveGroup())).SetGroupTarget(m.Target.Group),
		)
	}
}

// HandleWrongQuestion 处理添加错题逻辑
func HandleWrongQuestion(m *message.Message) {
	t, ok := m.Chain[0].(*message.Text)
	if !ok {
		return
	}

	if !strings.HasPrefix(t.Content, "/") {
		return
	}

	switch strings.TrimPrefix(t.Content, "/") {
	case "ct", "错题":
		args := strings.Split(t.Content, " ")

		defaultMsg := message.NewAtMessage(m.Target.ID).
			AddText("/ct add 添加错题\n" +
				"/ct del 删除错题\n" +
				"/ct zc 查看错题").
			SetGroupTarget(m.Target.Group)

		if len(args) == 1 {
			client.GetClient().SendMessage(defaultMsg)
		} else {
			switch args[1] {
			case "add":

			case "del":

			case "zc":

			default:
				client.GetClient().SendMessage(defaultMsg)
			}
		}
	}
}
