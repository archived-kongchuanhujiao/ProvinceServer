package wenda

import (
	"coding.net/kongchuanhujiao/server/internal/app/datahub/internal/maria"

	"github.com/elgris/sqrl"
	"go.uber.org/zap"
)

// SelectQuestions 获取问题
func SelectQuestions(v *QuestionsTab, page uint32) (data []*QuestionsTab, err error) {

	sqr := sqrl.Select("*").From("questions").OrderBy("id DESC")
	if v.Creator != "" {
		sqr = sqr.Where("creator=?", v.Creator).Limit(20).Offset(uint64(page * 20))
	}
	if v.ID != 0 {
		sqr = sqr.Where("id=?", v.ID).Limit(1)
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

	err = maria.DB.Select(&data, sql, args...)
	if err != nil {
		maria.Logger.Error("查询失败", zap.Error(err), zap.String("SQL语句", sql))
		return
	}
	return
}

// UpdateQuestionStatus 更新问题状态
func UpdateQuestionStatus(id uint32, status uint8) (err error) {
	sql, args, err := sqrl.Update("questions").Set("`status`", status).Where("id=?", id).ToSql()
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

// InsertQuestion 新建问题
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
