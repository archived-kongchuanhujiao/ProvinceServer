package wenda

import (
	"io/ioutil"
	"strings"
	"time"

	"coding.net/kongchuanhujiao/server/internal/app/client"
	"coding.net/kongchuanhujiao/server/internal/app/client/clientmsg"
	"coding.net/kongchuanhujiao/server/internal/app/datahub/datahubpkg/wenda"
	"coding.net/kongchuanhujiao/server/internal/pkg/logger"

	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
)

// StartQA 使用 i：问题ID(ID) 开始作答
func StartQA(i uint32) (err error) {

	q, err := wenda.SelectQuestions(&wenda.QuestionsTab{ID: i}, 0)
	if err != nil {
		return
	}
	if err = wenda.UpdateQuestionStatus(i, 1); err != nil {
		return
	}

	que := q[0]

	logger.Info("问题开始监听", zap.Uint32("ID", i))
	if err = sendQuestionMsg(que); err != nil {
		return
	}

	que.Status = 1
	// TODO 写到 datahub
	//QABasicSrvPoll[q.Target] = que
	return
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

// ReadQuestion 使用 i：问题ID(ID) 读取问答信息
func ReadQuestion(i uint32) (q *Question, err error) {
	// FIXME
	res, err := a.DB.Question().ReadQuestion(i)
	if err != nil {
		return
	}

	answer, err := a.DB.Answer().ReadAnswerList(i)
	if err != nil {
		return
	}

	groupInfo := a.Cli.C.FindGroupByUin(int64(res.Target))
	mems := a.ReadMemInfo(uint64(groupInfo.Uin))

	return &Question{res, mems, groupInfo.Name, answer}, nil
}

// insertAnswer 新增回答
// qid 问题 ID
// qnum 学生 QQ 号
// ans 学生回答内容
func insertAnswer(q *wenda.QuestionsTab, qnum uint64, ans string) {
	_ = wenda.InsertAnswer(&wenda.AnswersTab{Question: q.ID, QQ: qnum, Answer: ans, Time: time.Now()})
}

// handleAnswer 处理消息中可能存在的答案
func handleAnswer(m *clientmsg.Message) {
	// FIXME
	q, ok := QABasicSrvPoll[m.Group.ID]
	if !ok {
		return
	}

	for _, v := range q.Answer {
		if v.AnswererID == m.User.ID {
			return
		}
	}

	switch q.Type {
	// 选择题
	case 0:
		if checkAnswerForSelect(m.Chain[0].Text) {
			a.writeAnswer(q, m.User.ID, strings.ToUpper(m.Chain[0].Text))
		}
	// 简答题
	case 1:
		if checkAnswerForFill(m.Chain[0].Text) {
			a.writeAnswer(q, m.User.ID, strings.TrimPrefix(m.Chain[0].Text, "#"))
		}
	// 多选题
	case 2:

	// 填空题
	case 3:

	}

}

// deleteQABasicSrvPoll 使用 i：问题ID(ID) 删除问答基本服务池字段 FIXME 这里逻辑有点大问题
func deleteQABasicSrvPoll(i uint32) (err error) {
	// FIXME
	q, err := a.DB.Question().ReadQuestion(i)
	if err != nil {
		return
	}

	delete(QABasicSrvPoll, q.Target)
	return
}
