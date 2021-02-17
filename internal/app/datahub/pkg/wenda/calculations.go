package wenda

import (
	"coding.net/kongchuanhujiao/server/internal/app/datahub/internal/maria"
	"coding.net/kongchuanhujiao/server/internal/app/datahub/public/wenda"

	"github.com/elgris/sqrl"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
)

// InsertCalculations 插入计算结果
func InsertCalculations(c *wenda.CalculationsTab) (err error) {

	loggerr.Info("插入计算结果", zap.Uint32("问答ID", c.Question))

	sql, args, err := sqrl.Insert("calculations").Values(c.Question, c.Count, c.Right, c.Wrong).ToSql()
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

// SelectCalculations 获取计算结果
// qid 问题 ID
func SelectCalculations(qid uint32) (data []*wenda.CalculationsTab, err error) {

	sql, args, err := sqrl.Select("*").From("calculations").Where("question=?", qid).ToSql()
	if err != nil {
		loggerr.Error("生成SQL语句失败", zap.Error(err))
		return
	}

	type calculationsTab struct { // calculationsTab 计算结果
		Question    uint32 `db:"question"`     // 问题
		AnswerCount uint8  `db:"answer_count"` // 作答人数
		Right       string `db:"right"`        // 正确
		Wrong       string `db:"wrong"`        // 错误
	}

	var d []*calculationsTab
	err = maria.Select(&d, sql, args...)
	if err != nil {
		maria.Logger.Error("查询失败", zap.Error(err), zap.String("SQL语句", sql))
		return
	}

	for _, v := range d {

		var r []uint64
		var w [][]uint64

		err := jsoniter.UnmarshalFromString(v.Right, &r)
		if err != nil {
			loggerr.Error("解析问题字段失败", zap.Error(err))
			return nil, err
		}

		err = jsoniter.UnmarshalFromString(v.Wrong, &w)
		if err != nil {
			loggerr.Error("解析选项字段失败", zap.Error(err))
			return nil, err
		}

		data = append(data, &wenda.CalculationsTab{
			Question: v.Question, Count: v.AnswerCount,
			Right: r,
			Wrong: w,
		})
	}

	return
}
