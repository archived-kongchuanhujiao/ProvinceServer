package accounts

type Tab struct { // Tab 账号
	ID    string `json:"id" db:"id"`       // 唯一标识符
	Type  uint8  `json:"type" db:"type"`   // 类型
	QQ    uint64 `json:"qq" db:"qq"`       // QQ
	Class uint32 `json:"class" db:"class"` // 班级
	Push  string `json:"push" db:"push"`   // 推送
	Token string `json:"token" db:"token"` // Token
}
