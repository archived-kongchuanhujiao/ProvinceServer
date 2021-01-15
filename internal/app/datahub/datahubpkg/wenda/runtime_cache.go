package wenda

import (
	"io/ioutil"

	"coding.net/kongchuanhujiao/server/internal/app/client"
	"coding.net/kongchuanhujiao/server/internal/app/client/clientmsg"
	"coding.net/kongchuanhujiao/server/internal/app/kongchuanhujiao/public/wendapkg"
	"coding.net/kongchuanhujiao/server/internal/pkg/logger"

	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
)

var (
	Caches      = map[wendapkg.QuestionID]*wendapkg.WendaDetails{} // Caches 缓存
	ActiveGroup = map[uint64]wendapkg.QuestionID{}                 // ActiveGroup 活动的群
)

// sendQuestionMsg 发送问答题干 TODO 迁移至 datahub
func sendQuestionMsg(q *wendapkg.QuestionsTab) (err error) {
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
