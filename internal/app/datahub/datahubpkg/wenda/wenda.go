package wenda

import (
	"coding.net/kongchuanhujiao/server/internal/app/datahub/internal/maria"

	"github.com/elgris/sqrl"
	"go.uber.org/zap"
)

// GetQuestions 获取问题
func GetQuestions(v *QuestionsTab, page uint32) (data []*QuestionsTab, err error) {

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
	q, err := GetQuestions(&QuestionsTab{ID: id, Market: true}, 0)
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
