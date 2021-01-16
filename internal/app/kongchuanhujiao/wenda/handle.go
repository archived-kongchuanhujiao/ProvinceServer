package wenda

import (
	"coding.net/kongchuanhujiao/server/internal/pkg/logger"
	"fmt"
	"go.uber.org/zap"
	"strings"

	"coding.net/kongchuanhujiao/server/internal/app/client"
	"coding.net/kongchuanhujiao/server/internal/app/client/clientmsg"
	"coding.net/kongchuanhujiao/server/internal/app/datahub/datahubpkg/wenda"
	"coding.net/kongchuanhujiao/server/internal/app/kongchuanhujiao/public/wendapkg"
)

// HandleAnswer 处理回答
func HandleAnswer(m *clientmsg.Message) {

	qid, ok := wenda.ActiveGroup[m.Target.Group.ID]
	if !ok {
		return
	}

	ans, ok := m.Chain[0].(*clientmsg.Text)
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
		logger.Info("答题有效", zap.Uint32("问答ID", uint32(q.Questions.ID)))
		_ = wenda.InsertAnswer(&wendapkg.AnswersTab{
			Question: qid,
			QQ:       m.Target.ID,
			Answer:   strings.ToUpper(answer),
		})
	case 2: // 多选题

	case 3: // 简答题
		if !checkAnswerForFill(answer) {
			return
		}
		logger.Info("答题有效", zap.Uint32("问答ID", uint32(q.Questions.ID)))
		_ = wenda.InsertAnswer(&wendapkg.AnswersTab{
			Question: qid,
			QQ:       m.Target.ID,
			Answer:   strings.TrimPrefix(answer, "#"),
		})
	}
}

// HandleTest 处理测试
func HandleTest(m *clientmsg.Message) {

	t, ok := m.Chain[0].(*clientmsg.Text)
	if !ok {
		return
	}

	switch t.Content {
	case "你好":
		client.GetClient().SendMessage(
			clientmsg.NewAtMessage(m.Target.ID).AddText("你好").SetGroupTarget(m.Target.Group),
		)
	case "活动的群":
		client.GetClient().SendMessage(
			clientmsg.NewTextMessage(fmt.Sprintln(wenda.ActiveGroup)).SetGroupTarget(m.Target.Group),
		)
	}
}
