package wenda

import (
	"time"

	"github.com/kongchuanhujiao/server/internal/app/client"
	"github.com/kongchuanhujiao/server/internal/app/datahub/internal/maria"
	"github.com/kongchuanhujiao/server/internal/app/datahub/public/wenda"

	"github.com/elgris/sqrl"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
)

type (
	// questionsTab 问题表
	questionsTab struct {
		ID      uint32 `db:"id"`      // 标识号
		Type    uint8  `db:"type"`    // 类型
		Subject uint8  `db:"subject"` // 学科
		Creator string `db:"creator"` // 创建者
		Date    string `db:"date"`    // 创建日期
		Topic   string `db:"topic"`   // 主题
		Status  uint8  `db:"status"`  // 状态
		Market  bool   `db:"market"`  // 是否发布至问题市场
	}
)

// SelectQuestions 获取问题
func SelectQuestions(v *wenda.QuestionsTab, page uint32) (data []*wenda.QuestionsTab, err error) {

	var sqr *sqrl.SelectBuilder
	if v.ID != 0 {
		sqr = sqrl.Select("*").Where("id=?", v.ID).Limit(1)
	} else {
		sqr = sqrl.Select("id", "topic", "`status`").Limit(20).Offset(uint64(page * 20))
	}

	sqr = sqr.From("questions").OrderBy("id DESC")
	if v.Creator != "" {
		sqr = sqr.Where("creator=?", v.Creator)
	}
	if v.Market {
		sqr = sqr.Where("market=?", true)
	}
	if v.Subject != 0 {
		sqr = sqr.Where("`subject`=?", v.Subject)
	}

	sql, args, err := sqr.ToSql()
	if err != nil {
		loggerr.Error("生成SQL语句失败", zap.Error(err))
		return
	}

	var d []*questionsTab
	err = maria.Select(&d, sql, args...)
	if err != nil {
		maria.Logger.Error("查询失败", zap.Error(err), zap.String("SQL语句", sql))
		return
	}

	for _, v := range d {
		to := wenda.QuestionsTopicField{}
		err := jsoniter.UnmarshalFromString(v.Topic, &to)
		if err != nil {
			loggerr.Error("解析主题字段失败", zap.Error(err))
			return nil, err
		}

		data = append(data, &wenda.QuestionsTab{
			ID: v.ID, Type: v.Type, Subject: v.Subject, Topic: to, Creator: v.Creator, Status: v.Status,
			Market: v.Market,
		})
	}
	return
}

// UpdateQuestionStatus 更新问题状态
func UpdateQuestionStatus(q *wenda.QuestionsTab, status uint8) (err error) {

	sql, args, err := sqrl.Update("questions").Set("`status`", status).Where("id=?", q.ID).ToSql()
	if err != nil {
		loggerr.Error("生成SQL语句失败", zap.Error(err))
		return
	}

	_, err = maria.Exec(sql, args...)
	if err != nil {
		maria.Logger.Error("更新失败", zap.Error(err), zap.String("SQL语句", sql))
	}

	switch status {
	case 0, 2: // 准备作答
		DeleteActiveGroup(q.Topic.Target)
	case 1: // 开始作答
		a, err := SelectAnswers(q.ID)
		if err != nil {
			return err
		}

		WriteCaches(q.ID, &wenda.Detail{
			Questions: q, Answers: a, Members: client.GetClient().GetGroupMembers(q.Topic.Target),
		})
		WriteActiveGroup(q.Topic.Target, q.ID)
		return sendQuestionMsg(q)
	}

	return
}

// InsertQuestion 新增问题
func InsertQuestion(q *wenda.QuestionsTab) (err error) {

	sql, args, err := sqrl.Insert("questions").Values(nil, q.Type, q.Subject, q.Creator, q.Date, q.Topic,
		nil, q.Market).ToSql()
	if err != nil {
		loggerr.Error("生成SQL语句失败", zap.Error(err))
		return
	}

	_, err = maria.Exec(sql, args...)
	if err != nil {
		maria.Logger.Error("插入失败", zap.Error(err), zap.String("SQL语句", sql))
	}

	return
}

// UpdateQuestion 更新问题
func UpdateQuestion(q *wenda.QuestionsTab) (err error) {

	sql, args, err := sqrl.Update("questions").Where("id=?", q.ID).
		Set("`subject`", q.Subject).Set("topic", q.Topic).Set("market", q.Market).ToSql()
	if err != nil {
		loggerr.Error("生成SQL语句失败", zap.Error(err))
		return
	}

	_, err = maria.Exec(sql, args...)
	if err != nil {
		maria.Logger.Error("更新失败", zap.Error(err), zap.String("SQL语句", sql))
	}

	return
}

// CopyQuestions 复制问题
func CopyQuestions(id uint32, creator string, target uint64) (err error) {

	q, err := SelectQuestions(&wenda.QuestionsTab{ID: id, Market: true}, 0)
	if err != nil {
		return
	}
	que := q[0]
	que.Date = time.Now()
	que.Creator = creator
	que.Topic.Target = target
	que.Market = false
	err = InsertQuestion(que)

	return
}

// DeleteQuestion 删除问题
func DeleteQuestion(id uint32) (err error) {

	sql, args, err := sqrl.Delete("questions").Where("id=?", id).ToSql()
	if err != nil {
		loggerr.Error("生成SQL语句失败", zap.Error(err))
		return
	}
	_, err = maria.Exec(sql, args...)
	if err != nil {
		maria.Logger.Error("删除失败", zap.Error(err), zap.String("SQL语句", sql))
	}

	return
}
