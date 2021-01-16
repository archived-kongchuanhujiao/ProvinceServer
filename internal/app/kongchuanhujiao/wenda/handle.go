package wenda

import (
	"fmt"
	"strings"

	"coding.net/kongchuanhujiao/server/internal/app/client"
	"coding.net/kongchuanhujiao/server/internal/app/client/clientmsg"
	"coding.net/kongchuanhujiao/server/internal/app/datahub/datahubpkg/wenda"
	"coding.net/kongchuanhujiao/server/internal/app/kongchuanhujiao/public/wendapkg"
	"coding.net/kongchuanhujiao/server/internal/pkg/logger"
)

// HandleAnswer 处理消息中可能存在的答案
func HandleAnswer(m *clientmsg.Message) {

	logger.Debug("开始处理答题")

	qid, ok := wenda.ActiveGroup[m.Target.Group.ID]
	if !ok {
		return
	}

	logger.Debug("是活动的答题")

	ans, ok := m.Chain[0].(*clientmsg.Text)
	if !ok {
		return
	}
	answer := ans.Content

	logger.Debug("成功获取答题内容")

	q := wenda.Caches[qid]
	for _, v := range q.Answers {
		if v.QQ == m.Target.ID {
			return
		}
	}

	// TODO 检查答题是否有效

	logger.Debug("是有效的答题")

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
