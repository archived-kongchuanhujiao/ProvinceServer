package wendapkg

type (
	QuestionField []struct { // QuestionField 问题字段
		Type string `json:"type"` // 类型
		Data string `json:"data"` // 数据
	}

	OptionsField []string // OptionsField 选项字段
)
