package wenda

import (
	"time"

	"coding.net/kongchuanhujiao/server/internal/app/datahub/internal/maria"
	"coding.net/kongchuanhujiao/server/internal/app/kongchuanhujiao/public/wendapkg"

	"github.com/elgris/sqrl"
	"go.uber.org/zap"
)

// InsertAnswer 新增回答
// 结构体中无需传 ID、Time
func InsertAnswer(a *wendapkg.AnswersTab) (err error) {
	a.Time = time.Now()
	sql, args, err := sqrl.Insert("answers").Values(nil, a.Question, a.QQ, a.Answer, a.Time).ToSql()
	if err != nil {
		loggerr.Error("生成SQL语句失败", zap.Error(err))
		return
	}

	_, err = maria.DB.Exec(sql, args...)
	if err != nil {
		maria.Logger.Error("插入失败", zap.Error(err), zap.String("SQL语句", sql))
	}

	q := Caches[wendapkg.QuestionID(a.Question)]
	q.Answers = append(q.Answers, a)
	PushData(q.Questions.ID, q.Answers)
	return
}

// SelectAnswers 获取回答
// qid 问题 ID
func SelectAnswers(qid wendapkg.QuestionID) (data []*wendapkg.AnswersTab, err error) {
	sql, args, err := sqrl.Select("*").From("answers").
		Where("question=?", qid).OrderBy("id DESC").ToSql()
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
