package cuoti

import (
	"github.com/kongchuanhujiao/server/internal/app/datahub/public/cuoti"
)

// InsertWrongQuestion 插入错题
// TODO: 数据库交互
func InsertWrongQuestion(data *cuoti.Tab) (err error) {
	return
}

// SelectWrongQuestions 获取错题
// TODO: 数据库交互
func SelectWrongQuestions(id uint32, qid uint32) (data []*cuoti.Tab, err error) {
	return
}

// RemoveWrongQuestions 删除错题
// TODO: 数据库交互
func RemoveWrongQuestions(id uint32, qid uint64) (err error) {
	return
}
