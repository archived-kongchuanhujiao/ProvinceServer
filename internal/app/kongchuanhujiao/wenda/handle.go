package wenda

import (
	"fmt"
	"github.com/kongchuanhujiao/server/internal/app/datahub/pkg/cuoti"
	"github.com/kongchuanhujiao/server/internal/pkg/config"
	"strings"

	"github.com/kongchuanhujiao/server/internal/app/client"
	"github.com/kongchuanhujiao/server/internal/app/client/message"
	"github.com/kongchuanhujiao/server/internal/app/datahub/pkg/wenda"
	public "github.com/kongchuanhujiao/server/internal/app/datahub/public/wenda"
)

var sessionPool []uint64

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
	case "空传你好":
		client.GetClient().SendMessage(
			message.NewAtMessage(m.Target.ID).AddText("你好").SetGroupTarget(m.Target.Group),
		)
	case "空传版本":
		client.GetClient().SendMessage(
			message.NewTextMessage(config.Commit).SetGroupTarget(m.Target.Group),
		)
	default:
		if strings.HasPrefix(t.Content, "空传分词 ") {
			client.GetClient().SendMessage(
				message.NewTextMessage(
					strings.Join(client.GetClient().ExtractWords(t.Content), "／"),
				).SetGroupTarget(m.Target.Group),
			)
		}
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
				"/ct ck 查看错题").
			SetGroupTarget(m.Target.Group)

		if len(args) == 1 {
			client.GetClient().SendMessage(defaultMsg)
		} else {
			switch args[1] {
			case "add":
				for _, u := range sessionPool {
					if u == m.Target.ID {
						sessionPool = append(sessionPool, m.Target.ID)
						client.GetClient().SendMessage(message.NewAtMessage(m.Target.ID).AddText("请发送欲添加错题的题目:"))
						return
					}
				}

				client.GetClient().SendMessage(message.NewAtMessage(m.Target.ID).AddText("你还有正在进行添加的错题!"))
			case "del":
				return
			case "ck":
				wq, err := cuoti.SelectWrongQuestions(0, uint32(m.Target.ID))

				if err != nil {
					client.GetClient().SendMessage(message.NewAtMessage(m.Target.ID).AddText("发生了意外错误, 无法查询错题列表."))
					return
				}

				client.GetClient().SendMessage(message.NewAtMessage(m.Target.ID).AddText(fmt.Sprint(wq)))
			default:
				client.GetClient().SendMessage(defaultMsg)
			}
		}
	}
}
