package wenda

import (
	"time"

	"github.com/kongchuanhujiao/server/internal/app/datahub/internal/maria"
	"github.com/kongchuanhujiao/server/internal/app/datahub/internal/memory"
	"github.com/kongchuanhujiao/server/internal/app/datahub/public/wenda"

	"github.com/elgris/sqrl"
	"go.uber.org/zap"
)

// InsertAnswer 插入回答
// 结构体中无需传 ID、Time
func InsertAnswer(a *wenda.AnswersTab) (err error) {

	loggerr.Info("插入回答数据", zap.Uint32("问答ID", a.Question))
	sql, args, err := sqrl.Insert("answers").Values(a.Question, a.QQ, a.Answer,
		time.Now().Format("2006-01-02 15:04:05"), a.Mark).ToSql()
	if err != nil {
		loggerr.Error("生成SQL语句失败", zap.Error(err))
		return
	}

	_, err = maria.Exec(sql, args...)
	if err != nil {
		maria.Logger.Error("插入失败", zap.Error(err), zap.String("SQL语句", sql))
	}

	q := memory.Caches[a.Question]
	q.Answers = append(q.Answers, a)

	return
}

// SelectAnswers 获取回答
// qid 问题 ID
func SelectAnswers(qid uint32) (data []*wenda.AnswersTab, err error) {
	sql, args, err := sqrl.Select("*").From("answers").Where("question=?", qid).
		OrderBy("time DESC").ToSql()
	if err != nil {
		loggerr.Error("生成SQL语句失败", zap.Error(err))
		return
	}

	err = maria.Select(&data, sql, args...)
	if err != nil {
		maria.Logger.Error("查询失败", zap.Error(err), zap.String("SQL语句", sql))
		return
	}

	return
}
