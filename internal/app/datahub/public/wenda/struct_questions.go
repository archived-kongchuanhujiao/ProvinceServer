package wenda

import "time"

type (
	// QuestionsTab 问题表
	QuestionsTab struct {
		ID      uint32              `json:"id"      db:"id"`      // 标识号
		Type    uint8               `json:"type"    db:"type"`    // 类型
		Subject uint8               `json:"subject" db:"subject"` // 学科
		Creator string              `json:"creator" db:"creator"` // 创建者
		Date    time.Time           `json:"date"    db:"date"`    // 创建日期
		Topic   QuestionsTopicField `json:"topic"   db:"topic"`   // 主题
		Status  uint8               `json:"status"  db:"status"`  // 状态
		Market  bool                `json:"market"  db:"market"`  // 是否发布至问题市场
	}

	// QuestionsTopicField 问题主题字段
	QuestionsTopicField struct {
		Target uint64 `json:"target"` // 目标
		// 问题
		Question []struct {
			Type string `json:"type"`
			Data string `json:"data"`
		} `json:"question"`
		Options []string `json:"options"` // 选项
		Key     string   `json:"key"`     // 答案
	}
)
