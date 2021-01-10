package wenda

import (
	"coding.net/kongchuanhujiao/server/internal/pkg/logger"

	"github.com/gorilla/websocket"
)

var loggerr = logger.Named("数据总线").Named("问答")

type (
	QuestionsTab struct { // QuestionsTab 问题
		ID       uint32 `json:"id" db:"id"`             // 唯一标识符
		Type     uint8  `json:"type" db:"type"`         // 类型
		Subject  uint8  `json:"subject" db:"subject"`   // 学科
		Question string `json:"question" db:"question"` // 问题
		Creator  string `json:"creator" db:"creator"`   // 创建者
		Target   uint64 `json:"target" db:"target"`     // 目标
		Status   uint8  `json:"status" db:"status"`     // 状态
		Options  string `json:"options" db:"options"`   // 选项
		Key      string `json:"key" db:"key"`           // 答案
		Market   bool   `json:"market" db:"market"`     // 存在市场
	}

	runtime map[uint32][]*websocket.Conn // runtime 运行时
)
