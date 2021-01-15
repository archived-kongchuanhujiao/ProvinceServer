package wendapkg

import "coding.net/kongchuanhujiao/server/internal/app/datahub/datahubpkg/wenda"

type (
	QuestionID uint32 // 问题（问答） ID

	// WendaDetails 问答详情
	// 内部缓存结构体
	// 不应该被返回给 API 客户端
	WendaDetails struct {
		// 包含问题和回答还有群人的数据
		Questions *wenda.QuestionsTab // 问题
		Answers   []*wenda.AnswersTab // 回答
		Members   *GroupMembers       // 群成员
	}

	Groups       map[uint64]string // Groups 群
	GroupMembers map[uint64]string // GroupMembers 群成员

	QuestionField []struct { // QuestionField 问题字段
		Type string `json:"type"` // 类型
		Data string `json:"data"` // 数据
	}

	OptionsField []string // OptionsField 选项字段
)
