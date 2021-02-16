package wenda

import (
	"coding.net/kongchuanhujiao/server/internal/app/client"
	"coding.net/kongchuanhujiao/server/internal/app/datahub/internal/maria"

	"github.com/elgris/sqrl"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
)

// SelectQuestions 获取问题
func SelectQuestions(v *QuestionsTab, page uint32) (data []*QuestionsTab, err error) {

	var sqr *sqrl.SelectBuilder
	if v.ID != 0 {
		sqr = sqrl.Select("*").Where("id=?", v.ID).Limit(1)
	} else {
		sqr = sqrl.Select("id", "question", "target", "`status`", `options`).Limit(20).
			Offset(uint64(page * 20))
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

	type questionsTab struct { // questionsTab 问题
		ID       uint32 `db:"id"`       // 唯一标识符
		Type     uint8  `db:"type"`     // 类型
		Subject  uint8  `db:"subject"`  // 学科
		Question string `db:"question"` // 问题
		Creator  string `db:"creator"`  // 创建者
		Target   uint64 `db:"target"`   // 目标
		Status   uint8  `db:"status"`   // 状态
		Options  string `db:"options"`  // 选项
		Key      string `db:"key"`      // 答案
		Market   bool   `db:"market"`   // 存在市场
	}

	var d []*questionsTab
	err = maria.DB.Select(&d, sql, args...)
	if err != nil {
		maria.Logger.Error("查询失败", zap.Error(err), zap.String("SQL语句", sql))
		return
	}

	for _, v := range d {

		var (
			q   = QuestionField{}
			o   []string
			err = jsoniter.UnmarshalFromString(v.Question, &q)
		)
		if err != nil {
			loggerr.Error("解析问题字段失败", zap.Error(err))
			return nil, err
		}

		err = jsoniter.UnmarshalFromString(v.Options, &o)
		if err != nil {
			loggerr.Error("解析选项字段失败", zap.Error(err))
			return nil, err
		}

		data = append(data, &QuestionsTab{
			ID: v.ID, Type: v.Type, Subject: v.Subject,
			Question: q,
			Creator:  v.Creator, Target: v.Target, Status: v.Status,
			Options: o,
			Key:     v.Key, Market: v.Market,
		})
	}
	return
}

// UpdateQuestionStatus 更新问题状态
// 当 status = 1 时， q 必须传入由 SelectQuestions 获取的
func UpdateQuestionStatus(q *QuestionsTab, status uint8) (err error) {

	sql, args, err := sqrl.Update("questions").Set("`status`", status).Where("id=?", q.ID).ToSql()
	if err != nil {
		loggerr.Error("生成SQL语句失败", zap.Error(err))
		return
	}

	_, err = maria.DB.Exec(sql, args...)
	if err != nil {
		maria.Logger.Error("更新失败", zap.Error(err), zap.String("SQL语句", sql))
	}

	switch status {
	case 0, 2: // 准备作答
		delete(ActiveGroup, q.Target)
	case 1: // 开始作答
		a, err := SelectAnswers(q.ID)
		if err != nil {
			return err
		}
		Caches[q.ID] = &Detail{
			Questions: q, Answers: a,
			Members: client.GetClient().GetGroupMembers(q.Target),
		}
		ActiveGroup[q.Target] = q.ID
		return sendQuestionMsg(q)
	}

	return
}

// InsertQuestion 新增问题
func InsertQuestion(q *QuestionsTab) (err error) {

	sql, args, err := sqrl.Insert("questions").Values(nil, q.Type, q.Subject, q.Question, q.Creator,
		q.Target, 0, q.Options, q.Key, q.Market).ToSql()
	if err != nil {
		loggerr.Error("生成SQL语句失败", zap.Error(err))
		return
	}

	_, err = maria.DB.Exec(sql, args...)
	if err != nil {
		maria.Logger.Error("插入失败", zap.Error(err), zap.String("SQL语句", sql))
	}

	return
}

// UpdateQuestion 更新问题
func UpdateQuestion(q *QuestionsTab) (err error) {

	sql, args, err := sqrl.Update("questions").Where("id=?", q.ID).
		Set("`subject`", q.Subject).
		Set("question", q.Question).
		Set("target", q.Target).
		Set("`options`", q.Options).
		Set("`key`", q.Key).
		Set("market", q.Market).ToSql()
	if err != nil {
		loggerr.Error("生成SQL语句失败", zap.Error(err))
		return
	}

	_, err = maria.DB.Exec(sql, args...)
	if err != nil {
		maria.Logger.Error("更新失败", zap.Error(err), zap.String("SQL语句", sql))
	}

	return
}

// CopyQuestions 复制问题
func CopyQuestions(id uint32, creator string, target uint64) (err error) {

	q, err := SelectQuestions(&QuestionsTab{ID: id, Market: true}, 0)
	if err != nil {
		return
	}
	que := q[0]
	que.Creator = creator
	que.Target = target
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
	_, err = maria.DB.Exec(sql, args...)
	if err != nil {
		maria.Logger.Error("删除失败", zap.Error(err), zap.String("SQL语句", sql))
	}

	return
}
