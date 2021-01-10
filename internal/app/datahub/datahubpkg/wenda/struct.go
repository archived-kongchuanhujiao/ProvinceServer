package wenda

import (
	"coding.net/kongchuanhujiao/server/internal/pkg/logger"
	"time"

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

	AnswersTab struct { // AnswersTab 回答
		ID       uint32    `json:"id" db:"id"`             // 唯一标识符
		Question uint32    `json:"question" db:"question"` // 问题
		QQ       uint64    `json:"qq" db:"qq"`             // QQ
		Answer   string    `json:"answer" db:"answer"`     // 回答
		Time     time.Time `json:"time" db:"time"`         // 时刻
	}

	CalculationsTab struct { // CalculationsTab 计算结果
		Question uint32 `json:"question" db:"question"` // 问题
		Data     string `json:"data" db:"data"`         // 数据 TODO string 更改为具体的结构体
	}

	runtime map[uint32][]*websocket.Conn // runtime 运行时
)
