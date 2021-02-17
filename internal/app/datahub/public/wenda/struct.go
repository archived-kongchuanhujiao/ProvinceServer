package wenda

import "github.com/gorilla/websocket"

type (
	Runtime map[uint32][]*websocket.Conn // Runtime 运行时

	QuestionsTab struct { // QuestionsTab 问题
		ID       uint32        `json:"id" db:"id"`             // 唯一标识符
		Type     uint8         `json:"type" db:"type"`         // 类型
		Subject  uint8         `json:"subject" db:"subject"`   // 学科
		Question QuestionField `json:"question" db:"question"` // 问题
		Creator  string        `json:"creator" db:"creator"`   // 创建者
		Target   uint64        `json:"target" db:"target"`     // 目标
		Status   uint8         `json:"status" db:"status"`     // 状态
		Options  []string      `json:"options" db:"options"`   // 选项
		Key      string        `json:"key" db:"key"`           // 答案
		Market   bool          `json:"market" db:"market"`     // 存在市场
	}

	QuestionField []struct { // QuestionField 问题字段
		Type string `json:"type"` // 类型
		Data string `json:"data"` // 数据
	}

	AnswersTab struct { // AnswersTab 回答
		ID       uint32 `json:"id" db:"id"`             // 唯一标识符
		Question uint32 `json:"question" db:"question"` // 问题
		QQ       uint64 `json:"qq" db:"qq"`             // QQ
		Answer   string `json:"answer" db:"answer"`     // 回答
		Time     string `json:"time" db:"time"`         // 时刻
	}

	CalculationsTab struct { // CalculationsTab 计算结果
		Question uint32     `json:"question" db:"question"` // 问题
		Count    uint8      `json:"count" db:"count"`       // 作答人数
		Right    []uint64   `json:"right" db:"right"`       // 正确
		Wrong    [][]uint64 `json:"wrong" db:"wrong"`       // 错误
	}

	// Detail 问答详情
	// 内部缓存结构体
	// 不应该被返回给 API 客户端
	Detail struct {
		// 包含问题和回答还有群人的数据
		Questions *QuestionsTab // 问题
		Answers   []*AnswersTab // 回答
		Members   *GroupMembers // 群成员
	}

	Groups       map[uint64]string // Groups 群
	GroupMembers map[uint64]string // GroupMembers 群成员
)
