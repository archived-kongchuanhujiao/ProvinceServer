package wenda

import (
	"fmt"
	"strings"

	"coding.net/kongchuanhujiao/server/internal/app/client"
	"coding.net/kongchuanhujiao/server/internal/app/client/message"
	"coding.net/kongchuanhujiao/server/internal/app/datahub/pkg/wenda"
)

// HandleAnswer 处理回答
func HandleAnswer(m *message.Message) {

	qid, ok := wenda.ActiveGroup[m.Target.Group.ID]
	if !ok {
		return
	}

	ans, ok := m.Chain[0].(*message.Text)
	if !ok {
		return
	}
	answer := ans.Content

	q := wenda.Caches[qid]
	for _, v := range q.Answers {
		if v.QQ == m.Target.ID {
			return
		}
	}

	switch q.Questions.Type {
	case 0, 1: // 选择题、填空题
		if !checkAnswerForSelect(answer) {
			return
		}
		_ = wenda.InsertAnswer(&wenda.AnswersTab{
			Question: qid,
			QQ:       m.Target.ID,
			Answer:   strings.ToUpper(answer),
		})
	case 2: // 多选题

	case 3: // 简答题
		if !checkAnswerForFill(answer) {
			return
		}
		_ = wenda.InsertAnswer(&wenda.AnswersTab{
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
			message.NewTextMessage(fmt.Sprintln(wenda.ActiveGroup)).SetGroupTarget(m.Target.Group),
		)
	}
}
