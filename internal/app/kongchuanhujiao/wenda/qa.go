package wenda

import (
	"io/ioutil"
	"strings"

	"coding.net/kongchuanhujiao/server/internal/app/client"
	"coding.net/kongchuanhujiao/server/internal/app/client/clientmsg"
	"coding.net/kongchuanhujiao/server/internal/app/datahub/datahubpkg/wenda"
	"coding.net/kongchuanhujiao/server/internal/pkg/logger"

	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
)

// StartQA 开始作答
// i 问题（问答） ID
func StartQA(i uint32) (err error) {
	q, err := wenda.SelectQuestions(&wenda.QuestionsTab{ID: i}, 0)
	if err != nil {
		return
	}

	que := q[0]

	if err = wenda.UpdateQuestionStatus(que, 1); err != nil {
		return
	}

	logger.Info("问答开始", zap.Uint32("ID", i))
	return sendQuestionMsg(que)
}

// sendQuestionMsg 发送问答题干
func sendQuestionMsg(q *wenda.QuestionsTab) (err error) {
	var (
		question []struct {
			Type string `json:"type"` // 类型
			Data string `json:"data"`
		}
		options []string
		json    = jsoniter.ConfigCompatibleWithStandardLibrary
	)

	if err = json.UnmarshalFromString(q.Question, &question); err != nil {
		logger.Error("解析问题失败", zap.Error(err))
		return
	}
	if err = json.UnmarshalFromString(q.Options, &options); err != nil {
		logger.Error("解析选项失败", zap.Error(err))
		return
	}

	m := clientmsg.NewTextMessage("问题:\n")
	for _, v := range question {
		if v.Type == "img" {
			f, err := ioutil.ReadFile("assets/question/pictures/" + v.Data)
			if err != nil {
				logger.Error("读取题干图片失败", zap.Error(err))
				return err
			}
			m.AddImage(f).AddText("\n")
			continue
		}
		m.AddText(v.Data + "\n")
	}

	m.AddText("选项:\n")
	abc := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I"}
	for k, v := range options {
		m.AddText(abc[k] + ". " + v + "\n")
	}

	if q.Type == 0 {
		m.AddText("\n回复选项即可作答")
	} else {
		m.AddText("\n@+回答内容即可作答")
	}

	client.GetClient().SendMessage(m.SetGroupTarget(&clientmsg.Group{ID: q.Target}))
	return
}

// handleAnswer 处理消息中可能存在的答案
func handleAnswer(m *clientmsg.Message) {

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
		_ = wenda.InsertAnswer(&wenda.AnswersTab{
			Question: uint32(qid),
			QQ:       m.Target.ID,
			Answer:   strings.ToUpper(answer),
		})
	case 2: // 多选题

	case 3: // 简答题
		if !checkAnswerForFill(answer) {
			return
		}
		_ = wenda.InsertAnswer(&wenda.AnswersTab{
			Question: uint32(qid),
			QQ:       m.Target.ID,
			Answer:   strings.TrimPrefix(answer, "#"),
		})

	}
}
