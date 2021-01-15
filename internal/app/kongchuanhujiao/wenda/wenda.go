package wenda

import (
	"strings"

	"coding.net/kongchuanhujiao/server/internal/app/client/clientmsg"
	"coding.net/kongchuanhujiao/server/internal/app/datahub/datahubpkg/wenda"
	"coding.net/kongchuanhujiao/server/internal/app/kongchuanhujiao/public/wendapkg"
)

// HandleAnswer 处理消息中可能存在的答案
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
		_ = wenda.InsertAnswer(&wendapkg.AnswersTab{
			Question: uint32(qid),
			QQ:       m.Target.ID,
			Answer:   strings.ToUpper(answer),
		})
	case 2: // 多选题

	case 3: // 简答题
		if !checkAnswerForFill(answer) {
			return
		}
		_ = wenda.InsertAnswer(&wendapkg.AnswersTab{
			Question: uint32(qid),
			QQ:       m.Target.ID,
			Answer:   strings.TrimPrefix(answer, "#"),
		})
	}
}
