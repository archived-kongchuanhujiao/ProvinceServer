package wenda

type (
	QuestionListTab struct { // QuestionListTab 问题
		ID        uint32 `json:"id" db:"id"`                 // 唯一标识符
		Type      uint8  `json:"type" db:"type"`             // 类型
		Subject   uint8  `json:"subject" db:"subject"`       // 学科
		Question  string `json:"question" db:"question"`     // 问题
		CreatorID string `json:"creator_id" db:"creator_id"` // 创建者
		Target    uint64 `json:"target" db:"target"`         // 目标
		Status    uint8  `json:"status" db:"status"`         // 状态
		Options   string `json:"options" db:"options"`       // 选项
		Key       string `json:"key" db:"key"`               // 答案
		Market    bool   `json:"market" db:"market"`         // 存在市场
	}
)
