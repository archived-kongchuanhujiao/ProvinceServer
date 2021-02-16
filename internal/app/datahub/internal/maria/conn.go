package maria

import (
	"coding.net/kongchuanhujiao/server/internal/app/datahub/internal"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

var (
	DB     *sqlx.DB
	Logger = internal.Logger.Named("Maria")
)

// Connect 连接至 Maria 数据库
func Connect(url string) {
	db, err := sqlx.Connect("mysql", url)
	if err != nil {
		Logger.Panic("连接失败", zap.Error(err))
	}

	Logger.Debug("连接成功")
	DB = db
}
