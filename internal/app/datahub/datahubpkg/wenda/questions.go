package wenda

import (
	"coding.net/kongchuanhujiao/server/internal/app/datahub/internal/maria"

	"github.com/elgris/sqrl"
	"go.uber.org/zap"
)

// GetQuestions 获取问题
func GetQuestions(page uint32, id uint32, market bool, sub uint8) (data []QuestionListTab) {

	sqr := sqrl.Select("*").From("questions").OrderBy("questions.id DESC")
	if id != 0 {
		sqr = sqr.Where("questions.id=?", id).Limit(1)
	} else {
		if market {
			sqr = sqr.Where("questions.market=?", 1).Where("questions.`subject`=?", sub)
		}
		sqr = sqr.Limit(20).Offset(uint64(page * 20))
	}

	sql, args, err := sqr.ToSql()
	if err != nil {
		maria.Logger.Error("生成SQL语句失败", zap.Error(err))
		return nil
	}

	err = maria.DB.Select(&data, sql, args...)
	if err != nil {
		maria.Logger.Error("查询失败", zap.Error(err), zap.String("SQL语句", sql))
		return nil
	}
	return data
}
