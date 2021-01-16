package wendapkg

import (
	"time"
)

type (
	QuestionID uint32 // 问题（问答） ID

	QuestionsTab struct { // QuestionsTab 问题
		ID       QuestionID    `json:"id" db:"id"`             // 唯一标识符
		Type     uint8         `json:"type" db:"type"`         // 类型
		Subject  uint8         `json:"subject" db:"subject"`   // 学科
		Question QuestionField `json:"question" db:"question"` // 问题
		Creator  string        `json:"creator" db:"creator"`   // 创建者
		Target   uint64        `json:"target" db:"target"`     // 目标
		Status   uint8         `json:"status" db:"status"`     // 状态
		Options  OptionsField  `json:"options" db:"options"`   // 选项
		Key      string        `json:"key" db:"key"`           // 答案
		Market   bool          `json:"market" db:"market"`     // 存在市场
	}

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

	QuestionField []struct { // QuestionField 问题字段
		Type string `json:"type"` // 类型
		Data string `json:"data"` // 数据
	}

	OptionsField []string // OptionsField 选项字段

	Question struct { // 问题
		Type string `json:"type"` // 题目类型
		Text string `json:"text"` // 题目文本
		Path string `json:"path"` // 题目图片, 当 type 为 img 时不为空
	}

	Option struct { // 选项
		Type string `json:"type"` // 选项
		Body string `json:"body"` // 答案
	}
)
