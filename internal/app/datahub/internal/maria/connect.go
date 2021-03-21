package maria

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

var (
	db     *sqlx.DB
	Logger = zap.L().Named("Maria")
)

// Connect 连接至 Maria 数据库
func Connect(url string) {
	d, err := sqlx.Connect("mysql", url)
	if err != nil {
		Logger.Panic("连接失败", zap.Error(err))
	}

	Logger.Debug("连接成功")
	db = d
}

// Select using this db.
// Any placeholder parameters are replaced with supplied args.
func Select(dest interface{}, query string, args ...interface{}) error {
	return db.Select(dest, query, args...)
}

// Exec executes a query without returning any rows.
// The args are for any placeholder parameters in the query.
func Exec(query string, args ...interface{}) (sql.Result, error) {
	return db.Exec(query, args...)
}
