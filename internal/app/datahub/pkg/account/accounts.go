package account

import (
	"github.com/kongchuanhujiao/server/internal/app/datahub/internal/maria"
	"github.com/kongchuanhujiao/server/internal/app/datahub/public/account"
	public "github.com/kongchuanhujiao/server/internal/app/datahub/public/account"

	"github.com/Masterminds/squirrel"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
)

type tab struct {
	ID   string `json:"id" db:"id"`     // 标识号
	QQ   uint64 `json:"qq" db:"qq"`     // QQ
	Push string `json:"push" db:"push"` // 推送
}

var loggerr = zap.L().Named("数据总线").Named("账号")

// SelectAccount 获取账号
func SelectAccount(id string, qq uint64) (data []*account.Tab, err error) {

	sqr := squirrel.Select("*").From("accounts")
	if id != "" {
		sqr = sqr.Where("id=?", id)
	} else {
		sqr = sqr.Where("qq=?", qq)
	}

	sql, args, err := sqr.Limit(1).ToSql()
	if err != nil {
		loggerr.Error("生成SQL语句失败", zap.Error(err))
		return
	}

	var d []*tab
	err = maria.Select(&d, sql, args...)
	if err != nil {
		maria.Logger.Error("查询失败", zap.Error(err), zap.String("SQL语句", sql))
		return
	}

	for _, v := range d {
		pu := public.PushField{}
		err := jsoniter.UnmarshalFromString(v.Push, &pu)
		if err != nil {
			loggerr.Error("解析推送字段失败", zap.Error(err))
			return nil, err
		}

		data = append(data, &public.Tab{ID: v.ID, QQ: v.QQ, Push: pu})
	}
	return
}
