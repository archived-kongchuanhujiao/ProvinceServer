package wendapkg

type (
	AnswersTab struct { // AnswersTab 回答
		ID       uint32     `json:"id" db:"id"`             // 唯一标识符
		Question QuestionID `json:"question" db:"question"` // 问题
		QQ       uint64     `json:"qq" db:"qq"`             // QQ
		Answer   string     `json:"answer" db:"answer"`     // 回答
		Time     string     `json:"time" db:"time"`         // 时刻
	}

	CalculationsTab struct { // CalculationsTab 计算结果
		Question    QuestionID `json:"question" db:"question"`         // 问题
		AnswerCount uint8      `json:"answer_count" db:"answer_count"` // 作答人数
		Right       []uint64   `json:"right" db:"right"`               // 正确
		Wrong       [][]uint64 `json:"wrong" db:"wrong"`               // 错误
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
