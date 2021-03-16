package account

type (
	// Tab 账号
	Tab struct {
		ID    string    `json:"id" db:"id"`       // 标识号
		QQ    uint64    `json:"qq" db:"qq"`       // QQ
		Class uint32    `json:"class" db:"class"` // 班级 FIXME 是否需要
		Push  PushField `json:"push" db:"push"`   // 推送
	}

	// PushField 推送字段
	PushField []struct {
		Platform string `json:"platform"` // 平台
		Key      string `json:"key"`      // 键
	}
)
