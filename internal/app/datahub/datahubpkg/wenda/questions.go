package wenda

import (
	"coding.net/kongchuanhujiao/server/internal/app/datahub/internal/maria"

	"github.com/elgris/sqrl"
	"go.uber.org/zap"
)

// GetQuestions 获取问题
func GetQuestions(page uint32, id uint32, market bool, subject uint8) (data []*QuestionListTab, err error) {

	sqr := sqrl.Select("*").From("questions").OrderBy("id DESC")
	if id != 0 {
		sqr = sqr.Where("id=?", id).Limit(1)
	} else {
		sqr = sqr.Limit(20).Offset(uint64(page * 20))
	}
	if market {
		sqr = sqr.Where("market=?", true)
	}
	if subject != 0 {
		sqr = sqr.Where("`subject`=?", subject)
	}

	sql, args, err := sqr.ToSql()
	if err != nil {
		maria.Logger.Error("生成SQL语句失败", zap.Error(err))
		return
	}

	err = maria.DB.Select(&data, sql, args...)
	if err != nil {
		maria.Logger.Error("查询失败", zap.Error(err), zap.String("SQL语句", sql))
		return
	}
	return
}

// UpdateQuestions 更新问题
func UpdateQuestions(id uint32, status uint8) (err error) {
	sql, args, err := sqrl.Update("questions").Set("`status`", status).Where("id=?", id).ToSql()
	if err != nil {
		maria.Logger.Error("生成SQL语句失败", zap.Error(err))
		return err
	}

	_, err = maria.DB.Exec(sql, args...)
	if err != nil {
		maria.Logger.Error("更新失败", zap.Error(err), zap.String("SQL语句", sql))
	}
	return
}

// CopyQuestions 复制问题
func CopyQuestions(id uint32, creator string, target uint64) (err error) {
	q, err := GetQuestions(0, id, true, 0)
	if err != nil {
		return
	}
	que := q[0]
	sql, args, err := sqrl.Insert("questions").Values(nil, que.Type, que.Subject, que.Question,
		creator, target, 0, que.Options, que.Key, false).ToSql()
	_, err = maria.DB.Exec(sql, args...)
	if err != nil {
		maria.Logger.Error("复制失败", zap.Error(err), zap.String("SQL语句", sql))
	}
	return
}
