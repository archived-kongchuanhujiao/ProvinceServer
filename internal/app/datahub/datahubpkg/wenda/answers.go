package wenda

import (
	"coding.net/kongchuanhujiao/server/internal/app/datahub/internal/maria"

	"github.com/elgris/sqrl"
	"go.uber.org/zap"
)

// InsertAnswer 新增回答
func InsertAnswer(a *AnswersTab) (err error) {
	sql, args, err := sqrl.Insert("answers").Values(nil, a.Question, a.QQ, a.Answer, a.Time).ToSql()
	if err != nil {
		loggerr.Error("生成SQL语句失败", zap.Error(err))
		return
	}

	_, err = maria.DB.Exec(sql, args...)
	if err != nil {
		maria.Logger.Error("插入失败", zap.Error(err), zap.String("SQL语句", sql))
	}

	caches[a.Question].Answers = a
	return
}
