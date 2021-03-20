package wenda

import "github.com/gorilla/websocket"

type (
	Runtime []*WSConn // Runtime 运行时

	WSConn struct {
		Conn           *websocket.Conn // Websocket 连接
		ListenQuestion uint32          // 监听的问题题号
		ListenGroup    uint64          // 监听消息的群聊
	}

	// AnswersTab 回答表
	AnswersTab struct {
		Question uint32 `json:"question" db:"question"` // 问题标识号
		QQ       uint64 `json:"qq"       db:"qq"`       // QQ 标识号
		Answer   string `json:"answer"   db:"answer"`   // 作答内容
		Time     string `json:"time"     db:"time"`     // 时刻
		Mark     string `json:"mark"     db:"mark"`     // 标记
	}

	// Result 结果
	Result struct {
		Count uint8              `json:"count" db:"count"` // 作答人数
		Right []uint64           `json:"right" db:"right"` // 正确学生
		Wrong []ResultWrongField `json:"wrong" db:"wrong"` // 错误学生
	}

	// ResultWrongField 结果错误学生字段
	ResultWrongField struct {
		Type  string   `json:"type"`  // 类型
		Value []uint64 `json:"value"` // 学生
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
