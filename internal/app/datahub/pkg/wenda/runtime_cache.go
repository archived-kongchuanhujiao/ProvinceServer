package wenda

import (
	"io/ioutil"

	"github.com/kongchuanhujiao/server/internal/app/client"
	"github.com/kongchuanhujiao/server/internal/app/client/message"
	"github.com/kongchuanhujiao/server/internal/app/datahub/public/wenda"

	"go.uber.org/zap"
)

// sendQuestionMsg 发送问答题干
func sendQuestionMsg(q *wenda.QuestionsTab) (err error) {

	m := message.NewTextMessage("问题:\n")
	for _, v := range q.Topic.Question {
		if v.Type == "img" {
			f, err := ioutil.ReadFile("assets/pictures/questions/" + q.Creator + "/" + v.Data)
			if err != nil {
				zap.L().Error("读取题干图片失败", zap.Error(err))
				return err
			}
			m.AddImage(f).AddText("\n")
			continue
		}
		m.AddText(v.Data + "\n")
	}

	m.AddText("选项:\n")
	abc := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I"}
	for k, v := range q.Topic.Options {
		m.AddText(abc[k] + ". " + v + "\n")
	}

	if q.Type == 0 {
		m.AddText("\n回复选项即可作答")
	} else {
		m.AddText("\n@+回答内容即可作答")
	}

	client.GetClient().SendMessage(m.SetGroupTarget(&message.Group{ID: q.Topic.Target}))
	return
}
