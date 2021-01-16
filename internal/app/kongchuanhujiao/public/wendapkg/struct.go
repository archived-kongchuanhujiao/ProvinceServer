package wendapkg

import (
	"time"
)

type (
	AnswersTab struct { // AnswersTab 回答
		ID       uint32     `json:"id" db:"id"`             // 唯一标识符
		Question QuestionID `json:"question" db:"question"` // 问题
		QQ       uint64     `json:"qq" db:"qq"`             // QQ
		Answer   string     `json:"answer" db:"answer"`     // 回答
		Time     time.Time  `json:"time" db:"time"`         // 时刻
	}

	CalculationsTab struct { // CalculationsTab 计算结果
		Question uint32 `json:"question" db:"question"` // 问题
		Data     string `json:"data" db:"data"`         // 数据 TODO string 更改为具体的结构体
	}

	// WendaDetails 问答详情
	// 内部缓存结构体
	// 不应该被返回给 API 客户端
	WendaDetails struct {
		// 包含问题和回答还有群人的数据
		Questions *QuestionsTab // 问题
		Answers   []*AnswersTab // 回答
		Members   *GroupMembers // 群成员
	}

	Groups       map[uint64]string // Groups 群
	GroupMembers map[uint64]string // GroupMembers 群成员
)
