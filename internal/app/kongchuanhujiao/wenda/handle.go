package wenda

import (
	"fmt"
	"github.com/kongchuanhujiao/server/internal/app/datahub/pkg/cuoti"
	ct "github.com/kongchuanhujiao/server/internal/app/datahub/public/cuoti"
	"github.com/kongchuanhujiao/server/internal/pkg/config"
	"github.com/kongchuanhujiao/server/internal/pkg/logger"
	"go.uber.org/zap"
	"strconv"
	"strings"

	"github.com/kongchuanhujiao/server/internal/app/client"
	"github.com/kongchuanhujiao/server/internal/app/client/message"
	"github.com/kongchuanhujiao/server/internal/app/datahub/pkg/wenda"
	public "github.com/kongchuanhujiao/server/internal/app/datahub/public/wenda"
)

// sessionPool 储存添加错题进度
var sessionPool map[uint64]*ct.Tab

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
	}
}

// HandleWrongQuestion 处理添加错题逻辑
func HandleWrongQuestion(m *message.Message) {
	t, ok := m.Chain[0].(*message.Text)
	if !ok {
		return
	}

	if _, ok := sessionPool[m.Target.ID]; ok {
		handleAddCuoti(m.Target.ID, m.Target.Group, m.Chain[0])
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
				if _, ok := sessionPool[m.Target.ID]; !ok {
					sessionPool[m.Target.ID] = &ct.Tab{}
					client.GetClient().SendMessage(message.NewAtMessage(m.Target.ID).AddText("请发送欲添加错题的题目:").SetGroupTarget(m.Target.Group))
					return
				}

				client.GetClient().SendMessage(message.NewAtMessage(m.Target.ID).AddText("你还有正在进行添加的错题!").SetGroupTarget(m.Target.Group))
			case "del":
				return
			case "ck":
				wq, err := cuoti.SelectWrongQuestions(0, uint32(m.Target.ID))

				if err != nil {
					client.GetClient().SendMessage(message.NewAtMessage(m.Target.ID).AddText("发生了意外错误, 无法查询错题列表.").SetGroupTarget(m.Target.Group))
					return
				}

				client.GetClient().SendMessage(message.NewAtMessage(m.Target.ID).AddText(fmt.Sprint(wq)).SetGroupTarget(m.Target.Group))
			default:
				client.GetClient().SendMessage(defaultMsg)
			}
		}
	}
}

// handleAddCuoti 处理添加错题逻辑
func handleAddCuoti(user uint64, g *message.Group, chain message.Element) {
	stat, ok := sessionPool[user]

	if !ok {
		return
	}

	switch {
	case stat.Question == "":
		q := chain.(*message.Text).Content

		tab, err := cuoti.SelectWrongQuestions(0, uint32(user))

		if err != nil {
			client.GetClient().SendMessage(message.NewAtMessage(user).AddText("添加问题失败, 服务器异常.").SetGroupTarget(g))
			return
		}

		for _, t := range tab {
			if t.Question == q {
				client.GetClient().SendMessage(message.NewAtMessage(user).AddText("已经有相同题目的错题了!").SetGroupTarget(g))
				return
			}
		}

		stat.Question = q

		client.GetClient().SendMessage(message.NewAtMessage(user).AddText("设置错题题目成功! 下一步请设置错题的正确答案.").SetGroupTarget(g))
	case stat.QuestionAnswer == "":
		stat.QuestionAnswer = chain.(*message.Text).Content
		client.GetClient().SendMessage(message.NewAtMessage(user).AddText("设置错题答案成功! 如果该题目有图片需要添加的话, 接下来请发送图片, 反之则发送提醒周期 (单位为天).").SetGroupTarget(g))
	case stat.ImageURL == "" || stat.Duration == 0:
		i, success := chain.(*message.Image)
		if success && stat.ImageURL != "" {
			stat.ImageURL = i.URL
			client.GetClient().SendMessage(message.NewAtMessage(user).AddText("设置错题图片成功! 接下来请发送提醒时间 (单位为天).").SetGroupTarget(g))
		} else {
			d, err := strconv.Atoi(chain.(*message.Text).Content)
			if err != nil {
				client.GetClient().SendMessage(message.NewAtMessage(user).AddText("时间错误! 请填写正确的数字.").SetGroupTarget(g))
				return
			} else {
				stat.Duration = uint8(d)
				err := cuoti.InsertWrongQuestion(stat)
				if err != nil {
					client.GetClient().SendMessage(message.NewAtMessage(user).AddText("错题创建成功!").SetGroupTarget(g))
				} else {
					client.GetClient().SendMessage(message.NewAtMessage(user).AddText("添加错题失败!").SetGroupTarget(g))
					logger.Warn("添加错题时遇到了意外", zap.Error(err))
				}
			}
		}
	}
}
