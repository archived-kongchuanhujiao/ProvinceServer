package wenda

import (
	"coding.net/kongchuanhujiao/server/internal/app/datahub/internal/maria"
	"coding.net/kongchuanhujiao/server/internal/app/kongchuanhujiao/public/wendapkg"

	"github.com/elgris/sqrl"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
)

// SelectCalculations 获取数据
// qid 问题 ID
func SelectCalculations(qid wendapkg.QuestionID) (data []*wendapkg.CalculationsTab, err error) {
	sql, args, err := sqrl.Select("*").From("calculations").Where("question=?", qid).ToSql()
	if err != nil {
		loggerr.Error("生成SQL语句失败", zap.Error(err))
		return
	}

	type calculationsTab struct { // calculationsTab 计算结果
		Question    wendapkg.QuestionID `db:"question"`     // 问题
		AnswerCount uint8               `db:"answer_count"` // 作答人数
		Right       string              `db:"right"`        // 正确
		Wrong       string              `db:"wrong"`        // 错误
	}

	var d []*calculationsTab
	err = maria.DB.Select(&d, sql, args...)
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

		data = append(data, &wendapkg.CalculationsTab{
			Question: v.Question, AnswerCount: v.AnswerCount,
			Right: r,
			Wrong: w,
		})
	}
	return
}
