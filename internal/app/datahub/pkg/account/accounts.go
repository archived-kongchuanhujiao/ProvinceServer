package account

import (
	"github.com/kongchuanhujiao/server/internal/app/datahub/internal/maria"
	"github.com/kongchuanhujiao/server/internal/app/datahub/public/account"
	"github.com/kongchuanhujiao/server/internal/pkg/logger"

	"github.com/Masterminds/squirrel"
	"go.uber.org/zap"
)

var loggerr = logger.Named("数据总线").Named("账号")

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

	err = maria.Select(&data, sql, args...)
	if err != nil {
		maria.Logger.Error("查询失败", zap.Error(err), zap.String("SQL语句", sql))
		return
	}

	return
}
