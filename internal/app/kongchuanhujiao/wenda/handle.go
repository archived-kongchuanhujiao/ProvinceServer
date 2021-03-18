package wenda

import (
	"strings"

	"github.com/kongchuanhujiao/server/internal/app/client"
	"github.com/kongchuanhujiao/server/internal/app/client/message"
	"github.com/kongchuanhujiao/server/internal/app/datahub/pkg/wenda"
	public "github.com/kongchuanhujiao/server/internal/app/datahub/public/wenda"
	"github.com/kongchuanhujiao/server/internal/pkg/config"
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
